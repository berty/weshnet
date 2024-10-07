package weshnet

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/libp2p/go-libp2p/core/crypto"
	peer "github.com/libp2p/go-libp2p/core/peer"
	"google.golang.org/protobuf/proto"

	"berty.tech/go-ipfs-log/enc"
	"berty.tech/go-ipfs-log/entry"
	"berty.tech/go-orbit-db/iface"
	"berty.tech/weshnet/v2/pkg/protocoltypes"
	"berty.tech/weshnet/v2/pkg/rendezvous"
	"berty.tech/weshnet/v2/pkg/secretstore"
)

type PeerDeviceGroup struct {
	Group    *protocoltypes.Group
	DevicePK crypto.PubKey
}

type OrbitDBMessageMarshaler struct {
	rp           *rendezvous.RotationInterval
	sharedKeys   map[string]enc.SharedKey
	topicGroup   map[string]*protocoltypes.Group
	deviceCaches map[peer.ID]*PeerDeviceGroup
	muMarshall   sync.RWMutex
	selfid       peer.ID
	secretStore  secretstore.SecretStore

	// in Replication Mode DeviceKey should not be sent
	useReplicationMode bool
}

func NewOrbitDBMessageMarshaler(selfid peer.ID, secretStore secretstore.SecretStore, rp *rendezvous.RotationInterval, useReplicationMode bool) *OrbitDBMessageMarshaler {
	return &OrbitDBMessageMarshaler{
		selfid:             selfid,
		sharedKeys:         make(map[string]enc.SharedKey),
		deviceCaches:       make(map[peer.ID]*PeerDeviceGroup),
		topicGroup:         make(map[string]*protocoltypes.Group),
		rp:                 rp,
		secretStore:        secretStore,
		useReplicationMode: useReplicationMode,
	}
}

func (m *OrbitDBMessageMarshaler) RegisterSharedKeyForTopic(topic string, sk enc.SharedKey) {
	m.muMarshall.Lock()
	m.sharedKeys[topic] = sk
	m.muMarshall.Unlock()
}

func (m *OrbitDBMessageMarshaler) RegisterGroup(sid string, group *protocoltypes.Group) {
	m.muMarshall.Lock()
	m.topicGroup[sid] = group
	m.muMarshall.Unlock()
}

func (m *OrbitDBMessageMarshaler) GetDevicePKForPeerID(id peer.ID) (pdg *PeerDeviceGroup, ok bool) {
	m.muMarshall.RLock()
	pdg, ok = m.deviceCaches[id]
	m.muMarshall.RUnlock()
	return
}

func (m *OrbitDBMessageMarshaler) getSharedKeyFor(topic string) (sk enc.SharedKey, ok bool) {
	sk, ok = m.sharedKeys[topic]
	return
}

func (m *OrbitDBMessageMarshaler) Marshal(msg *iface.MessageExchangeHeads) ([]byte, error) {
	topic := msg.Address

	m.muMarshall.RLock()
	defer m.muMarshall.RUnlock()

	// marshall binary always return nil has error
	pid, _ := m.selfid.MarshalBinary()

	point, err := m.rp.PointForTopic(topic)
	if err != nil {
		return nil, fmt.Errorf("unable to get rendezvous for period: %w", err)
	}

	group, ok := m.topicGroup[topic]
	if !ok {
		return nil, fmt.Errorf("unknown group for topic: %s", topic)
	}

	var ownPK []byte

	// in replication mode, it doesn't make sense to send DevicePK
	if !m.useReplicationMode {
		ownDevice, err := m.secretStore.GetOwnMemberDeviceForGroup(group)
		if err != nil {
			return nil, fmt.Errorf("unable to get own member device key for group: %w", err)
		}

		ownPK, err = ownDevice.Device().Raw()
		if err != nil {
			return nil, fmt.Errorf("unable to get raw pk for device: %w", err)
		}
	}

	// @TODO(gfanton): use protobuf for this ?
	heads, err := json.Marshal(msg.Heads)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal heads: %w", err)
	}

	box := &protocoltypes.OrbitDBMessageHeads_Box{
		Address:  msg.Address,
		Heads:    heads,
		DevicePk: ownPK,
		PeerId:   pid,
	}

	sealedBox, err := m.sealBox(msg.Address, box)
	if err != nil {
		return nil, fmt.Errorf("unable to seal box: %w", err)
	}

	msghead := protocoltypes.OrbitDBMessageHeads{
		RawRotation: point.RawRotationTopic(),
		SealedBox:   sealedBox,
	}

	payload, err := proto.Marshal(&msghead)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal payload: %w", err)
	}

	return payload, nil
}

