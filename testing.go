package weshnet

import (
	"bytes"
	"context"
	crand "crypto/rand"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	datastore "github.com/ipfs/go-datastore"
	ds_sync "github.com/ipfs/go-datastore/sync"
	keystore "github.com/ipfs/go-ipfs-keystore"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	libp2p_mocknet "github.com/libp2p/go-libp2p/p2p/net/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	encrepo "berty.tech/go-ipfs-repo-encrypted"
	orbitdb "berty.tech/go-orbit-db"
	"berty.tech/go-orbit-db/pubsub/pubsubraw"
	"berty.tech/weshnet/pkg/cryptoutil"
	"berty.tech/weshnet/internal/datastoreutil"
	"berty.tech/weshnet/pkg/ipfsutil"
	"berty.tech/berty/v2/go/pkg/testutil"
	"berty.tech/weshnet/pkg/errcode"
	"berty.tech/weshnet/pkg/protocoltypes"
)

const InMemoryDir = ":memory:"

func TestHelperIPFSSetUp(t *testing.T) (context.Context, context.CancelFunc, libp2p_mocknet.Mocknet, host.Host) {
	t.Helper()

	ctx, cancel := context.WithCancel(context.Background())

	mn := libp2p_mocknet.New()
	t.Cleanup(func() { mn.Close() })

	rdvp, err := mn.GenPeer()
	require.NoError(t, err, "failed to generate mocked peer")

	return ctx, cancel, mn, rdvp
}

func NewTestOrbitDB(ctx context.Context, t *testing.T, logger *zap.Logger, node ipfsutil.CoreAPIMock, baseDS datastore.Batching) *BertyOrbitDB {
	t.Helper()

	api := node.API()
	selfKey, err := api.Key().Self(ctx)
	if err != nil {
		t.Fatal(err)
	}

	baseDS = datastoreutil.NewNamespacedDatastore(baseDS, datastore.NewKey(selfKey.ID().String()))

	pubSub := pubsubraw.NewPubSub(node.PubSub(), selfKey.ID(), logger, nil)

	odb, err := NewBertyOrbitDB(ctx, api, &NewOrbitDBOptions{
		Datastore: baseDS,
		NewOrbitDBOptions: orbitdb.NewOrbitDBOptions{
			Logger: logger,
			PubSub: pubSub,
		},
	})
	require.NoError(t, err)

	return odb
}

type mockedPeer struct {
	CoreAPI ipfsutil.CoreAPIMock
	DB      *BertyOrbitDB
	GC      *GroupContext
	MKS     *cryptoutil.MessageKeystore
	DevKS   cryptoutil.DeviceKeystore
}

func (m *mockedPeer) PeerInfo() peer.AddrInfo {
	return m.CoreAPI.MockNode().Peerstore.PeerInfo(m.CoreAPI.MockNode().Identity)
}

type TestingProtocol struct {
	Opts *Opts

	Service Service
	Client  Client

	RootDatastore  datastore.Batching
	DeviceKeystore cryptoutil.DeviceKeystore
	IpfsCoreAPI    ipfsutil.ExtendedCoreAPI
	OrbitDB        *BertyOrbitDB
	GroupDatastore *cryptoutil.GroupDatastore
}

type TestingOpts struct {
	Logger         *zap.Logger
	Mocknet        libp2p_mocknet.Mocknet
	RDVPeer        peer.AddrInfo
	DeviceKeystore cryptoutil.DeviceKeystore
	CoreAPIMock    ipfsutil.CoreAPIMock
	OrbitDB        *BertyOrbitDB
	ConnectFunc    ConnectTestingProtocolFunc
	PushSK         *[32]byte
}

