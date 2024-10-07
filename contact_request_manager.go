package weshnet

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
	"sync"
	"time"

	ipfscid "github.com/ipfs/go-cid"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/network"
	peer "github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/host/eventbus"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	"berty.tech/weshnet/v2/internal/handshake"
	"berty.tech/weshnet/v2/pkg/errcode"
	"berty.tech/weshnet/v2/pkg/ipfsutil"
	"berty.tech/weshnet/v2/pkg/logutil"
	"berty.tech/weshnet/v2/pkg/protocoltypes"
	"berty.tech/weshnet/v2/pkg/protoio"
	"berty.tech/weshnet/v2/pkg/tyber"
)

const contactRequestV1 = "/wesh/contact_req/1.0.0"

type contactRequestsManager struct {
	muManager sync.Mutex

	ctx            context.Context
	cancel         context.CancelFunc
	announceCancel context.CancelFunc

	lookupProcess   map[string]context.CancelFunc
	muLookupProcess sync.Mutex

	logger *zap.Logger

	enabled bool

	ownRendezvousSeed []byte
	accountPrivateKey crypto.PrivKey

	ipfs          ipfsutil.ExtendedCoreAPI
	swiper        *Swiper
	metadataStore *MetadataStore
}

func newContactRequestsManager(s *Swiper, store *MetadataStore, ipfs ipfsutil.ExtendedCoreAPI, logger *zap.Logger) (*contactRequestsManager, error) {
	accountPrivateKey, err := store.secretStore.GetAccountPrivateKey()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	cm := &contactRequestsManager{
		lookupProcess:     make(map[string]context.CancelFunc),
		metadataStore:     store,
		ipfs:              ipfs,
		logger:            logger.Named("req-mngr"),
		accountPrivateKey: accountPrivateKey,
		ctx:               ctx,
		cancel:            cancel,
		swiper:            s,
	}

	go cm.metadataWatcher(ctx)

	return cm, nil
}

func (c *contactRequestsManager) close() {
	if c.isClosed() {
		c.logger.Warn("contactRequestsManager already closed")
		return
	}

	c.cancel()

	c.muManager.Lock()
	defer c.muManager.Unlock()

	c.enabled = false

	c.disableAnnounce()

	c.ipfs.RemoveStreamHandler(contactRequestV1)
}

func (c *contactRequestsManager) isClosed() bool {
	select {
	case <-c.ctx.Done():
		return true
	default:
		return false
	}
}

func (c *contactRequestsManager) metadataWatcher(ctx context.Context) {
	handlers := map[protocoltypes.EventType]func(context.Context, *protocoltypes.GroupMetadataEvent) error{
		protocoltypes.EventType_EventTypeAccountContactRequestDisabled:         c.metadataRequestDisabled,
		protocoltypes.EventType_EventTypeAccountContactRequestEnabled:          c.metadataRequestEnabled,
		protocoltypes.EventType_EventTypeAccountContactRequestReferenceReset:   c.metadataRequestReset,
		protocoltypes.EventType_EventTypeAccountContactRequestOutgoingEnqueued: c.metadataRequestEnqueued,

		// @FIXME: looks like we don't need those events
		protocoltypes.EventType_EventTypeAccountContactRequestOutgoingSent:     c.metadataRequestSent,
		protocoltypes.EventType_EventTypeAccountContactRequestIncomingReceived: c.metadataRequestReceived,
	}

	// subscribe to new event
	sub, err := c.metadataStore.EventBus().Subscribe(new(*protocoltypes.GroupMetadataEvent),
		eventbus.Name("weshnet/rqmngr/metadata-watcher"))
	if err != nil {
		c.logger.Warn("unable to subscribe to group metadata event", zap.Error(err))
		return
	}

	// recreate previous contact request state
	enabled, contact := c.metadataStore.GetIncomingContactRequestsStatus()
	if contact != nil {
		c.ownRendezvousSeed = contact.PublicRendezvousSeed
	}

	c.muManager.Lock()
	if enabled {
		if err := c.enableContactRequest(ctx); err != nil {
			c.logger.Warn("unable to enable contact request", zap.Error(err))
		}
	}
	c.muManager.Unlock()

	// enqueue all contact with the `ToRequest` state
	for _, contact := range c.metadataStore.ListContactsByStatus(protocoltypes.ContactState_ContactStateToRequest) {
		if err := c.enqueueRequest(ctx, contact); err != nil {
			c.logger.Warn("unable to enqueue contact request", logutil.PrivateBinary("pk", contact.Pk), zap.Error(err))
		}
	}

	defer sub.Close()
	for {
		var evt interface{}
		select {
		case evt = <-sub.Out():
		case <-ctx.Done():
			return
		}

		// handle new events
		e := evt.(*protocoltypes.GroupMetadataEvent)
		typ := e.GetMetadata().GetEventType()
		hctx, _, endSection := tyber.Section(ctx, c.logger, fmt.Sprintf("handling event - %s", typ.String()))

		c.muManager.Lock()

		var err error
		if handler, ok := handlers[typ]; ok {
			if err = handler(hctx, e); err != nil {
				c.logger.Error("metadata store event handler", zap.String("event", typ.String()), zap.Error(err))
			}
		}

		c.muManager.Unlock()

		endSection(err, "")
	}
}

