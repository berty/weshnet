package protocoltypes

func (m *AccountGroupJoined) SetDevicePK(pk []byte) {
	m.DevicePK = pk
}

func (m *AccountGroupLeft) SetDevicePK(pk []byte) {
	m.DevicePK = pk
}

func (m *AccountContactRequestDisabled) SetDevicePK(pk []byte) {
	m.DevicePK = pk
}

func (m *AccountContactRequestEnabled) SetDevicePK(pk []byte) {
	m.DevicePK = pk
}

func (m *AccountContactRequestReferenceReset) SetDevicePK(pk []byte) {
	m.DevicePK = pk
}

func (m *AccountContactRequestOutgoingEnqueued) SetDevicePK(pk []byte) {
	m.DevicePK = pk
}

func (m *AccountContactRequestOutgoingSent) SetDevicePK(pk []byte) {
	m.DevicePK = pk
}

func (m *AccountContactRequestIncomingReceived) SetDevicePK(pk []byte) {
	m.DevicePK = pk
}

func (m *AccountContactRequestIncomingDiscarded) SetDevicePK(pk []byte) {
	m.DevicePK = pk
}

func (m *AccountContactRequestIncomingAccepted) SetDevicePK(pk []byte) {
	m.DevicePK = pk
}

func (m *AccountContactBlocked) SetDevicePK(pk []byte) {
	m.DevicePK = pk
}

func (m *AccountContactUnblocked) SetDevicePK(pk []byte) {
	m.DevicePK = pk
}

func (m *AccountContactRequestOutgoingSent) SetContactPK(pk []byte) {
	m.ContactPK = pk
}

func (m *AccountContactRequestIncomingDiscarded) SetContactPK(pk []byte) {
	m.ContactPK = pk
}

func (m *AccountContactRequestIncomingAccepted) SetContactPK(pk []byte) {
	m.ContactPK = pk
}

func (m *AccountContactBlocked) SetContactPK(pk []byte) {
	m.ContactPK = pk
}

func (m *AccountContactUnblocked) SetContactPK(pk []byte) {
	m.ContactPK = pk
}

func (m *AccountGroupLeft) SetGroupPK(pk []byte) {
	m.GroupPK = pk
}

func (m *ContactAliasKeyAdded) SetDevicePK(pk []byte) {
	m.DevicePK = pk
}

func (m *MultiMemberGroupAliasResolverAdded) SetDevicePK(pk []byte) {
	m.DevicePK = pk
}

func (m *MultiMemberGroupAdminRoleGranted) SetDevicePK(pk []byte) {
	m.DevicePK = pk
}

func (m *GroupMetadataPayloadSent) SetDevicePK(pk []byte) {
	m.DevicePK = pk
}

func (m *GroupReplicating) SetDevicePK(pk []byte) {
	m.DevicePK = pk
}

func (m *AccountVerifiedCredentialRegistered) SetDevicePK(pk []byte) {
	m.DevicePK = pk
}