func NewTestingProtocol(ctx context.Context, t testing.TB, opts *TestingOpts, ds datastore.Batching) (*TestingProtocol, func()) {
	if opts == nil {
		opts = &TestingOpts{}
	}
	opts.applyDefaults(ctx, t)

	if ds == nil {
		ds = ds_sync.MutexWrap(datastore.NewMapDatastore())
	}

	ipfsopts := &ipfsutil.TestingAPIOpts{
		Logger:    opts.Logger,
		Mocknet:   opts.Mocknet,
		RDVPeer:   opts.RDVPeer,
		Datastore: ds,
	}

	node := opts.CoreAPIMock
	cleanupNode := func() {}
	if node == nil {
		node, cleanupNode = ipfsutil.TestingCoreAPIUsingMockNet(ctx, t, ipfsopts)
	}

	deviceKeystore := opts.DeviceKeystore
	if deviceKeystore == nil {
		deviceKeystore = cryptoutil.NewDeviceKeystore(ipfsutil.NewDatastoreKeystore(datastoreutil.NewNamespacedDatastore(ds, datastore.NewKey(NamespaceDeviceKeystore))), nil)
	}

	odb := opts.OrbitDB
	if odb == nil {
		var err error

		pubSub := pubsubraw.NewPubSub(node.PubSub(), node.MockNode().PeerHost.ID(), opts.Logger, nil)

		odb, err = NewBertyOrbitDB(ctx, node.API(), &NewOrbitDBOptions{
			NewOrbitDBOptions: orbitdb.NewOrbitDBOptions{
				PubSub: pubSub,
				Logger: opts.Logger,
			},
			Datastore:      ds,
			DeviceKeystore: deviceKeystore,
		})
		require.NoError(t, err)
	}

	groupDatastore, err := cryptoutil.NewGroupDatastore(ds)
	require.NoError(t, err)

	serviceOpts := Opts{
		Host:           node.MockNode().PeerHost,
		PubSub:         node.PubSub(),
		Logger:         opts.Logger,
		RootDatastore:  ds,
		DeviceKeystore: deviceKeystore,
		IpfsCoreAPI:    node.API(),
		OrbitDB:        odb,
		TinderService:  node.Tinder(),
		PushKey:        opts.PushSK,
		GroupDatastore: groupDatastore,
	}

	service, cleanupService := TestingService(ctx, t, serviceOpts)

	// setup client
	grpcLogger := opts.Logger.Named("grpc")
	zapOpts := []grpc_zap.Option{}

	serverOpts := []grpc.ServerOption{
		grpc_middleware.WithUnaryServerChain(
			grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpc_zap.UnaryServerInterceptor(grpcLogger, zapOpts...),
		),
		grpc_middleware.WithStreamServerChain(
			grpc_ctxtags.StreamServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpc_zap.StreamServerInterceptor(grpcLogger, zapOpts...),
		),
	}

	clientOpts := []grpc.DialOption{
		grpc.WithChainUnaryInterceptor(),
		grpc.WithChainStreamInterceptor(),
	}

	server := grpc.NewServer(serverOpts...)
	client, cleanupClient := TestingClientFromServer(ctx, t, server, service, clientOpts...)

	tp := &TestingProtocol{
		Opts:    &serviceOpts,
		Client:  client,
		Service: service,

		RootDatastore:  ds,
		DeviceKeystore: deviceKeystore,
		IpfsCoreAPI:    node.API(),
		OrbitDB:        odb,
		GroupDatastore: serviceOpts.GroupDatastore,
	}
	cleanup := func() {
		server.Stop()
		cleanupClient()
		cleanupService()
		cleanupNode()
	}
	return tp, cleanup
}

func (opts *TestingOpts) applyDefaults(ctx context.Context, t testing.TB) {
	if opts.Logger == nil {
		opts.Logger = zap.NewNop()
	}

	if opts.Mocknet == nil {
		opts.Mocknet = libp2p_mocknet.New()
		t.Cleanup(func() { opts.Mocknet.Close() })
	}

	if opts.ConnectFunc == nil {
		opts.ConnectFunc = ConnectAll
	}
}

