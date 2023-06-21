package weshnet

import (
	"testing"
	"time"

	mocknet "github.com/berty/go-libp2p-mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"berty.tech/go-ipfs-log/enc"
	"berty.tech/go-ipfs-log/entry"
	"berty.tech/go-orbit-db/iface"
	"berty.tech/weshnet/pkg/rendezvous"
	"berty.tech/weshnet/pkg/secretstore"
)

var (
	testSeed1 = []byte("secretsecretsecretsecretsecretse") // 32 bytes seed
	testSeed2 = []byte("badbadbadbadbadbadbadbadbadbadba") // 32 bytes seed
)

func TestRotationMessageMarshaler(t *testing.T) {
	key, err := enc.NewSecretbox(testSeed1)
	require.NoError(t, err)

	msg := &iface.MessageExchangeHeads{
		Address: "address_1",
		Heads:   []*entry.Entry{},
	}

	mn := mocknet.New()
	defer mn.Close()

	p, err := mn.GenPeer()
	require.NoError(t, err)

	// generate keystore
	acc1, err := secretstore.NewInMemSecretStore(nil)
	require.NoError(t, err)

	rp := rendezvous.NewStaticRotationInterval()
	m := NewOrbitDBMessageMarshaler(p.ID(), acc1, rp, false)

	g, _, err := NewGroupMultiMember()
	require.NoError(t, err)

	m.RegisterGroup(msg.Address, g)

	rp.RegisterRotation(time.Now(), msg.Address, testSeed1)
	m.RegisterSharedKeyForTopic(msg.Address, key)

	// marshal with register topic, should succeed
	payload, err := m.Marshal(msg)
	require.NoError(t, err)
	require.NotNil(t, payload)

	ret := iface.MessageExchangeHeads{}
	err = m.Unmarshal(payload, &ret)
	require.NoError(t, err)

	assert.Equal(t, msg.Address, ret.Address)
}

func TestRotationMessageMarshalUnknownTopic(t *testing.T) {
	mn := mocknet.New()
	defer mn.Close()

	msg := &iface.MessageExchangeHeads{
		Address: "address_1",
		Heads:   []*entry.Entry{},
	}

	p, err := mn.GenPeer()
	require.NoError(t, err)

	// generate keystore
	acc, err := secretstore.NewInMemSecretStore(nil)
	require.NoError(t, err)

	rp := rendezvous.NewStaticRotationInterval()
	m := NewOrbitDBMessageMarshaler(p.ID(), acc, rp, false)

	// marshal without register topic, should fail
	payload, err := m.Marshal(msg)
	require.Error(t, err)
	require.Nil(t, payload)
}

func TestRotationMessageUnmarshalUnknownTopic(t *testing.T) {
	mn := mocknet.New()
	defer mn.Close()

	msg := &iface.MessageExchangeHeads{
		Address: "address_1",
		Heads:   []*entry.Entry{},
	}
	key1, err := enc.NewSecretbox(testSeed1)
	require.NoError(t, err)

	key2, err := enc.NewSecretbox(testSeed2)
	require.NoError(t, err)

	p1, err := mn.GenPeer()
	require.NoError(t, err)

	p2, err := mn.GenPeer()
	require.NoError(t, err)

	// generate keystore
	acc1, err := secretstore.NewInMemSecretStore(nil)
	require.NoError(t, err)

	acc2, err := secretstore.NewInMemSecretStore(nil)
	require.NoError(t, err)

	g1, _, err := NewGroupMultiMember()
	require.NoError(t, err)

	g2, _, err := NewGroupMultiMember()
	require.NoError(t, err)

	rp1 := rendezvous.NewStaticRotationInterval()
	rp1.RegisterRotation(time.Now(), msg.Address, testSeed1)

	rp2 := rendezvous.NewStaticRotationInterval()

	m1 := NewOrbitDBMessageMarshaler(p1.ID(), acc1, rp1, false)
	m1.RegisterSharedKeyForTopic(msg.Address, key1)

	m1.RegisterGroup(msg.Address, g1)

	payload, err := m1.Marshal(msg)
	require.NoError(t, err)

	m2 := NewOrbitDBMessageMarshaler(p2.ID(), acc2, rp2, false)
	m2.RegisterSharedKeyForTopic(msg.Address, key2)

	m2.RegisterGroup(msg.Address, g2)

	var ret iface.MessageExchangeHeads

	// marshal with wrong key should fail
	err = m2.Unmarshal(payload, &ret)
	require.Error(t, err)
	assert.NotEqual(t, ret.Address, msg.Address)
}

func TestRotationMessageMarshalWrongKey(t *testing.T) {
	mn := mocknet.New()
	defer mn.Close()

	msg := &iface.MessageExchangeHeads{
		Address: "address_1",
		Heads:   []*entry.Entry{},
	}

	key1, err := enc.NewSecretbox(testSeed1)
	require.NoError(t, err)

	key2, err := enc.NewSecretbox(testSeed2)
	require.NoError(t, err)

	p1, err := mn.GenPeer()
	require.NoError(t, err)

	p2, err := mn.GenPeer()
	require.NoError(t, err)

	g1, _, err := NewGroupMultiMember()
	require.NoError(t, err)

	g2, _, err := NewGroupMultiMember()
	require.NoError(t, err)

	// generate keystore
	acc1, err := secretstore.NewInMemSecretStore(nil)
	require.NoError(t, err)

	acc2, err := secretstore.NewInMemSecretStore(nil)
	require.NoError(t, err)

	rp1 := rendezvous.NewStaticRotationInterval()
	rp1.RegisterRotation(time.Now(), msg.Address, testSeed1)
	rp2 := rendezvous.NewStaticRotationInterval()
	rp2.RegisterRotation(time.Now(), msg.Address, testSeed2)

	m1 := NewOrbitDBMessageMarshaler(p1.ID(), acc1, rp1, false)
	m1.RegisterSharedKeyForTopic(msg.Address, key1)
	m1.RegisterGroup(msg.Address, g1)

	payload, err := m1.Marshal(msg)
	require.NoError(t, err)

	m2 := NewOrbitDBMessageMarshaler(p2.ID(), acc2, rp2, false)
	m2.RegisterSharedKeyForTopic(msg.Address, key2)
	m2.RegisterGroup(msg.Address, g2)

	var ret iface.MessageExchangeHeads

	// marshal with wrong key should fail
	err = m2.Unmarshal(payload, &ret)
	require.Error(t, err)
	assert.NotEqual(t, ret.Address, msg.Address)

	// marshal with good key should succeed
	err = m1.Unmarshal(payload, &ret)
	require.NoError(t, err)
	assert.Equal(t, ret.Address, msg.Address)
}