func (c *contactRequestsManager) metadataRequestDisabled(_ context.Context, _ *protocoltypes.GroupMetadataEvent) error {
	if !c.enabled {
		c.logger.Warn("contact request already disabled")
		return nil
	}

	c.enabled = false

	c.disableAnnounce()

	c.ipfs.RemoveStreamHandler(contactRequestV1)

	return nil
}

func (c *contactRequestsManager) metadataRequestEnabled(ctx context.Context, evt *protocoltypes.GroupMetadataEvent) error {
	e := &protocoltypes.AccountContactRequestEnabled{}
	if err := proto.Unmarshal(evt.Event, e); err != nil {
		return errcode.ErrCode_ErrDeserialization.Wrap(err)
	}

	return c.enableContactRequest(ctx)
}

func (c *contactRequestsManager) enableContactRequest(ctx context.Context) error {
	if c.enabled {
		c.logger.Warn("contact request already enabled")
		return nil
	}

	pkBytes, err := c.accountPrivateKey.GetPublic().Raw()
	if err != nil {
		return fmt.Errorf("unable to get raw pk: %w", err)
	}

	c.ipfs.SetStreamHandler(contactRequestV1, func(s network.Stream) {
		ctx, _, endSection := tyber.Section(c.ctx, c.logger, "receiving incoming contact request")

		if err := c.handleIncomingRequest(ctx, s); err != nil {
			c.logger.Error("unable to handle incoming contact request", zap.Error(err))
		}

		endSection(err, "")

		if err := s.Reset(); err != nil {
			c.logger.Error("unable to reset stream", zap.Error(err))
		}
	})

	c.enabled = true
	tyber.LogStep(ctx, c.logger, "enabled contact request")

	// announce on swiper if we already got seed
	if c.ownRendezvousSeed != nil {
		return c.enableAnnounce(ctx, c.ownRendezvousSeed, pkBytes)
	}

	c.logger.Warn("no seed registered, reset will be needed before announcing")
	return nil
}

func (c *contactRequestsManager) metadataRequestReset(ctx context.Context, evt *protocoltypes.GroupMetadataEvent) error {
	e := &protocoltypes.AccountContactRequestReferenceReset{}
	if err := proto.Unmarshal(evt.Event, e); err != nil {
		return errcode.ErrCode_ErrDeserialization.Wrap(err)
	}

	accPK, err := c.accountPrivateKey.GetPublic().Raw()
	if err != nil {
		return fmt.Errorf("unable to get raw pk: %w", err)
	}

	switch {
	case e.PublicRendezvousSeed == nil:
		return fmt.Errorf("unable to reset with an empty seed")
	case bytes.Equal(e.PublicRendezvousSeed, c.ownRendezvousSeed):
		return fmt.Errorf("unable to reset twice with the same seed")
	}

	// updating rendezvous seed
	tyber.LogStep(ctx, c.logger, "update rendezvous seed")
	c.ownRendezvousSeed = e.PublicRendezvousSeed

	// if contact request manager is disable don't run announce
	if !c.enabled {
		return nil
	}

	return c.enableAnnounce(ctx, c.ownRendezvousSeed, accPK)
}

func (c *contactRequestsManager) metadataRequestEnqueued(ctx context.Context, evt *protocoltypes.GroupMetadataEvent) error {
	ctx = tyber.ContextWithConstantTraceID(ctx, "msgrcvd-"+cidBytesString(evt.EventContext.Id))

	traceName := fmt.Sprintf("Received %s on group %s",
		strings.TrimPrefix(evt.Metadata.EventType.String(), "EventType"), base64.RawURLEncoding.EncodeToString(evt.EventContext.GroupPk))
	c.logger.Debug(traceName, tyber.FormatStepLogFields(ctx, []tyber.Detail{}, tyber.UpdateTraceName(traceName))...)

	e := &protocoltypes.AccountContactRequestOutgoingEnqueued{}
	if err := proto.Unmarshal(evt.Event, e); err != nil {
		return tyber.LogError(ctx, c.logger, "Failed to unmarshal event", err)
	}

	// enqueue contact request
	if err := c.enqueueRequest(ctx, e.Contact); err != nil {
		return tyber.LogError(ctx, c.logger, "Failed to enqueue request", err)
	}

	return nil
}