func NewTestingProtocolWithMockedPeers(ctx context.Context, t testing.TB, opts *TestingOpts, ds datastore.Batching, amount int) ([]*TestingProtocol, func()) {
	t.Helper()
	opts.applyDefaults(ctx, t)
	logger := opts.Logger

	if ds == nil {
		ds = ds_sync.MutexWrap(datastore.NewMapDatastore())
	}

	cleanupRDVP := func() {}
	closeRDVP := func() error { return nil }

	if opts.RDVPeer.ID == "" {
		rdvpeer, err := opts.Mocknet.GenPeer()
		require.NoError(t, err)
		require.NotNil(t, rdvpeer)

		_, cleanupRDVP = ipfsutil.TestingRDVP(ctx, t, rdvpeer)
		closeRDVP = rdvpeer.Close

		opts.RDVPeer = rdvpeer.Peerstore().PeerInfo(rdvpeer.ID())
	}

	rdvpnet := opts.Mocknet.Net(opts.RDVPeer.ID)
	require.NotNil(t, rdvpnet)

	cls := make([]func(), amount)
	tps := make([]*TestingProtocol, amount)
	for i := range tps {
		svcName := fmt.Sprintf("mock%d", i)
		opts.Logger = logger.Named(svcName)
		ds := datastoreutil.NewNamespacedDatastore(ds, datastore.NewKey(fmt.Sprintf("%d", i)))

		tps[i], cls[i] = NewTestingProtocol(ctx, t, opts, ds)
	}

	opts.ConnectFunc(t, opts.Mocknet)

	for _, net := range opts.Mocknet.Nets() {
		if net != rdvpnet {
			_, err := opts.Mocknet.LinkNets(net, rdvpnet)
			assert.NoError(t, err)

			_, err = opts.Mocknet.ConnectNets(net, rdvpnet)
			assert.NoError(t, err)
		}
	}

	cleanup := func() {
		for i := range cls {
			cls[i]()
		}

		cleanupRDVP()

		rdvpnet.Close()
		closeRDVP()
	}
	return tps, cleanup
}

// TestingService returns a configured Client struct with in-memory contexts.
func TestingService(ctx context.Context, t testing.TB, opts Opts) (Service, func()) {
	t.Helper()

	if opts.Logger == nil {
		opts.Logger = zap.NewNop()
	}

	cleanupNode := func() {}

	if opts.IpfsCoreAPI == nil {
		var mn ipfsutil.CoreAPIMock
		mn, cleanupNode = ipfsutil.TestingCoreAPI(ctx, t)
		opts.IpfsCoreAPI = mn.API()
	}

	service, err := New(opts)
	if err != nil {
		t.Fatalf("failed to initialize client: %v", err)
	}

	cleanup := func() {
		service.Close()
		cleanupNode()
	}

	return service, cleanup
}

func TestingClientFromServer(ctx context.Context, t testing.TB, s *grpc.Server, svc Service, dialOpts ...grpc.DialOption) (client Client, cleanup func()) {
	t.Helper()

	var err error

	client, err = NewClientFromServer(ctx, s, svc, dialOpts...)
	require.NoError(t, err)
	cleanup = func() {
		client.Close()
	}

	return
}

func TestingClient(ctx context.Context, t testing.TB, svc Service, clientOpts []grpc.DialOption, serverOpts []grpc.ServerOption) (client Client, cleanup func()) {
	t.Helper()

	var err error

	client, err = NewClient(ctx, svc, clientOpts, serverOpts)
	require.NoError(t, err)
	cleanup = func() { client.Close() }

	return
}

// Connect Peers Helper
type ConnectTestingProtocolFunc func(testing.TB, libp2p_mocknet.Mocknet)

// ConnectAll peers between themselves
func ConnectAll(t testing.TB, m libp2p_mocknet.Mocknet) {
	t.Helper()

	err := m.LinkAll()
	require.NoError(t, err)

	err = m.ConnectAllButSelf()
	require.NoError(t, err)
}

// ConnectInLine, connect peers one by one in order to make a straight line:
// ┌───┐    ┌───┐    ┌───┐         ┌───┐
// │ 1 │───▶│ 2 │───▶│ 3 │─ ─ ─ ─ ▶│ x │
// └───┘    └───┘    └───┘         └───┘

func ConnectInLine(t testing.TB, m libp2p_mocknet.Mocknet) {
	t.Helper()

	peers := m.Peers()

	for i := 0; i < len(peers)-1; i++ {
		_, err := m.LinkPeers(peers[i], peers[i+1])
		require.NoError(t, err)

		_, err = m.ConnectPeers(peers[i], peers[i+1])
		require.NoError(t, err)
	}
}