func (m *OrbitDBMessageMarshaler) Unmarshal(payload []byte, msg *iface.MessageExchangeHeads) error {
	m.muMarshall.Lock()
	defer m.muMarshall.Unlock()

	if msg == nil {
		msg = &iface.MessageExchangeHeads{}
	}

	msghead := protocoltypes.OrbitDBMessageHeads{}
	if err := proto.Unmarshal(payload, &msghead); err != nil {
		return fmt.Errorf("unable to unmarshal payload `%x`: %w", payload, err)
	}

	rotation := msghead.GetRawRotation()
	point, err := m.rp.PointForRawRotation(rotation)
	if err != nil {
		return fmt.Errorf("unable to get topic for rendezvous: %w", err)
	}

	box, err := m.openBox(point.Topic(), msghead.GetSealedBox())
	if err != nil {
		return fmt.Errorf("unable to open sealed box: %w", err)
	}

	var entries []*entry.Entry
	if err := json.Unmarshal(box.Heads, &entries); err != nil {
		return fmt.Errorf("unable to unmarshal entries: %w", err)
	}

	msg.Address = box.Address
	msg.Heads = entries

	if box.DevicePk == nil {
		// @NOTE(gfanton): this is probably a message from a replication server
		// which should not have a DevicePK
		return nil
	}

	pid, err := peer.IDFromBytes(box.PeerId)
	if err != nil {
		return fmt.Errorf("unable to parse peer id: %w", err)
	}

	// store device into cache
	var pdg PeerDeviceGroup

	pub, err := crypto.UnmarshalEd25519PublicKey(box.DevicePk)
	if err != nil {
		return fmt.Errorf("unable to unmarshal remote device pk: %w", err)
	}

	pdg.DevicePK = pub
	group, ok := m.topicGroup[msg.Address]
	if ok {
		// @FIXME(gfanton): do we need to raise an error here ?
		pdg.Group = group
	}
	m.deviceCaches[pid] = &pdg

	return nil
}

func (m *OrbitDBMessageMarshaler) sealBox(topic string, box *protocoltypes.OrbitDBMessageHeads_Box) ([]byte, error) {
	sk, ok := m.getSharedKeyFor(topic)
	if !ok {
		return nil, fmt.Errorf("unable to get shared key for topic")
	}

	rawBox, err := proto.Marshal(box)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal box %w", err)
	}

	sealedBox, err := sk.Seal(rawBox)
	if err != nil {
		return nil, fmt.Errorf("unable to seal box: %w", err)
	}

	return sealedBox, nil
}

func (m *OrbitDBMessageMarshaler) openBox(topic string, payload []byte) (*protocoltypes.OrbitDBMessageHeads_Box, error) {
	sk, ok := m.getSharedKeyFor(topic)
	if !ok {
		return nil, fmt.Errorf("unable to get shared key for topic")
	}

	rawBox, err := sk.Open(payload)
	if err != nil {
		return nil, fmt.Errorf("unable to open sealed box: %w", err)
	}

	box := &protocoltypes.OrbitDBMessageHeads_Box{}
	if err := proto.Unmarshal(rawBox, box); err != nil {
		return nil, fmt.Errorf("unable to unmarshal box: %w", err)
	}

	return box, nil
}