func (c *contactRequestsManager) metadataRequestSent(_ context.Context, evt *protocoltypes.GroupMetadataEvent) error {
	e := &protocoltypes.AccountContactRequestOutgoingSent{}
	if err := proto.Unmarshal(evt.Event, e); err != nil {
		return errcode.ErrCode_ErrDeserialization.Wrap(err)
	}

	// another device may have successfully sent contact request, try to cancel
	// lookup if needed
	c.cancelContactLookup(e.ContactPk)
	return nil
}

func (c *contactRequestsManager) metadataRequestReceived(_ context.Context, evt *protocoltypes.GroupMetadataEvent) error {
	e := &protocoltypes.AccountContactRequestIncomingReceived{}
	if err := proto.Unmarshal(evt.Event, e); err != nil {
		return errcode.ErrCode_ErrDeserialization.Wrap(err)
	}

	// another device may have successfully sent contact request, try to cancel
	// lookup if needed
	c.cancelContactLookup(e.ContactPk)
	return nil
}

func (c *contactRequestsManager) registerContactLookup(ctx context.Context, contactPK []byte) context.Context {
	c.muLookupProcess.Lock()

	key := hex.EncodeToString(contactPK)

	// make sure to only have one process for this pk running
	ctx, cancel := context.WithCancel(ctx)
	if cancelProvious, ok := c.lookupProcess[key]; ok {
		cancelProvious() // cancel previous lookup if needed
	}
	c.lookupProcess[key] = cancel

	c.muLookupProcess.Unlock()
	return ctx
}

func (c *contactRequestsManager) cancelContactLookup(contactPK []byte) {
	c.muLookupProcess.Lock()

	key := hex.EncodeToString(contactPK)

	// cancel current lookup if needed
	if cancel, ok := c.lookupProcess[key]; ok {
		cancel()
		delete(c.lookupProcess, key)
	}

	c.muLookupProcess.Unlock()
}

func (c *contactRequestsManager) enableAnnounce(ctx context.Context, seed, accPK []byte) error {
	if seed == nil {
		return fmt.Errorf("announcing with empty seed")
	}

	if c.announceCancel != nil { // is already enable
		tyber.LogStep(ctx, c.logger, "canceling previous announce")
		c.announceCancel()
	}

	ctx, c.announceCancel = context.WithCancel(ctx)
	c.enabled = true

	tyber.LogStep(ctx, c.logger, "announcing on swipper")

	// start announcing on swiper, this method should take care ton announce as
	// many time as needed
	c.swiper.Announce(ctx, accPK, seed)
	return nil
}

func (c *contactRequestsManager) disableAnnounce() {
	if c.announceCancel != nil {
		c.announceCancel()
		c.announceCancel = nil
	}
}

func (c *contactRequestsManager) enqueueRequest(ctx context.Context, to *protocoltypes.ShareableContact) (err error) {
	ctx, _, endSection := tyber.Section(ctx, c.logger, "Enqueue contact request: "+base64.RawURLEncoding.EncodeToString(to.Pk))

	otherPK, err := crypto.UnmarshalEd25519PublicKey(to.Pk)
	if err != nil {
		return err
	}

	if ok := c.metadataStore.checkContactStatus(otherPK, protocoltypes.ContactState_ContactStateAdded); ok {
		err = fmt.Errorf("contact already added")
		endSection(err, "")
		// contact already added,
		return err
	}

	// register lookup process
	ctx = c.registerContactLookup(ctx, to.Pk)

	// start watching topic on swiper, this method should take care of calling
	// `FindPeer` as many times as needed
	cpeers := c.swiper.WatchTopic(ctx, to.Pk, to.PublicRendezvousSeed)
	go func() {
		var err error
		for peer := range cpeers {
			// get our sharable contact to send to other contact
			if err = c.SendContactRequest(ctx, to, otherPK, peer); err != nil {
				c.logger.Warn("unable to send contact request", zap.Error(err))
			} else {
				// succefully send contact request, leave the loop and cancel lookup
				break
			}

			// wait one second to avoid infinity loop on send contact request
			// ex: when we dont have any network, send request can fail instantly
			time.Sleep(time.Second)
		}

		// cancel lookup process
		c.cancelContactLookup(to.Pk)

		endSection(err, "")
	}()

	return nil
}