func CreatePeersWithGroupTest(ctx context.Context, t testing.TB, pathBase string, memberCount int, deviceCount int) ([]*mockedPeer, crypto.PrivKey, func()) {
	t.Helper()

	logger, cleanupLogger := testutil.Logger(t)

	var devKS cryptoutil.DeviceKeystore

	mockedPeers := make([]*mockedPeer, memberCount*deviceCount)

	g, groupSK, err := NewGroupMultiMember()
	if err != nil {
		t.Fatal(err)
	}

	mn := libp2p_mocknet.New()
	t.Cleanup(func() { mn.Close() })

	rdvp, err := mn.GenPeer()
	require.NoError(t, err, "failed to generate mocked peer")

	_, cleanuprdvp := ipfsutil.TestingRDVP(ctx, t, rdvp)

	ipfsopts := ipfsutil.TestingAPIOpts{
		Logger:  logger,
		Mocknet: mn,
		RDVPeer: rdvp.Peerstore().PeerInfo(rdvp.ID()),
	}
	deviceIndex := 0

	cls := make([]func(), memberCount)
	for i := 0; i < memberCount; i++ {
		for j := 0; j < deviceCount; j++ {
			ca, cleanupNode := ipfsutil.TestingCoreAPIUsingMockNet(ctx, t, &ipfsopts)

			if j == 0 {
				devKS = cryptoutil.NewDeviceKeystore(keystore.NewMemKeystore(), nil)
			} else {
				accSK, err := devKS.AccountPrivKey()
				require.NoError(t, err, "deviceKeystore private key")

				accProofSK, err := devKS.AccountProofPrivKey()
				require.NoError(t, err, "deviceKeystore private proof key")

				devKS, err = cryptoutil.NewWithExistingKeys(keystore.NewMemKeystore(), accSK, accProofSK)
				require.NoError(t, err, "deviceKeystore from existing keys")
			}

			mk, cleanupMessageKeystore := cryptoutil.NewInMemMessageKeystore()

			db, err := NewBertyOrbitDB(ctx, ca.API(), &NewOrbitDBOptions{
				NewOrbitDBOptions: orbitdb.NewOrbitDBOptions{
					Logger: logger,
				},
				DeviceKeystore:  devKS,
				MessageKeystore: mk,
			})
			if err != nil {
				t.Fatal(err)
			}

			gc, err := db.OpenGroup(ctx, g, nil)
			if err != nil {
				t.Fatalf("err: creating new group context, %v", err)
			}

			mp := &mockedPeer{
				CoreAPI: ca,
				DB:      db,
				GC:      gc,
				MKS:     mk,
				DevKS:   devKS,
			}

			// setup cleanup
			cls[i] = func() {
				if ms := mp.GC.MetadataStore(); ms != nil {
					err := ms.Drop()
					assert.NoError(t, err)
				}

				gc.Close()

				if db := mp.DB; db != nil {
					assert.NoError(t, err)

					err = db.Close()
					assert.NoError(t, err)
				}

				cleanupNode()
				cleanupMessageKeystore()
			}

			mockedPeers[deviceIndex] = mp
			deviceIndex++
		}
	}

	connectPeers(ctx, t, ipfsopts.Mocknet)

	return mockedPeers, groupSK, func() {
		for _, cleanup := range cls {
			cleanup()
		}

		cleanuprdvp()
		cleanupLogger()

		_ = rdvp.Close()
	}
}

func connectPeers(ctx context.Context, t testing.TB, mn libp2p_mocknet.Mocknet) {
	t.Helper()

	err := mn.LinkAll()
	require.NoError(t, err)

	err = mn.ConnectAllButSelf()
	require.NoError(t, err)
}

