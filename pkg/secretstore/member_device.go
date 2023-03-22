package secretstore

import (
	"github.com/libp2p/go-libp2p/core/crypto"
)

type ownMemberDevice struct {
	member crypto.PrivKey
	device crypto.PrivKey
	public *memberDevice
}

// newOwnMemberDevice instantiate a new ownMemberDevice allowing signing
// and encrypting data as both a device or a member part of a group.
// It also contains the public counterpart of the member and device keys.
func newOwnMemberDevice(member, device crypto.PrivKey) *ownMemberDevice {
	return &ownMemberDevice{
		member: member,
		device: device,
		public: newMemberDevice(member.GetPublic(), device.GetPublic()),
	}
}

func (d *ownMemberDevice) MemberSign(data []byte) ([]byte, error) {
	return d.member.Sign(data)
}

func (d *ownMemberDevice) DeviceSign(data []byte) ([]byte, error) {
	return d.device.Sign(data)
}

func (d *ownMemberDevice) Member() crypto.PubKey {
	return d.public.member
}

func (d *ownMemberDevice) Device() crypto.PubKey {
	return d.public.device
}

type memberDevice struct {
	member crypto.PubKey
	device crypto.PubKey
}

func NewMemberDevice(member, device crypto.PubKey) MemberDevice {
	return newMemberDevice(member, device)
}

func newMemberDevice(member, device crypto.PubKey) *memberDevice {
	return &memberDevice{
		member: member,
		device: device,
	}
}

func (m *memberDevice) Member() crypto.PubKey {
	return m.member
}

func (m *memberDevice) Device() crypto.PubKey {
	return m.device
}