// SendContactRequest try to perform contact request with the given remote peer
func (c *contactRequestsManager) SendContactRequest(ctx context.Context, to *protocoltypes.ShareableContact, otherPK crypto.PubKey, peer peer.AddrInfo) (err error) {
	ctx, _, endSection := tyber.Section(ctx, c.logger, "sending contact request")
	defer func() {
		endSection(err, "")
	}()

	_, own := c.metadataStore.GetIncomingContactRequestsStatus()
	if own == nil {
		err = fmt.Errorf("unable to retrieve own contact information")
		return err
	}

	// get own metadata for contact
	ownMetadata, err := c.metadataStore.GetRequestOwnMetadataForContact(to.Pk)
	if err != nil {
		c.logger.Warn("unable to get own metadata for contact", zap.Error(err))
		ownMetadata = nil
	}
	own.Metadata = ownMetadata

	// make sure to have connection with the remote peer
	if err := c.ipfs.Swarm().Connect(ctx, peer); err != nil {
		return fmt.Errorf("unable to connect: %w", err)
	}

	// create a new stream with the remote peer
	stream, err := c.ipfs.NewStream(network.WithAllowLimitedConn(ctx, "req_mngr"), peer.ID, contactRequestV1)
	if err != nil {
		return fmt.Errorf("unable to open stream: %w", err)
	}

	defer func() {
		if err := stream.Close(); err != nil {
			c.logger.Warn("error while closing stream with other peer", zap.Error(err))
		}
	}()

	reader := protoio.NewDelimitedReader(stream, 2048)
	writer := protoio.NewDelimitedWriter(stream)

	c.logger.Debug("performing handshake")

	tyber.LogStep(ctx, c.logger, "performing handshake")
	if err := handshake.RequestUsingReaderWriter(ctx, c.logger, reader, writer, c.accountPrivateKey, otherPK); err != nil {
		return fmt.Errorf("an error occurred during handshake: %w", err)
	}

	tyber.LogStep(ctx, c.logger, "sending own contact")
	// send own contact information
	if err := writer.WriteMsg(own); err != nil {
		return fmt.Errorf("an error occurred while sending own contact information: %w", err)
	}

	tyber.LogStep(ctx, c.logger, "mark contact request has sent")
	// mark this contact request as sent
	if _, err := c.metadataStore.ContactRequestOutgoingSent(ctx, otherPK); err != nil {
		return fmt.Errorf("an error occurred while marking contact request as sent: %w", err)
	}

	return nil
}

func (c *contactRequestsManager) handleIncomingRequest(ctx context.Context, stream network.Stream) (err error) {
	reader := protoio.NewDelimitedReader(stream, 2048)
	writer := protoio.NewDelimitedWriter(stream)

	tyber.LogStep(ctx, c.logger, "responding to handshake")

	otherPK, err := handshake.ResponseUsingReaderWriter(ctx, c.logger, reader, writer, c.accountPrivateKey)
	if err != nil {
		return fmt.Errorf("handshake failed: %w", err)
	}

	otherPKBytes, err := otherPK.Raw()
	if err != nil {
		return fmt.Errorf("failed to marshal contact public key: %w", err)
	}

	contact := &protocoltypes.ShareableContact{}

	tyber.LogStep(ctx, c.logger, "checking remote contact information")

	// read remote contact information
	if err := reader.ReadMsg(contact); err != nil {
		return fmt.Errorf("failed to read contact information: %w", err)
	}

	// validate contact pk
	if !bytes.Equal(otherPKBytes, contact.Pk) {
		return fmt.Errorf("contact information does not match handshake data")
	}

	// check contact information format
	if err := contact.CheckFormat(protocoltypes.ShareableContactOptionsAllowMissingRDVSeed); err != nil {
		return fmt.Errorf("invalid contact information format: %w", err)
	}

	tyber.LogStep(ctx, c.logger, "marking contact request has received")

	// mark contact request as received
	_, err = c.metadataStore.ContactRequestIncomingReceived(ctx, &protocoltypes.ShareableContact{
		Pk:                   otherPKBytes,
		PublicRendezvousSeed: contact.PublicRendezvousSeed,
		Metadata:             contact.Metadata,
	})
	if err != nil {
		return fmt.Errorf("invalid contact information format: %w", err)
	}

	return nil
}

func cidBytesString(bytes []byte) string {
	cid, err := ipfscid.Cast(bytes)
	if err != nil {
		return "error"
	}
	return cid.String()
}