func dropPeers(t *testing.T, mockedPeers []*mockedPeer) {
	t.Helper()

	for _, m := range mockedPeers {
		if ms := m.GC.MetadataStore(); ms != nil {
			if err := ms.Drop(); err != nil {
				t.Fatal(err)
			}
		}

		if db := m.DB; db != nil {
			if err := db.Close(); err != nil {
				t.Fatal(err)
			}
		}

		if ca := m.CoreAPI; ca != nil {
			if err := ca.MockNode().Close(); err != nil {
				t.Fatal(err)
			}
		}
	}
}

type ServiceMethods interface {
	GetContextGroupForID(id []byte) (*GroupContext, error)
	GetCurrentDevicePushConfig() (*protocoltypes.PushServiceReceiver, *protocoltypes.PushServer)
}

func CreateVirtualOtherPeerSecretsShareSecret(t testing.TB, ctx context.Context, membersStores []*MetadataStore) (*cryptoutil.OwnMemberDevice, *protocoltypes.DeviceSecret) {
	// Manually adding another member to the group
	otherMemberSK, _, err := crypto.GenerateEd25519Key(crand.Reader)
	require.NoError(t, err)

	otherDeviceSK, _, err := crypto.GenerateEd25519Key(crand.Reader)
	require.NoError(t, err)

	otherMD := cryptoutil.NewOwnMemberDevice(otherMemberSK, otherDeviceSK)
	ds, err := cryptoutil.NewDeviceSecret()
	require.NoError(t, err)

	for _, store := range membersStores {
		_, err = MetadataStoreAddDeviceToGroup(ctx, store, store.Group(), otherMD)
		require.NoError(t, err)

		memPK, err := store.MemberPK()
		require.NoError(t, err)

		_, err = MetadataStoreSendSecret(ctx, store, store.Group(), otherMD, memPK, ds)
		require.NoError(t, err)
	}

	time.Sleep(time.Millisecond * 200)

	return otherMD, ds
}

func GetRootDatastoreForPath(dir string, key []byte, salt []byte, logger *zap.Logger) (datastore.Batching, error) {
	inMemory := dir == InMemoryDir

	var ds datastore.Batching
	if inMemory {
		ds = datastore.NewMapDatastore()
	} else {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return nil, errcode.TODO.Wrap(err)
		}

		dbPath := filepath.Join(dir, "datastore.sqlite")
		sqldsOpts := encrepo.SQLCipherDatastoreOptions{JournalMode: "WAL", PlaintextHeader: len(salt) != 0, Salt: salt}
		ds, err = encrepo.NewSQLCipherDatastore("sqlite3", dbPath, "blocks", key, sqldsOpts)
		if err != nil {
			return nil, errcode.TODO.Wrap(err)
		}
	}

	ds = ds_sync.MutexWrap(ds)

	return ds, nil
}

func CreateMultiMemberGroupInstance(ctx context.Context, t *testing.T, tps ...*TestingProtocol) *protocoltypes.Group {
	testutil.LogTree(t, "Create and Join MultiMember Group", 0, true)
	start := time.Now()

	ntps := len(tps)

	// Create group
	group, _, err := NewGroupMultiMember()
	require.NoError(t, err)

	// Get Instance Configurations
	{
		testutil.LogTree(t, "Get Instance Configuration", 1, true)
		start := time.Now()

		// check if everything is ready
		for _, pt := range tps {
			_, err := pt.Client.InstanceGetConfiguration(ctx, &protocoltypes.InstanceGetConfiguration_Request{})
			require.NoError(t, err)
		}

		testutil.LogTree(t, "duration: %s", 1, false, time.Since(start))
	}

	// Join Group
	{
		testutil.LogTree(t, "Join Group", 1, true)
		start := time.Now()

		for _, pt := range tps {
			req := protocoltypes.MultiMemberGroupJoin_Request{
				Group: group,
			}

			// pt join group
			_, err = pt.Client.MultiMemberGroupJoin(ctx, &req)
			require.NoError(t, err)
		}

		testutil.LogTree(t, "duration: %s", 1, false, time.Since(start))
	}

	// Get Member/Device PKs
	memberPKs := make([][]byte, ntps)
	devicePKs := make([][]byte, ntps)
	{
		testutil.LogTree(t, "Get Member/Device PKs", 1, true)
		start := time.Now()

		for i, pt := range tps {
			res, err := pt.Client.GroupInfo(ctx, &protocoltypes.GroupInfo_Request{
				GroupPK: group.PublicKey,
			})
			require.NoError(t, err)
			assert.Equal(t, group.PublicKey, res.Group.PublicKey)

			memberPKs[i] = res.MemberPK
			devicePKs[i] = res.DevicePK
		}

		testutil.LogTree(t, "duration: %s", 1, false, time.Since(start))
	}

	// Activate Group
	{
		testutil.LogTree(t, "Activate Group", 1, true)
		start := time.Now()

		for i, pt := range tps {
			_, err := pt.Client.ActivateGroup(ctx, &protocoltypes.ActivateGroup_Request{
				GroupPK: group.PublicKey,
			})

			assert.NoError(t, err, fmt.Sprintf("error for client %d", i))
		}

		testutil.LogTree(t, "duration: %s", 1, false, time.Since(start))
	}

	// Exchange Secrets
	{
		testutil.LogTree(t, "Exchange Secrets", 1, true)
		start := time.Now()

		wg := sync.WaitGroup{}
		secretsReceivedLock := sync.Mutex{}
		secretsReceived := make([]map[string]struct{}, ntps)
		wg.Add(ntps)

		nSuccess := int64(0)
		for i := range tps {
			go func(i int) {
				tp := tps[i]

				defer wg.Done()

				secretsReceived[i] = map[string]struct{}{}

				ctx, cancel := context.WithCancel(ctx)
				defer cancel()

				sub, inErr := tp.Client.GroupMetadataList(ctx, &protocoltypes.GroupMetadataList_Request{
					GroupPK: group.PublicKey,
				})
				if inErr != nil {
					assert.NoError(t, err, fmt.Sprintf("error for client %d", i))
					return
				}

				for {
					evt, inErr := sub.Recv()
					if inErr != nil {
						if inErr != io.EOF {
							assert.NoError(t, err, fmt.Sprintf("error for client %d", i))
						}

						break
					}

					if source, err := isEventAddSecretTargetedToMember(memberPKs[i], evt); err != nil {
						tps[i].Opts.Logger.Error("err:", zap.Error(inErr))
						assert.NoError(t, err, fmt.Sprintf("error for client %d", i))

						break
					} else if source != nil {
						secretsReceivedLock.Lock()
						secretsReceived[i][string(source)] = struct{}{}
						done := len(secretsReceived[i]) == ntps
						secretsReceivedLock.Unlock()

						if done {
							atomic.AddInt64(&nSuccess, 1)
							nSuccess := atomic.LoadInt64(&nSuccess)

							got := fmt.Sprintf("%d/%d", nSuccess, ntps)
							tps[i].Opts.Logger.Debug("received all secrets", zap.String("ok", got))
							return
						}
					}
				}
			}(i)
		}

		wg.Wait()

		secretsReceivedLock.Lock()
		ok := true
		for i := range secretsReceived {
			if !assert.Equal(t, ntps, len(secretsReceived[i]), fmt.Sprintf("mismatch for client %d", i)) {
				ok = false
			}
		}
		require.True(t, ok)
		secretsReceivedLock.Unlock()

		testutil.LogTree(t, "duration: %s", 1, false, time.Since(start))
	}

	testutil.LogTree(t, "duration: %s", 0, false, time.Since(start))

	return group
}

func isEventAddSecretTargetedToMember(ownRawPK []byte, evt *protocoltypes.GroupMetadataEvent) ([]byte, error) {
	// Only count EventTypeGroupDeviceSecretAdded events
	if evt.Metadata.EventType != protocoltypes.EventTypeGroupDeviceSecretAdded {
		return nil, nil
	}

	sec := &protocoltypes.GroupAddDeviceSecret{}
	err := sec.Unmarshal(evt.Event)
	if err != nil {
		return nil, err
	}

	// Filter out events targeted at other members
	if !bytes.Equal(ownRawPK, sec.DestMemberPK) {
		return nil, nil
	}

	return sec.DevicePK, nil
}
