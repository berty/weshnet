# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [protocoltypes.proto](#protocoltypes-proto)
    - [Account](#weshnet-protocol-v1-Account)
    - [AccountContactBlocked](#weshnet-protocol-v1-AccountContactBlocked)
    - [AccountContactRequestDisabled](#weshnet-protocol-v1-AccountContactRequestDisabled)
    - [AccountContactRequestEnabled](#weshnet-protocol-v1-AccountContactRequestEnabled)
    - [AccountContactRequestIncomingAccepted](#weshnet-protocol-v1-AccountContactRequestIncomingAccepted)
    - [AccountContactRequestIncomingDiscarded](#weshnet-protocol-v1-AccountContactRequestIncomingDiscarded)
    - [AccountContactRequestIncomingReceived](#weshnet-protocol-v1-AccountContactRequestIncomingReceived)
    - [AccountContactRequestOutgoingEnqueued](#weshnet-protocol-v1-AccountContactRequestOutgoingEnqueued)
    - [AccountContactRequestOutgoingSent](#weshnet-protocol-v1-AccountContactRequestOutgoingSent)
    - [AccountContactRequestReferenceReset](#weshnet-protocol-v1-AccountContactRequestReferenceReset)
    - [AccountContactUnblocked](#weshnet-protocol-v1-AccountContactUnblocked)
    - [AccountGroupJoined](#weshnet-protocol-v1-AccountGroupJoined)
    - [AccountGroupLeft](#weshnet-protocol-v1-AccountGroupLeft)
    - [AccountServiceTokenAdded](#weshnet-protocol-v1-AccountServiceTokenAdded)
    - [AccountServiceTokenRemoved](#weshnet-protocol-v1-AccountServiceTokenRemoved)
    - [AccountVerifiedCredentialRegistered](#weshnet-protocol-v1-AccountVerifiedCredentialRegistered)
    - [ActivateGroup](#weshnet-protocol-v1-ActivateGroup)
    - [ActivateGroup.Reply](#weshnet-protocol-v1-ActivateGroup-Reply)
    - [ActivateGroup.Request](#weshnet-protocol-v1-ActivateGroup-Request)
    - [AppMessageSend](#weshnet-protocol-v1-AppMessageSend)
    - [AppMessageSend.Reply](#weshnet-protocol-v1-AppMessageSend-Reply)
    - [AppMessageSend.Request](#weshnet-protocol-v1-AppMessageSend-Request)
    - [AppMetadataSend](#weshnet-protocol-v1-AppMetadataSend)
    - [AppMetadataSend.Reply](#weshnet-protocol-v1-AppMetadataSend-Reply)
    - [AppMetadataSend.Request](#weshnet-protocol-v1-AppMetadataSend-Request)
    - [AuthExchangeResponse](#weshnet-protocol-v1-AuthExchangeResponse)
    - [AuthExchangeResponse.ServicesEntry](#weshnet-protocol-v1-AuthExchangeResponse-ServicesEntry)
    - [AuthServiceCompleteFlow](#weshnet-protocol-v1-AuthServiceCompleteFlow)
    - [AuthServiceCompleteFlow.Reply](#weshnet-protocol-v1-AuthServiceCompleteFlow-Reply)
    - [AuthServiceCompleteFlow.Request](#weshnet-protocol-v1-AuthServiceCompleteFlow-Request)
    - [AuthServiceInitFlow](#weshnet-protocol-v1-AuthServiceInitFlow)
    - [AuthServiceInitFlow.Reply](#weshnet-protocol-v1-AuthServiceInitFlow-Reply)
    - [AuthServiceInitFlow.Request](#weshnet-protocol-v1-AuthServiceInitFlow-Request)
    - [ContactAliasKeyAdded](#weshnet-protocol-v1-ContactAliasKeyAdded)
    - [ContactAliasKeySend](#weshnet-protocol-v1-ContactAliasKeySend)
    - [ContactAliasKeySend.Reply](#weshnet-protocol-v1-ContactAliasKeySend-Reply)
    - [ContactAliasKeySend.Request](#weshnet-protocol-v1-ContactAliasKeySend-Request)
    - [ContactBlock](#weshnet-protocol-v1-ContactBlock)
    - [ContactBlock.Reply](#weshnet-protocol-v1-ContactBlock-Reply)
    - [ContactBlock.Request](#weshnet-protocol-v1-ContactBlock-Request)
    - [ContactRequestAccept](#weshnet-protocol-v1-ContactRequestAccept)
    - [ContactRequestAccept.Reply](#weshnet-protocol-v1-ContactRequestAccept-Reply)
    - [ContactRequestAccept.Request](#weshnet-protocol-v1-ContactRequestAccept-Request)
    - [ContactRequestDisable](#weshnet-protocol-v1-ContactRequestDisable)
    - [ContactRequestDisable.Reply](#weshnet-protocol-v1-ContactRequestDisable-Reply)
    - [ContactRequestDisable.Request](#weshnet-protocol-v1-ContactRequestDisable-Request)
    - [ContactRequestDiscard](#weshnet-protocol-v1-ContactRequestDiscard)
    - [ContactRequestDiscard.Reply](#weshnet-protocol-v1-ContactRequestDiscard-Reply)
    - [ContactRequestDiscard.Request](#weshnet-protocol-v1-ContactRequestDiscard-Request)
    - [ContactRequestEnable](#weshnet-protocol-v1-ContactRequestEnable)
    - [ContactRequestEnable.Reply](#weshnet-protocol-v1-ContactRequestEnable-Reply)
    - [ContactRequestEnable.Request](#weshnet-protocol-v1-ContactRequestEnable-Request)
    - [ContactRequestReference](#weshnet-protocol-v1-ContactRequestReference)
    - [ContactRequestReference.Reply](#weshnet-protocol-v1-ContactRequestReference-Reply)
    - [ContactRequestReference.Request](#weshnet-protocol-v1-ContactRequestReference-Request)
    - [ContactRequestResetReference](#weshnet-protocol-v1-ContactRequestResetReference)
    - [ContactRequestResetReference.Reply](#weshnet-protocol-v1-ContactRequestResetReference-Reply)
    - [ContactRequestResetReference.Request](#weshnet-protocol-v1-ContactRequestResetReference-Request)
    - [ContactRequestSend](#weshnet-protocol-v1-ContactRequestSend)
    - [ContactRequestSend.Reply](#weshnet-protocol-v1-ContactRequestSend-Reply)
    - [ContactRequestSend.Request](#weshnet-protocol-v1-ContactRequestSend-Request)
    - [ContactUnblock](#weshnet-protocol-v1-ContactUnblock)
    - [ContactUnblock.Reply](#weshnet-protocol-v1-ContactUnblock-Reply)
    - [ContactUnblock.Request](#weshnet-protocol-v1-ContactUnblock-Request)
    - [CredentialVerificationServiceCompleteFlow](#weshnet-protocol-v1-CredentialVerificationServiceCompleteFlow)
    - [CredentialVerificationServiceCompleteFlow.Reply](#weshnet-protocol-v1-CredentialVerificationServiceCompleteFlow-Reply)
    - [CredentialVerificationServiceCompleteFlow.Request](#weshnet-protocol-v1-CredentialVerificationServiceCompleteFlow-Request)
    - [CredentialVerificationServiceInitFlow](#weshnet-protocol-v1-CredentialVerificationServiceInitFlow)
    - [CredentialVerificationServiceInitFlow.Reply](#weshnet-protocol-v1-CredentialVerificationServiceInitFlow-Reply)
    - [CredentialVerificationServiceInitFlow.Request](#weshnet-protocol-v1-CredentialVerificationServiceInitFlow-Request)
    - [DeactivateGroup](#weshnet-protocol-v1-DeactivateGroup)
    - [DeactivateGroup.Reply](#weshnet-protocol-v1-DeactivateGroup-Reply)
    - [DeactivateGroup.Request](#weshnet-protocol-v1-DeactivateGroup-Request)
    - [DebugAuthServiceSetToken](#weshnet-protocol-v1-DebugAuthServiceSetToken)
    - [DebugAuthServiceSetToken.Reply](#weshnet-protocol-v1-DebugAuthServiceSetToken-Reply)
    - [DebugAuthServiceSetToken.Request](#weshnet-protocol-v1-DebugAuthServiceSetToken-Request)
    - [DebugGroup](#weshnet-protocol-v1-DebugGroup)
    - [DebugGroup.Reply](#weshnet-protocol-v1-DebugGroup-Reply)
    - [DebugGroup.Request](#weshnet-protocol-v1-DebugGroup-Request)
    - [DebugInspectGroupStore](#weshnet-protocol-v1-DebugInspectGroupStore)
    - [DebugInspectGroupStore.Reply](#weshnet-protocol-v1-DebugInspectGroupStore-Reply)
    - [DebugInspectGroupStore.Request](#weshnet-protocol-v1-DebugInspectGroupStore-Request)
    - [DebugListGroups](#weshnet-protocol-v1-DebugListGroups)
    - [DebugListGroups.Reply](#weshnet-protocol-v1-DebugListGroups-Reply)
    - [DebugListGroups.Request](#weshnet-protocol-v1-DebugListGroups-Request)
    - [DecodeContact](#weshnet-protocol-v1-DecodeContact)
    - [DecodeContact.Reply](#weshnet-protocol-v1-DecodeContact-Reply)
    - [DecodeContact.Request](#weshnet-protocol-v1-DecodeContact-Request)
    - [DeviceChainKey](#weshnet-protocol-v1-DeviceChainKey)
    - [EncryptedMessage](#weshnet-protocol-v1-EncryptedMessage)
    - [EventContext](#weshnet-protocol-v1-EventContext)
    - [FirstLastCounters](#weshnet-protocol-v1-FirstLastCounters)
    - [Group](#weshnet-protocol-v1-Group)
    - [GroupAddAdditionalRendezvousSeed](#weshnet-protocol-v1-GroupAddAdditionalRendezvousSeed)
    - [GroupDeviceChainKeyAdded](#weshnet-protocol-v1-GroupDeviceChainKeyAdded)
    - [GroupDeviceStatus](#weshnet-protocol-v1-GroupDeviceStatus)
    - [GroupDeviceStatus.Reply](#weshnet-protocol-v1-GroupDeviceStatus-Reply)
    - [GroupDeviceStatus.Reply.PeerConnected](#weshnet-protocol-v1-GroupDeviceStatus-Reply-PeerConnected)
    - [GroupDeviceStatus.Reply.PeerDisconnected](#weshnet-protocol-v1-GroupDeviceStatus-Reply-PeerDisconnected)
    - [GroupDeviceStatus.Reply.PeerReconnecting](#weshnet-protocol-v1-GroupDeviceStatus-Reply-PeerReconnecting)
    - [GroupDeviceStatus.Request](#weshnet-protocol-v1-GroupDeviceStatus-Request)
    - [GroupEnvelope](#weshnet-protocol-v1-GroupEnvelope)
    - [GroupHeadsExport](#weshnet-protocol-v1-GroupHeadsExport)
    - [GroupInfo](#weshnet-protocol-v1-GroupInfo)
    - [GroupInfo.Reply](#weshnet-protocol-v1-GroupInfo-Reply)
    - [GroupInfo.Request](#weshnet-protocol-v1-GroupInfo-Request)
    - [GroupMemberDeviceAdded](#weshnet-protocol-v1-GroupMemberDeviceAdded)
    - [GroupMessageEvent](#weshnet-protocol-v1-GroupMessageEvent)
    - [GroupMessageList](#weshnet-protocol-v1-GroupMessageList)
    - [GroupMessageList.Request](#weshnet-protocol-v1-GroupMessageList-Request)
    - [GroupMetadata](#weshnet-protocol-v1-GroupMetadata)
    - [GroupMetadataEvent](#weshnet-protocol-v1-GroupMetadataEvent)
    - [GroupMetadataList](#weshnet-protocol-v1-GroupMetadataList)
    - [GroupMetadataList.Request](#weshnet-protocol-v1-GroupMetadataList-Request)
    - [GroupMetadataPayloadSent](#weshnet-protocol-v1-GroupMetadataPayloadSent)
    - [GroupRemoveAdditionalRendezvousSeed](#weshnet-protocol-v1-GroupRemoveAdditionalRendezvousSeed)
    - [GroupReplicating](#weshnet-protocol-v1-GroupReplicating)
    - [MemberWithDevices](#weshnet-protocol-v1-MemberWithDevices)
    - [MessageEnvelope](#weshnet-protocol-v1-MessageEnvelope)
    - [MessageHeaders](#weshnet-protocol-v1-MessageHeaders)
    - [MessageHeaders.MetadataEntry](#weshnet-protocol-v1-MessageHeaders-MetadataEntry)
    - [MultiMemberGroupAdminRoleGrant](#weshnet-protocol-v1-MultiMemberGroupAdminRoleGrant)
    - [MultiMemberGroupAdminRoleGrant.Reply](#weshnet-protocol-v1-MultiMemberGroupAdminRoleGrant-Reply)
    - [MultiMemberGroupAdminRoleGrant.Request](#weshnet-protocol-v1-MultiMemberGroupAdminRoleGrant-Request)
    - [MultiMemberGroupAdminRoleGranted](#weshnet-protocol-v1-MultiMemberGroupAdminRoleGranted)
    - [MultiMemberGroupAliasResolverAdded](#weshnet-protocol-v1-MultiMemberGroupAliasResolverAdded)
    - [MultiMemberGroupAliasResolverDisclose](#weshnet-protocol-v1-MultiMemberGroupAliasResolverDisclose)
    - [MultiMemberGroupAliasResolverDisclose.Reply](#weshnet-protocol-v1-MultiMemberGroupAliasResolverDisclose-Reply)
    - [MultiMemberGroupAliasResolverDisclose.Request](#weshnet-protocol-v1-MultiMemberGroupAliasResolverDisclose-Request)
    - [MultiMemberGroupCreate](#weshnet-protocol-v1-MultiMemberGroupCreate)
    - [MultiMemberGroupCreate.Reply](#weshnet-protocol-v1-MultiMemberGroupCreate-Reply)
    - [MultiMemberGroupCreate.Request](#weshnet-protocol-v1-MultiMemberGroupCreate-Request)
    - [MultiMemberGroupInitialMemberAnnounced](#weshnet-protocol-v1-MultiMemberGroupInitialMemberAnnounced)
    - [MultiMemberGroupInvitationCreate](#weshnet-protocol-v1-MultiMemberGroupInvitationCreate)
    - [MultiMemberGroupInvitationCreate.Reply](#weshnet-protocol-v1-MultiMemberGroupInvitationCreate-Reply)
    - [MultiMemberGroupInvitationCreate.Request](#weshnet-protocol-v1-MultiMemberGroupInvitationCreate-Request)
    - [MultiMemberGroupJoin](#weshnet-protocol-v1-MultiMemberGroupJoin)
    - [MultiMemberGroupJoin.Reply](#weshnet-protocol-v1-MultiMemberGroupJoin-Reply)
    - [MultiMemberGroupJoin.Request](#weshnet-protocol-v1-MultiMemberGroupJoin-Request)
    - [MultiMemberGroupLeave](#weshnet-protocol-v1-MultiMemberGroupLeave)
    - [MultiMemberGroupLeave.Reply](#weshnet-protocol-v1-MultiMemberGroupLeave-Reply)
    - [MultiMemberGroupLeave.Request](#weshnet-protocol-v1-MultiMemberGroupLeave-Request)
    - [OrbitDBMessageHeads](#weshnet-protocol-v1-OrbitDBMessageHeads)
    - [OrbitDBMessageHeads.Box](#weshnet-protocol-v1-OrbitDBMessageHeads-Box)
    - [OutOfStoreMessage](#weshnet-protocol-v1-OutOfStoreMessage)
    - [OutOfStoreReceive](#weshnet-protocol-v1-OutOfStoreReceive)
    - [OutOfStoreReceive.Reply](#weshnet-protocol-v1-OutOfStoreReceive-Reply)
    - [OutOfStoreReceive.Request](#weshnet-protocol-v1-OutOfStoreReceive-Request)
    - [OutOfStoreSeal](#weshnet-protocol-v1-OutOfStoreSeal)
    - [OutOfStoreSeal.Reply](#weshnet-protocol-v1-OutOfStoreSeal-Reply)
    - [OutOfStoreSeal.Request](#weshnet-protocol-v1-OutOfStoreSeal-Request)
    - [PeerList](#weshnet-protocol-v1-PeerList)
    - [PeerList.Peer](#weshnet-protocol-v1-PeerList-Peer)
    - [PeerList.Reply](#weshnet-protocol-v1-PeerList-Reply)
    - [PeerList.Request](#weshnet-protocol-v1-PeerList-Request)
    - [PeerList.Route](#weshnet-protocol-v1-PeerList-Route)
    - [PeerList.Stream](#weshnet-protocol-v1-PeerList-Stream)
    - [Progress](#weshnet-protocol-v1-Progress)
    - [ProtocolMetadata](#weshnet-protocol-v1-ProtocolMetadata)
    - [PushDeviceServerRegistered](#weshnet-protocol-v1-PushDeviceServerRegistered)
    - [PushDeviceTokenRegistered](#weshnet-protocol-v1-PushDeviceTokenRegistered)
    - [PushMemberTokenUpdate](#weshnet-protocol-v1-PushMemberTokenUpdate)
    - [PushServer](#weshnet-protocol-v1-PushServer)
    - [PushServiceReceiver](#weshnet-protocol-v1-PushServiceReceiver)
    - [RefreshContactRequest](#weshnet-protocol-v1-RefreshContactRequest)
    - [RefreshContactRequest.Peer](#weshnet-protocol-v1-RefreshContactRequest-Peer)
    - [RefreshContactRequest.Reply](#weshnet-protocol-v1-RefreshContactRequest-Reply)
    - [RefreshContactRequest.Request](#weshnet-protocol-v1-RefreshContactRequest-Request)
    - [ReplicationServiceRegisterGroup](#weshnet-protocol-v1-ReplicationServiceRegisterGroup)
    - [ReplicationServiceRegisterGroup.Reply](#weshnet-protocol-v1-ReplicationServiceRegisterGroup-Reply)
    - [ReplicationServiceRegisterGroup.Request](#weshnet-protocol-v1-ReplicationServiceRegisterGroup-Request)
    - [ReplicationServiceReplicateGroup](#weshnet-protocol-v1-ReplicationServiceReplicateGroup)
    - [ReplicationServiceReplicateGroup.Reply](#weshnet-protocol-v1-ReplicationServiceReplicateGroup-Reply)
    - [ReplicationServiceReplicateGroup.Request](#weshnet-protocol-v1-ReplicationServiceReplicateGroup-Request)
    - [ServiceExportData](#weshnet-protocol-v1-ServiceExportData)
    - [ServiceExportData.Reply](#weshnet-protocol-v1-ServiceExportData-Reply)
    - [ServiceExportData.Request](#weshnet-protocol-v1-ServiceExportData-Request)
    - [ServiceGetConfiguration](#weshnet-protocol-v1-ServiceGetConfiguration)
    - [ServiceGetConfiguration.Reply](#weshnet-protocol-v1-ServiceGetConfiguration-Reply)
    - [ServiceGetConfiguration.Request](#weshnet-protocol-v1-ServiceGetConfiguration-Request)
    - [ServiceToken](#weshnet-protocol-v1-ServiceToken)
    - [ServiceTokenSupportedService](#weshnet-protocol-v1-ServiceTokenSupportedService)
    - [ServicesTokenCode](#weshnet-protocol-v1-ServicesTokenCode)
    - [ServicesTokenList](#weshnet-protocol-v1-ServicesTokenList)
    - [ServicesTokenList.Reply](#weshnet-protocol-v1-ServicesTokenList-Reply)
    - [ServicesTokenList.Request](#weshnet-protocol-v1-ServicesTokenList-Request)
    - [ShareContact](#weshnet-protocol-v1-ShareContact)
    - [ShareContact.Reply](#weshnet-protocol-v1-ShareContact-Reply)
    - [ShareContact.Request](#weshnet-protocol-v1-ShareContact-Request)
    - [ShareableContact](#weshnet-protocol-v1-ShareableContact)
    - [SystemInfo](#weshnet-protocol-v1-SystemInfo)
    - [SystemInfo.OrbitDB](#weshnet-protocol-v1-SystemInfo-OrbitDB)
    - [SystemInfo.OrbitDB.ReplicationStatus](#weshnet-protocol-v1-SystemInfo-OrbitDB-ReplicationStatus)
    - [SystemInfo.P2P](#weshnet-protocol-v1-SystemInfo-P2P)
    - [SystemInfo.Process](#weshnet-protocol-v1-SystemInfo-Process)
    - [SystemInfo.Reply](#weshnet-protocol-v1-SystemInfo-Reply)
    - [SystemInfo.Request](#weshnet-protocol-v1-SystemInfo-Request)
    - [VerifiedCredentialsList](#weshnet-protocol-v1-VerifiedCredentialsList)
    - [VerifiedCredentialsList.Reply](#weshnet-protocol-v1-VerifiedCredentialsList-Reply)
    - [VerifiedCredentialsList.Request](#weshnet-protocol-v1-VerifiedCredentialsList-Request)
  
    - [ContactState](#weshnet-protocol-v1-ContactState)
    - [DebugInspectGroupLogType](#weshnet-protocol-v1-DebugInspectGroupLogType)
    - [Direction](#weshnet-protocol-v1-Direction)
    - [EventType](#weshnet-protocol-v1-EventType)
    - [GroupDeviceStatus.Transport](#weshnet-protocol-v1-GroupDeviceStatus-Transport)
    - [GroupDeviceStatus.Type](#weshnet-protocol-v1-GroupDeviceStatus-Type)
    - [GroupType](#weshnet-protocol-v1-GroupType)
    - [PeerList.Feature](#weshnet-protocol-v1-PeerList-Feature)
    - [ServiceGetConfiguration.SettingState](#weshnet-protocol-v1-ServiceGetConfiguration-SettingState)
  
    - [ProtocolService](#weshnet-protocol-v1-ProtocolService)
  
- [Scalar Value Types](#scalar-value-types)

<a name="protocoltypes-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## protocoltypes.proto

<a name="weshnet-protocol-v1-Account"></a>

### Account
Account describes all the secrets that identifies an Account

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| group | [Group](#weshnet-protocol-v1-Group) |  | group specifies which group is used to manage the account |
| account_private_key | [bytes](#bytes) |  | account_private_key, private part is used to signs handshake, signs device, create contacts group keys via ECDH -- public part is used to have a shareable identity |
| alias_private_key | [bytes](#bytes) |  | alias_private_key, private part is use to derive group members private keys, signs alias proofs, public part can be shared to contacts to prove identity |
| public_rendezvous_seed | [bytes](#bytes) |  | public_rendezvous_seed, rendezvous seed used for direct communication |

<a name="weshnet-protocol-v1-AccountContactBlocked"></a>

### AccountContactBlocked
AccountContactBlocked indicates that a contact is blocked

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| device_pk | [bytes](#bytes) |  | device_pk is the device sending the event, signs the message |
| contact_pk | [bytes](#bytes) |  | contact_pk is the contact blocked |

<a name="weshnet-protocol-v1-AccountContactRequestDisabled"></a>

### AccountContactRequestDisabled
AccountContactRequestDisabled indicates that the account should not be advertised on a public rendezvous point

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| device_pk | [bytes](#bytes) |  | device_pk is the device sending the event, signs the message |

<a name="weshnet-protocol-v1-AccountContactRequestEnabled"></a>

### AccountContactRequestEnabled
AccountContactRequestEnabled indicates that the account should be advertised on a public rendezvous point

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| device_pk | [bytes](#bytes) |  | device_pk is the device sending the event, signs the message |

<a name="weshnet-protocol-v1-AccountContactRequestIncomingAccepted"></a>

### AccountContactRequestIncomingAccepted
This event should be followed by an AccountGroupJoined event
This event should be followed by GroupMemberDeviceAdded and GroupDeviceChainKeyAdded events within the AccountGroup
AccountContactRequestIncomingAccepted indicates that a contact request has been accepted

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| device_pk | [bytes](#bytes) |  | device_pk is the device sending the event, signs the message |
| contact_pk | [bytes](#bytes) |  | contact_pk is the contact whom request is accepted |
| group_pk | [bytes](#bytes) |  | group_pk is the 1to1 group with the requester user |

<a name="weshnet-protocol-v1-AccountContactRequestIncomingDiscarded"></a>

### AccountContactRequestIncomingDiscarded
AccountContactRequestIncomingDiscarded indicates that a contact request has been refused

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| device_pk | [bytes](#bytes) |  | device_pk is the device sending the event, signs the message |
| contact_pk | [bytes](#bytes) |  | contact_pk is the contact whom request is refused |

<a name="weshnet-protocol-v1-AccountContactRequestIncomingReceived"></a>

### AccountContactRequestIncomingReceived
AccountContactRequestIncomingReceived indicates that the account has received a new contact request

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| device_pk | [bytes](#bytes) |  | device_pk is the device sending the account event (which received the contact request), signs the message |
| contact_pk | [bytes](#bytes) |  | contact_pk is the account sending the request |
| contact_rendezvous_seed | [bytes](#bytes) |  | TODO: is this necessary? contact_rendezvous_seed is the rendezvous seed of the contact sending the request |
| contact_metadata | [bytes](#bytes) |  | TODO: is this necessary? contact_metadata is the metadata specific to the app to identify the contact for the request |

<a name="weshnet-protocol-v1-AccountContactRequestOutgoingEnqueued"></a>

### AccountContactRequestOutgoingEnqueued
This event should be followed by an AccountGroupJoined event
This event should be followed by a GroupMemberDeviceAdded event within the AccountGroup
This event should be followed by a GroupDeviceChainKeyAdded event within the AccountGroup
AccountContactRequestOutgoingEnqueued indicates that the account will attempt to send a contact request when a matching peer is discovered

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| device_pk | [bytes](#bytes) |  | device_pk is the device sending the event, signs the message |
| group_pk | [bytes](#bytes) |  | group_pk is the 1to1 group with the requested user |
| contact | [ShareableContact](#weshnet-protocol-v1-ShareableContact) |  | contact is a message describing how to connect to the other account |
| own_metadata | [bytes](#bytes) |  | own_metadata is the identifying metadata that will be shared to the other account |

<a name="weshnet-protocol-v1-AccountContactRequestOutgoingSent"></a>

### AccountContactRequestOutgoingSent
AccountContactRequestOutgoingSent indicates that the account has sent a contact request

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| device_pk | [bytes](#bytes) |  | device_pk is the device sending the account event, signs the message |
| contact_pk | [bytes](#bytes) |  | contact_pk is the contacted account |

<a name="weshnet-protocol-v1-AccountContactRequestReferenceReset"></a>

### AccountContactRequestReferenceReset
AccountContactRequestReferenceReset indicates that the account should be advertised on different public rendezvous points

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| device_pk | [bytes](#bytes) |  | device_pk is the device sending the event, signs the message |
| public_rendezvous_seed | [bytes](#bytes) |  | public_rendezvous_seed is the new rendezvous point seed |

<a name="weshnet-protocol-v1-AccountContactUnblocked"></a>

### AccountContactUnblocked
AccountContactUnblocked indicates that a contact is unblocked

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| device_pk | [bytes](#bytes) |  | device_pk is the device sending the event, signs the message |
| contact_pk | [bytes](#bytes) |  | contact_pk is the contact unblocked |

<a name="weshnet-protocol-v1-AccountGroupJoined"></a>

### AccountGroupJoined
AccountGroupJoined indicates that the account is now part of a new group

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| device_pk | [bytes](#bytes) |  | device_pk is the device sending the event, signs the message |
| group | [Group](#weshnet-protocol-v1-Group) |  | group describe the joined group |

<a name="weshnet-protocol-v1-AccountGroupLeft"></a>

### AccountGroupLeft
AccountGroupLeft indicates that the account has left a group

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| device_pk | [bytes](#bytes) |  | device_pk is the device sending the event, signs the message |
| group_pk | [bytes](#bytes) |  | group_pk references the group left |

<a name="weshnet-protocol-v1-AccountServiceTokenAdded"></a>

### AccountServiceTokenAdded
AccountServiceTokenAdded indicates a token has been added to the account

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| device_pk | [bytes](#bytes) |  | device_pk is the device sending the event, signs the message |
| service_token | [ServiceToken](#weshnet-protocol-v1-ServiceToken) |  |  |

<a name="weshnet-protocol-v1-AccountServiceTokenRemoved"></a>

### AccountServiceTokenRemoved
AccountServiceTokenRemoved indicates a token has removed

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| device_pk | [bytes](#bytes) |  | device_pk is the device sending the event, signs the message |
| token_id | [string](#string) |  |  |

<a name="weshnet-protocol-v1-AccountVerifiedCredentialRegistered"></a>

### AccountVerifiedCredentialRegistered

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| device_pk | [bytes](#bytes) |  | device_pk is the public key of the device sending the message |
| signed_identity_public_key | [bytes](#bytes) |  |  |
| verified_credential | [string](#string) |  |  |
| registration_date | [int64](#int64) |  |  |
| expiration_date | [int64](#int64) |  |  |
| identifier | [string](#string) |  |  |
| issuer | [string](#string) |  |  |

<a name="weshnet-protocol-v1-ActivateGroup"></a>

### ActivateGroup

<a name="weshnet-protocol-v1-ActivateGroup-Reply"></a>

### ActivateGroup.Reply

<a name="weshnet-protocol-v1-ActivateGroup-Request"></a>

### ActivateGroup.Request

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| group_pk | [bytes](#bytes) |  | group_pk is the identifier of the group |
| local_only | [bool](#bool) |  | local_only will open the group without enabling network interactions with other members |

<a name="weshnet-protocol-v1-AppMessageSend"></a>

### AppMessageSend

<a name="weshnet-protocol-v1-AppMessageSend-Reply"></a>

### AppMessageSend.Reply

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| cid | [bytes](#bytes) |  |  |

<a name="weshnet-protocol-v1-AppMessageSend-Request"></a>

### AppMessageSend.Request

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| group_pk | [bytes](#bytes) |  | group_pk is the identifier of the group |
| payload | [bytes](#bytes) |  | payload is the payload to send |

<a name="weshnet-protocol-v1-AppMetadataSend"></a>

### AppMetadataSend

<a name="weshnet-protocol-v1-AppMetadataSend-Reply"></a>

### AppMetadataSend.Reply

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| cid | [bytes](#bytes) |  |  |

<a name="weshnet-protocol-v1-AppMetadataSend-Request"></a>

### AppMetadataSend.Request

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| group_pk | [bytes](#bytes) |  | group_pk is the identifier of the group |
| payload | [bytes](#bytes) |  | payload is the payload to send |

<a name="weshnet-protocol-v1-AuthExchangeResponse"></a>

### AuthExchangeResponse

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| access_token | [string](#string) |  |  |
| scope | [string](#string) |  |  |
| error | [string](#string) |  |  |
| error_description | [string](#string) |  |  |
| services | [AuthExchangeResponse.ServicesEntry](#weshnet-protocol-v1-AuthExchangeResponse-ServicesEntry) | repeated |  |

<a name="weshnet-protocol-v1-AuthExchangeResponse-ServicesEntry"></a>

### AuthExchangeResponse.ServicesEntry

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |

<a name="weshnet-protocol-v1-AuthServiceCompleteFlow"></a>

### AuthServiceCompleteFlow

<a name="weshnet-protocol-v1-AuthServiceCompleteFlow-Reply"></a>

### AuthServiceCompleteFlow.Reply

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| token_id | [string](#string) |  |  |

<a name="weshnet-protocol-v1-AuthServiceCompleteFlow-Request"></a>

### AuthServiceCompleteFlow.Request

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| callback_url | [string](#string) |  |  |

<a name="weshnet-protocol-v1-AuthServiceInitFlow"></a>

### AuthServiceInitFlow

<a name="weshnet-protocol-v1-AuthServiceInitFlow-Reply"></a>

### AuthServiceInitFlow.Reply

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| url | [string](#string) |  |  |
| secure_url | [bool](#bool) |  |  |

<a name="weshnet-protocol-v1-AuthServiceInitFlow-Request"></a>

### AuthServiceInitFlow.Request

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| auth_url | [string](#string) |  |  |
| services | [string](#string) | repeated |  |

<a name="weshnet-protocol-v1-ContactAliasKeyAdded"></a>

### ContactAliasKeyAdded
ContactAliasKeyAdded is an event type where ones shares their alias public key

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| device_pk | [bytes](#bytes) |  | device_pk is the device sending the event, signs the message |
| alias_pk | [bytes](#bytes) |  | alias_pk is the alias key which will be used to verify a contact identity |

<a name="weshnet-protocol-v1-ContactAliasKeySend"></a>

### ContactAliasKeySend

<a name="weshnet-protocol-v1-ContactAliasKeySend-Reply"></a>

### ContactAliasKeySend.Reply

<a name="weshnet-protocol-v1-ContactAliasKeySend-Request"></a>

### ContactAliasKeySend.Request

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| group_pk | [bytes](#bytes) |  | contact_pk is the identifier of the contact to send the alias public key to |

<a name="weshnet-protocol-v1-ContactBlock"></a>

### ContactBlock

<a name="weshnet-protocol-v1-ContactBlock-Reply"></a>

### ContactBlock.Reply

<a name="weshnet-protocol-v1-ContactBlock-Request"></a>

### ContactBlock.Request

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| contact_pk | [bytes](#bytes) |  | contact_pk is the identifier of the contact to block |

<a name="weshnet-protocol-v1-ContactRequestAccept"></a>

### ContactRequestAccept

<a name="weshnet-protocol-v1-ContactRequestAccept-Reply"></a>

### ContactRequestAccept.Reply

<a name="weshnet-protocol-v1-ContactRequestAccept-Request"></a>

### ContactRequestAccept.Request

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| contact_pk | [bytes](#bytes) |  | contact_pk is the identifier of the contact to accept the request from |

<a name="weshnet-protocol-v1-ContactRequestDisable"></a>

### ContactRequestDisable

<a name="weshnet-protocol-v1-ContactRequestDisable-Reply"></a>

### ContactRequestDisable.Reply

<a name="weshnet-protocol-v1-ContactRequestDisable-Request"></a>

### ContactRequestDisable.Request

<a name="weshnet-protocol-v1-ContactRequestDiscard"></a>

### ContactRequestDiscard

<a name="weshnet-protocol-v1-ContactRequestDiscard-Reply"></a>

### ContactRequestDiscard.Reply

<a name="weshnet-protocol-v1-ContactRequestDiscard-Request"></a>

### ContactRequestDiscard.Request

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| contact_pk | [bytes](#bytes) |  | contact_pk is the identifier of the contact to ignore the request from |

<a name="weshnet-protocol-v1-ContactRequestEnable"></a>

### ContactRequestEnable

<a name="weshnet-protocol-v1-ContactRequestEnable-Reply"></a>

### ContactRequestEnable.Reply

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| public_rendezvous_seed | [bytes](#bytes) |  | public_rendezvous_seed is the rendezvous seed used by the current account |

<a name="weshnet-protocol-v1-ContactRequestEnable-Request"></a>

### ContactRequestEnable.Request

<a name="weshnet-protocol-v1-ContactRequestReference"></a>

### ContactRequestReference

<a name="weshnet-protocol-v1-ContactRequestReference-Reply"></a>

### ContactRequestReference.Reply

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| public_rendezvous_seed | [bytes](#bytes) |  | public_rendezvous_seed is the rendezvous seed used by the current account |
| enabled | [bool](#bool) |  | enabled indicates if incoming contact requests are enabled |

<a name="weshnet-protocol-v1-ContactRequestReference-Request"></a>

### ContactRequestReference.Request

<a name="weshnet-protocol-v1-ContactRequestResetReference"></a>

### ContactRequestResetReference

<a name="weshnet-protocol-v1-ContactRequestResetReference-Reply"></a>

### ContactRequestResetReference.Reply

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| public_rendezvous_seed | [bytes](#bytes) |  | public_rendezvous_seed is the rendezvous seed used by the current account |

<a name="weshnet-protocol-v1-ContactRequestResetReference-Request"></a>

### ContactRequestResetReference.Request

<a name="weshnet-protocol-v1-ContactRequestSend"></a>

### ContactRequestSend

<a name="weshnet-protocol-v1-ContactRequestSend-Reply"></a>

### ContactRequestSend.Reply

<a name="weshnet-protocol-v1-ContactRequestSend-Request"></a>

### ContactRequestSend.Request

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| contact | [ShareableContact](#weshnet-protocol-v1-ShareableContact) |  | contact is a message describing how to connect to the other account |
| own_metadata | [bytes](#bytes) |  | own_metadata is the identifying metadata that will be shared to the other account |

<a name="weshnet-protocol-v1-ContactUnblock"></a>

### ContactUnblock

<a name="weshnet-protocol-v1-ContactUnblock-Reply"></a>

### ContactUnblock.Reply

<a name="weshnet-protocol-v1-ContactUnblock-Request"></a>

### ContactUnblock.Request

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| contact_pk | [bytes](#bytes) |  | contact_pk is the identifier of the contact to unblock |

<a name="weshnet-protocol-v1-CredentialVerificationServiceCompleteFlow"></a>

### CredentialVerificationServiceCompleteFlow

<a name="weshnet-protocol-v1-CredentialVerificationServiceCompleteFlow-Reply"></a>

### CredentialVerificationServiceCompleteFlow.Reply

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| identifier | [string](#string) |  |  |

<a name="weshnet-protocol-v1-CredentialVerificationServiceCompleteFlow-Request"></a>

### CredentialVerificationServiceCompleteFlow.Request

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| callback_uri | [string](#string) |  |  |

<a name="weshnet-protocol-v1-CredentialVerificationServiceInitFlow"></a>

### CredentialVerificationServiceInitFlow

<a name="weshnet-protocol-v1-CredentialVerificationServiceInitFlow-Reply"></a>

### CredentialVerificationServiceInitFlow.Reply

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| url | [string](#string) |  |  |
| secure_url | [bool](#bool) |  |  |

<a name="weshnet-protocol-v1-CredentialVerificationServiceInitFlow-Request"></a>

### CredentialVerificationServiceInitFlow.Request

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| service_url | [string](#string) |  |  |
| public_key | [bytes](#bytes) |  |  |
| link | [string](#string) |  |  |

<a name="weshnet-protocol-v1-DeactivateGroup"></a>

### DeactivateGroup

<a name="weshnet-protocol-v1-DeactivateGroup-Reply"></a>

### DeactivateGroup.Reply

<a name="weshnet-protocol-v1-DeactivateGroup-Request"></a>

### DeactivateGroup.Request

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| group_pk | [bytes](#bytes) |  | group_pk is the identifier of the group |

<a name="weshnet-protocol-v1-DebugAuthServiceSetToken"></a>

### DebugAuthServiceSetToken

<a name="weshnet-protocol-v1-DebugAuthServiceSetToken-Reply"></a>

### DebugAuthServiceSetToken.Reply

<a name="weshnet-protocol-v1-DebugAuthServiceSetToken-Request"></a>

### DebugAuthServiceSetToken.Request

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| token | [AuthExchangeResponse](#weshnet-protocol-v1-AuthExchangeResponse) |  |  |
| authentication_url | [string](#string) |  |  |

<a name="weshnet-protocol-v1-DebugGroup"></a>

### DebugGroup

<a name="weshnet-protocol-v1-DebugGroup-Reply"></a>

### DebugGroup.Reply

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| peer_ids | [string](#string) | repeated | peer_ids is the list of peer ids connected to the same group |

<a name="weshnet-protocol-v1-DebugGroup-Request"></a>

### DebugGroup.Request

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| group_pk | [bytes](#bytes) |  | group_pk is the identifier of the group |

<a name="weshnet-protocol-v1-DebugInspectGroupStore"></a>

### DebugInspectGroupStore

<a name="weshnet-protocol-v1-DebugInspectGroupStore-Reply"></a>

### DebugInspectGroupStore.Reply

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| cid | [bytes](#bytes) |  | cid is the CID of the IPFS log entry |
| parent_cids | [bytes](#bytes) | repeated | parent_cids is the list of the parent entries |
| metadata_event_type | [EventType](#weshnet-protocol-v1-EventType) |  | event_type metadata event type if subscribed to metadata events |
| device_pk | [bytes](#bytes) |  | device_pk is the public key of the device signing the entry |
| payload | [bytes](#bytes) |  | payload is the un encrypted entry payload if available |

<a name="weshnet-protocol-v1-DebugInspectGroupStore-Request"></a>

### DebugInspectGroupStore.Request

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| group_pk | [bytes](#bytes) |  | group_pk is the identifier of the group |
| log_type | [DebugInspectGroupLogType](#weshnet-protocol-v1-DebugInspectGroupLogType) |  | log_type is the log to inspect |

<a name="weshnet-protocol-v1-DebugListGroups"></a>

### DebugListGroups

<a name="weshnet-protocol-v1-DebugListGroups-Reply"></a>

### DebugListGroups.Reply

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| group_pk | [bytes](#bytes) |  | group_pk is the public key of the group |
| group_type | [GroupType](#weshnet-protocol-v1-GroupType) |  | group_type is the type of the group |
| contact_pk | [bytes](#bytes) |  | contact_pk is the contact public key if appropriate |

<a name="weshnet-protocol-v1-DebugListGroups-Request"></a>

### DebugListGroups.Request

<a name="weshnet-protocol-v1-DecodeContact"></a>

### DecodeContact

<a name="weshnet-protocol-v1-DecodeContact-Reply"></a>

### DecodeContact.Reply

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| contact | [ShareableContact](#weshnet-protocol-v1-ShareableContact) |  | shareable_contact is the decoded shareable contact. |

<a name="weshnet-protocol-v1-DecodeContact-Request"></a>

### DecodeContact.Request

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| encoded_contact | [bytes](#bytes) |  | encoded_contact is the Protobuf encoding of the shareable contact (as returned by ShareContact). |

<a name="weshnet-protocol-v1-DeviceChainKey"></a>

### DeviceChainKey
DeviceChainKey is a chain key, which will be encrypted for a specific member of the group

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| chain_key | [bytes](#bytes) |  | chain_key is the current value of the chain key of the group device |
| counter | [uint64](#uint64) |  | counter is the current value of the counter of the group device |

<a name="weshnet-protocol-v1-EncryptedMessage"></a>

### EncryptedMessage
EncryptedMessage is used in MessageEnvelope and only readable by groups members that joined before the message was sent

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| plaintext | [bytes](#bytes) |  | plaintext is the app layer data |
| protocol_metadata | [ProtocolMetadata](#weshnet-protocol-v1-ProtocolMetadata) |  | protocol_metadata is protocol layer data |

<a name="weshnet-protocol-v1-EventContext"></a>

### EventContext
EventContext adds context (its id, its parents and its attachments) to an event

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [bytes](#bytes) |  | id is the CID of the underlying OrbitDB event |
| parent_ids | [bytes](#bytes) | repeated | id are the the CIDs of the underlying parents of the OrbitDB event |
| group_pk | [bytes](#bytes) |  | group_pk receiving the event |

<a name="weshnet-protocol-v1-FirstLastCounters"></a>

### FirstLastCounters

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| first | [uint64](#uint64) |  |  |
| last | [uint64](#uint64) |  |  |

<a name="weshnet-protocol-v1-Group"></a>

### Group
Group define a group and is enough to invite someone to it

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| public_key | [bytes](#bytes) |  | public_key is the identifier of the group, it signs the group secret and the initial member of a multi-member group |
| secret | [bytes](#bytes) |  | secret is the symmetric secret of the group, which is used to encrypt the metadata |
| secret_sig | [bytes](#bytes) |  | secret_sig is the signature of the secret used to ensure the validity of the group |
| group_type | [GroupType](#weshnet-protocol-v1-GroupType) |  | group_type specifies the type of the group, used to determine how device chain key is generated |
| sign_pub | [bytes](#bytes) |  | sign_pub is the signature public key used to verify entries, not required when secret and secret_sig are provided |
| link_key | [bytes](#bytes) |  | link_key is the secret key used to exchange group updates and links to attachments, useful for replication services |
| link_key_sig | [bytes](#bytes) |  | link_key_sig is the signature of the link_key using the group private key |

<a name="weshnet-protocol-v1-GroupAddAdditionalRendezvousSeed"></a>

### GroupAddAdditionalRendezvousSeed
GroupAddAdditionalRendezvousSeed indicates that an additional rendezvous point should be used for data synchronization

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| device_pk | [bytes](#bytes) |  | device_pk is the device sending the event, signs the message, must be the device of an admin of the group |
| seed | [bytes](#bytes) |  | seed is the additional rendezvous point seed which should be used |

<a name="weshnet-protocol-v1-GroupDeviceChainKeyAdded"></a>

### GroupDeviceChainKeyAdded
GroupDeviceChainKeyAdded is an event which indicates to a group member a device chain key

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| device_pk | [bytes](#bytes) |  | device_pk is the device sending the event, signs the message |
| dest_member_pk | [bytes](#bytes) |  | dest_member_pk is the member who should receive the secret |
| payload | [bytes](#bytes) |  | payload is the serialization of Payload encrypted for the specified member |

<a name="weshnet-protocol-v1-GroupDeviceStatus"></a>

### GroupDeviceStatus

<a name="weshnet-protocol-v1-GroupDeviceStatus-Reply"></a>

### GroupDeviceStatus.Reply

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| type | [GroupDeviceStatus.Type](#weshnet-protocol-v1-GroupDeviceStatus-Type) |  |  |
| event | [bytes](#bytes) |  |  |

<a name="weshnet-protocol-v1-GroupDeviceStatus-Reply-PeerConnected"></a>

### GroupDeviceStatus.Reply.PeerConnected

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| peer_id | [string](#string) |  |  |
| device_pk | [bytes](#bytes) |  |  |
| transports | [GroupDeviceStatus.Transport](#weshnet-protocol-v1-GroupDeviceStatus-Transport) | repeated |  |
| maddrs | [string](#string) | repeated |  |

<a name="weshnet-protocol-v1-GroupDeviceStatus-Reply-PeerDisconnected"></a>

### GroupDeviceStatus.Reply.PeerDisconnected

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| peer_id | [string](#string) |  |  |

<a name="weshnet-protocol-v1-GroupDeviceStatus-Reply-PeerReconnecting"></a>

### GroupDeviceStatus.Reply.PeerReconnecting

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| peer_id | [string](#string) |  |  |

<a name="weshnet-protocol-v1-GroupDeviceStatus-Request"></a>

### GroupDeviceStatus.Request

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| group_pk | [bytes](#bytes) |  |  |

<a name="weshnet-protocol-v1-GroupEnvelope"></a>

### GroupEnvelope
GroupEnvelope is a publicly exposed structure containing a group metadata event

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| nonce | [bytes](#bytes) |  | nonce is used to encrypt the message |
| event | [bytes](#bytes) |  | event is encrypted using a symmetric key shared among group members |

<a name="weshnet-protocol-v1-GroupHeadsExport"></a>

### GroupHeadsExport

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| public_key | [bytes](#bytes) |  | public_key is the identifier of the group, it signs the group secret and the initial member of a multi-member group |
| sign_pub | [bytes](#bytes) |  | sign_pub is the signature public key used to verify entries |
| metadata_heads_cids | [bytes](#bytes) | repeated | metadata_heads_cids are the heads of the metadata store that should be restored from an export |
| messages_heads_cids | [bytes](#bytes) | repeated | messages_heads_cids are the heads of the metadata store that should be restored from an export |
| link_key | [bytes](#bytes) |  | link_key |

<a name="weshnet-protocol-v1-GroupInfo"></a>

### GroupInfo

<a name="weshnet-protocol-v1-GroupInfo-Reply"></a>

### GroupInfo.Reply

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| group | [Group](#weshnet-protocol-v1-Group) |  | group is the group invitation, containing the group pk and its type |
| member_pk | [bytes](#bytes) |  | member_pk is the identifier of the current member in the group |
| device_pk | [bytes](#bytes) |  | device_pk is the identifier of the current device in the group |

<a name="weshnet-protocol-v1-GroupInfo-Request"></a>

### GroupInfo.Request

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| group_pk | [bytes](#bytes) |  | group_pk is the identifier of the group |
| contact_pk | [bytes](#bytes) |  | contact_pk is the identifier of the contact |

<a name="weshnet-protocol-v1-GroupMemberDeviceAdded"></a>

### GroupMemberDeviceAdded
GroupMemberDeviceAdded is an event which indicates to a group a new device (and eventually a new member) is joining it
When added on AccountGroup, this event should be followed by appropriate GroupMemberDeviceAdded and GroupDeviceChainKeyAdded events

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| member_pk | [bytes](#bytes) |  | member_pk is the member sending the event |
| device_pk | [bytes](#bytes) |  | device_pk is the device sending the event, signs the message |
| member_sig | [bytes](#bytes) |  | member_sig is used to prove the ownership of the member pk

TODO: signature of what ??? ensure it can&#39;t be replayed |

<a name="weshnet-protocol-v1-GroupMessageEvent"></a>

### GroupMessageEvent

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| event_context | [EventContext](#weshnet-protocol-v1-EventContext) |  | event_context contains context information about the event |
| headers | [MessageHeaders](#weshnet-protocol-v1-MessageHeaders) |  | headers contains headers of the secure message |
| message | [bytes](#bytes) |  | message contains the secure message payload |

<a name="weshnet-protocol-v1-GroupMessageList"></a>

### GroupMessageList

<a name="weshnet-protocol-v1-GroupMessageList-Request"></a>

### GroupMessageList.Request

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| group_pk | [bytes](#bytes) |  | group_pk is the identifier of the group |
| since_id | [bytes](#bytes) |  | since is the lower ID bound used to filter events if not set, will return events since the beginning |
| since_now | [bool](#bool) |  | since_now will list only new event to come since_id must not be set |
| until_id | [bytes](#bytes) |  | until is the upper ID bound used to filter events if not set, will subscribe to new events to come |
| until_now | [bool](#bool) |  | until_now will not list new event to come until_id must not be set |
| reverse_order | [bool](#bool) |  | reverse_order indicates whether the previous events should be returned in reverse chronological order |

<a name="weshnet-protocol-v1-GroupMetadata"></a>

### GroupMetadata
GroupMetadata is used in GroupEnvelope and only readable by invited group members

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| event_type | [EventType](#weshnet-protocol-v1-EventType) |  | event_type defines which event type is used |
| payload | [bytes](#bytes) |  | the serialization depends on event_type, event is symmetrically encrypted |
| sig | [bytes](#bytes) |  | sig is the signature of the payload, it depends on the event_type for the used key |
| protocol_metadata | [ProtocolMetadata](#weshnet-protocol-v1-ProtocolMetadata) |  | protocol_metadata is protocol layer data |

<a name="weshnet-protocol-v1-GroupMetadataEvent"></a>

### GroupMetadataEvent

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| event_context | [EventContext](#weshnet-protocol-v1-EventContext) |  | event_context contains context information about the event |
| metadata | [GroupMetadata](#weshnet-protocol-v1-GroupMetadata) |  | metadata contains the newly available metadata |
| event | [bytes](#bytes) |  | event_clear clear bytes for the event |

<a name="weshnet-protocol-v1-GroupMetadataList"></a>

### GroupMetadataList

<a name="weshnet-protocol-v1-GroupMetadataList-Request"></a>

### GroupMetadataList.Request

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| group_pk | [bytes](#bytes) |  | group_pk is the identifier of the group |
| since_id | [bytes](#bytes) |  | since is the lower ID bound used to filter events if not set, will return events since the beginning |
| since_now | [bool](#bool) |  | since_now will list only new event to come since_id must not be set |
| until_id | [bytes](#bytes) |  | until is the upper ID bound used to filter events if not set, will subscribe to new events to come |
| until_now | [bool](#bool) |  | until_now will not list new event to come until_id must not be set |
| reverse_order | [bool](#bool) |  | reverse_order indicates whether the previous events should be returned in reverse chronological order |

<a name="weshnet-protocol-v1-GroupMetadataPayloadSent"></a>

### GroupMetadataPayloadSent
GroupMetadataPayloadSent is an app defined message, accessible to future group members

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| device_pk | [bytes](#bytes) |  | device_pk is the device sending the event, signs the message |
| message | [bytes](#bytes) |  | message is the payload |

<a name="weshnet-protocol-v1-GroupRemoveAdditionalRendezvousSeed"></a>

### GroupRemoveAdditionalRendezvousSeed
GroupRemoveAdditionalRendezvousSeed indicates that a previously added rendezvous point should be removed

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| device_pk | [bytes](#bytes) |  | device_pk is the device sending the event, signs the message, must be the device of an admin of the group |
| seed | [bytes](#bytes) |  | seed is the additional rendezvous point seed which should be removed |

<a name="weshnet-protocol-v1-GroupReplicating"></a>

### GroupReplicating

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| device_pk | [bytes](#bytes) |  | device_pk is the device sending the event, signs the message |
| authentication_url | [string](#string) |  | authentication_url indicates which server has been used for authentication |
| replication_server | [string](#string) |  | replication_server indicates which server will be used for replication |

<a name="weshnet-protocol-v1-MemberWithDevices"></a>

### MemberWithDevices

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| member_pk | [bytes](#bytes) |  |  |
| devices_pks | [bytes](#bytes) | repeated |  |

<a name="weshnet-protocol-v1-MessageEnvelope"></a>

### MessageEnvelope
MessageEnvelope is a publicly exposed structure containing a group secure message

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| message_headers | [bytes](#bytes) |  | message_headers is an encrypted serialization using a symmetric key of a MessageHeaders message |
| message | [bytes](#bytes) |  | message is an encrypted message, only readable by group members who previously received the appropriate chain key |
| nonce | [bytes](#bytes) |  | nonce is a nonce for message headers |

<a name="weshnet-protocol-v1-MessageHeaders"></a>

### MessageHeaders
MessageHeaders is used in MessageEnvelope and only readable by invited group members

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| counter | [uint64](#uint64) |  | counter is the current counter value for the specified device |
| device_pk | [bytes](#bytes) |  | device_pk is the public key of the device sending the message |
| sig | [bytes](#bytes) |  | sig is the signature of the encrypted message using the device&#39;s private key |
| metadata | [MessageHeaders.MetadataEntry](#weshnet-protocol-v1-MessageHeaders-MetadataEntry) | repeated | metadata allow to pass custom informations |

<a name="weshnet-protocol-v1-MessageHeaders-MetadataEntry"></a>

### MessageHeaders.MetadataEntry

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |

<a name="weshnet-protocol-v1-MultiMemberGroupAdminRoleGrant"></a>

### MultiMemberGroupAdminRoleGrant

<a name="weshnet-protocol-v1-MultiMemberGroupAdminRoleGrant-Reply"></a>

### MultiMemberGroupAdminRoleGrant.Reply

<a name="weshnet-protocol-v1-MultiMemberGroupAdminRoleGrant-Request"></a>

### MultiMemberGroupAdminRoleGrant.Request

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| group_pk | [bytes](#bytes) |  | group_pk is the identifier of the group |
| member_pk | [bytes](#bytes) |  | member_pk is the identifier of the member which will be granted the admin role |

<a name="weshnet-protocol-v1-MultiMemberGroupAdminRoleGranted"></a>

### MultiMemberGroupAdminRoleGranted
MultiMemberGroupAdminRoleGranted indicates that a group admin allows another group member to act as an admin

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| device_pk | [bytes](#bytes) |  | device_pk is the device sending the event, signs the message, must be the device of an admin of the group |
| grantee_member_pk | [bytes](#bytes) |  | grantee_member_pk is the member public key of the member granted of the admin role |

<a name="weshnet-protocol-v1-MultiMemberGroupAliasResolverAdded"></a>

### MultiMemberGroupAliasResolverAdded
MultiMemberGroupAliasResolverAdded indicates that a group member want to disclose their presence in the group to their contacts

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| device_pk | [bytes](#bytes) |  | device_pk is the device sending the event, signs the message |
| alias_resolver | [bytes](#bytes) |  | alias_resolver allows contact of an account to resolve the real identity behind an alias (Multi-Member Group Member) Generated by both contacts and account independently using: hmac(aliasPK, GroupID) |
| alias_proof | [bytes](#bytes) |  | alias_proof ensures that the associated alias_resolver has been issued by the right account Generated using aliasSKSig(GroupID) |

<a name="weshnet-protocol-v1-MultiMemberGroupAliasResolverDisclose"></a>

### MultiMemberGroupAliasResolverDisclose

<a name="weshnet-protocol-v1-MultiMemberGroupAliasResolverDisclose-Reply"></a>

### MultiMemberGroupAliasResolverDisclose.Reply

<a name="weshnet-protocol-v1-MultiMemberGroupAliasResolverDisclose-Request"></a>

### MultiMemberGroupAliasResolverDisclose.Request

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| group_pk | [bytes](#bytes) |  | group_pk is the identifier of the group |

<a name="weshnet-protocol-v1-MultiMemberGroupCreate"></a>

### MultiMemberGroupCreate

<a name="weshnet-protocol-v1-MultiMemberGroupCreate-Reply"></a>

### MultiMemberGroupCreate.Reply

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| group_pk | [bytes](#bytes) |  | group_pk is the identifier of the newly created group |

<a name="weshnet-protocol-v1-MultiMemberGroupCreate-Request"></a>

### MultiMemberGroupCreate.Request

<a name="weshnet-protocol-v1-MultiMemberGroupInitialMemberAnnounced"></a>

### MultiMemberGroupInitialMemberAnnounced
MultiMemberGroupInitialMemberAnnounced indicates that a member is the group creator, this event is signed using the group ID private key

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| member_pk | [bytes](#bytes) |  | member_pk is the public key of the member who is the group creator |

<a name="weshnet-protocol-v1-MultiMemberGroupInvitationCreate"></a>

### MultiMemberGroupInvitationCreate

<a name="weshnet-protocol-v1-MultiMemberGroupInvitationCreate-Reply"></a>

### MultiMemberGroupInvitationCreate.Reply

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| group | [Group](#weshnet-protocol-v1-Group) |  | group is the invitation to the group |

<a name="weshnet-protocol-v1-MultiMemberGroupInvitationCreate-Request"></a>

### MultiMemberGroupInvitationCreate.Request

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| group_pk | [bytes](#bytes) |  | group_pk is the identifier of the group |

<a name="weshnet-protocol-v1-MultiMemberGroupJoin"></a>

### MultiMemberGroupJoin

<a name="weshnet-protocol-v1-MultiMemberGroupJoin-Reply"></a>

### MultiMemberGroupJoin.Reply

<a name="weshnet-protocol-v1-MultiMemberGroupJoin-Request"></a>

### MultiMemberGroupJoin.Request

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| group | [Group](#weshnet-protocol-v1-Group) |  | group is the information of the group to join |

<a name="weshnet-protocol-v1-MultiMemberGroupLeave"></a>

### MultiMemberGroupLeave

<a name="weshnet-protocol-v1-MultiMemberGroupLeave-Reply"></a>

### MultiMemberGroupLeave.Reply

<a name="weshnet-protocol-v1-MultiMemberGroupLeave-Request"></a>

### MultiMemberGroupLeave.Request

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| group_pk | [bytes](#bytes) |  |  |

<a name="weshnet-protocol-v1-OrbitDBMessageHeads"></a>

### OrbitDBMessageHeads
OrbitDBMessageHeads is the payload sent on orbitdb to share peer&#39;s heads

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| sealed_box | [bytes](#bytes) |  | sealed box should contain encrypted Box |
| raw_rotation | [bytes](#bytes) |  | current topic used |

<a name="weshnet-protocol-v1-OrbitDBMessageHeads-Box"></a>

### OrbitDBMessageHeads.Box

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| address | [string](#string) |  |  |
| heads | [bytes](#bytes) |  |  |
| device_pk | [bytes](#bytes) |  |  |
| peer_id | [bytes](#bytes) |  |  |

<a name="weshnet-protocol-v1-OutOfStoreMessage"></a>

### OutOfStoreMessage

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| cid | [bytes](#bytes) |  |  |
| device_pk | [bytes](#bytes) |  |  |
| counter | [fixed64](#fixed64) |  |  |
| sig | [bytes](#bytes) |  |  |
| flags | [fixed32](#fixed32) |  |  |
| encrypted_payload | [bytes](#bytes) |  |  |
| nonce | [bytes](#bytes) |  |  |

<a name="weshnet-protocol-v1-OutOfStoreReceive"></a>

### OutOfStoreReceive

<a name="weshnet-protocol-v1-OutOfStoreReceive-Reply"></a>

### OutOfStoreReceive.Reply

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| message | [OutOfStoreMessage](#weshnet-protocol-v1-OutOfStoreMessage) |  |  |
| cleartext | [bytes](#bytes) |  |  |
| group_public_key | [bytes](#bytes) |  |  |
| already_received | [bool](#bool) |  |  |

<a name="weshnet-protocol-v1-OutOfStoreReceive-Request"></a>

### OutOfStoreReceive.Request

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| payload | [bytes](#bytes) |  |  |

<a name="weshnet-protocol-v1-OutOfStoreSeal"></a>

### OutOfStoreSeal

<a name="weshnet-protocol-v1-OutOfStoreSeal-Reply"></a>

### OutOfStoreSeal.Reply

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| encrypted | [bytes](#bytes) |  |  |

<a name="weshnet-protocol-v1-OutOfStoreSeal-Request"></a>

### OutOfStoreSeal.Request

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| cid | [bytes](#bytes) |  |  |
| group_public_key | [bytes](#bytes) |  |  |

<a name="weshnet-protocol-v1-PeerList"></a>

### PeerList

<a name="weshnet-protocol-v1-PeerList-Peer"></a>

### PeerList.Peer

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  | id is the libp2p.PeerID. |
| routes | [PeerList.Route](#weshnet-protocol-v1-PeerList-Route) | repeated | routes are the list of active and known maddr. |
| errors | [string](#string) | repeated | errors is a list of errors related to the peer. |
| features | [PeerList.Feature](#weshnet-protocol-v1-PeerList-Feature) | repeated | Features is a list of available features. |
| min_latency | [int64](#int64) |  | MinLatency is the minimum latency across all the peer routes. |
| is_active | [bool](#bool) |  | IsActive is true if at least one of the route is active. |
| direction | [Direction](#weshnet-protocol-v1-Direction) |  | Direction is the aggregate of all the routes&#39;s direction. |

<a name="weshnet-protocol-v1-PeerList-Reply"></a>

### PeerList.Reply

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| peers | [PeerList.Peer](#weshnet-protocol-v1-PeerList-Peer) | repeated |  |

<a name="weshnet-protocol-v1-PeerList-Request"></a>

### PeerList.Request

<a name="weshnet-protocol-v1-PeerList-Route"></a>

### PeerList.Route

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| is_active | [bool](#bool) |  | IsActive indicates whether the address is currently used or just known. |
| address | [string](#string) |  | Address is the multiaddress via which we are connected with the peer. |
| direction | [Direction](#weshnet-protocol-v1-Direction) |  | Direction is which way the connection was established. |
| latency | [int64](#int64) |  | Latency is the last known round trip time to the peer in ms. |
| streams | [PeerList.Stream](#weshnet-protocol-v1-PeerList-Stream) | repeated | Streams returns list of streams established with the peer. |

<a name="weshnet-protocol-v1-PeerList-Stream"></a>

### PeerList.Stream

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  | id is an identifier used to write protocol headers in streams. |

<a name="weshnet-protocol-v1-Progress"></a>

### Progress
Progress define a generic object that can be used to display a progress bar for long-running actions.

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| state | [string](#string) |  |  |
| doing | [string](#string) |  |  |
| progress | [float](#float) |  |  |
| completed | [uint64](#uint64) |  |  |
| total | [uint64](#uint64) |  |  |
| delay | [uint64](#uint64) |  |  |

<a name="weshnet-protocol-v1-ProtocolMetadata"></a>

### ProtocolMetadata

<a name="weshnet-protocol-v1-PushDeviceServerRegistered"></a>

### PushDeviceServerRegistered

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| server | [PushServer](#weshnet-protocol-v1-PushServer) |  |  |
| device_pk | [bytes](#bytes) |  | device_pk is the public key of the device sending the message |

<a name="weshnet-protocol-v1-PushDeviceTokenRegistered"></a>

### PushDeviceTokenRegistered

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| token | [PushServiceReceiver](#weshnet-protocol-v1-PushServiceReceiver) |  |  |
| device_pk | [bytes](#bytes) |  | device_pk is the public key of the device sending the message |

<a name="weshnet-protocol-v1-PushMemberTokenUpdate"></a>

### PushMemberTokenUpdate

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| server | [PushServer](#weshnet-protocol-v1-PushServer) |  |  |
| token | [bytes](#bytes) |  |  |
| device_pk | [bytes](#bytes) |  | device_pk is the public key of the device sending the message |

<a name="weshnet-protocol-v1-PushServer"></a>

### PushServer

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| server_key | [bytes](#bytes) |  |  |
| service_addr | [string](#string) |  |  |

<a name="weshnet-protocol-v1-PushServiceReceiver"></a>

### PushServiceReceiver

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| token_type | [weshnet.push.v1.PushServiceTokenType](#weshnet-push-v1-PushServiceTokenType) |  | token_type is the type of the token used, it allows us to act as a proxy to the appropriate push server |
| bundle_id | [string](#string) |  | bundle_id is the app identifier |
| token | [bytes](#bytes) |  | token is the device identifier used |
| recipient_public_key | [bytes](#bytes) |  | recipient_public_key is the public key which will be used to encrypt the payload |

<a name="weshnet-protocol-v1-RefreshContactRequest"></a>

### RefreshContactRequest

<a name="weshnet-protocol-v1-RefreshContactRequest-Peer"></a>

### RefreshContactRequest.Peer

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  | id is the libp2p.PeerID. |
| addrs | [string](#string) | repeated | list of peers multiaddrs. |

<a name="weshnet-protocol-v1-RefreshContactRequest-Reply"></a>

### RefreshContactRequest.Reply

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| peers_found | [RefreshContactRequest.Peer](#weshnet-protocol-v1-RefreshContactRequest-Peer) | repeated | peers found and successfully connected. |

<a name="weshnet-protocol-v1-RefreshContactRequest-Request"></a>

### RefreshContactRequest.Request

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| contact_pk | [bytes](#bytes) |  |  |
| timeout | [int64](#int64) |  | timeout in second |

<a name="weshnet-protocol-v1-ReplicationServiceRegisterGroup"></a>

### ReplicationServiceRegisterGroup

<a name="weshnet-protocol-v1-ReplicationServiceRegisterGroup-Reply"></a>

### ReplicationServiceRegisterGroup.Reply

<a name="weshnet-protocol-v1-ReplicationServiceRegisterGroup-Request"></a>

### ReplicationServiceRegisterGroup.Request

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| token_id | [string](#string) |  |  |
| group_pk | [bytes](#bytes) |  |  |

<a name="weshnet-protocol-v1-ReplicationServiceReplicateGroup"></a>

### ReplicationServiceReplicateGroup

<a name="weshnet-protocol-v1-ReplicationServiceReplicateGroup-Reply"></a>

### ReplicationServiceReplicateGroup.Reply

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ok | [bool](#bool) |  |  |

<a name="weshnet-protocol-v1-ReplicationServiceReplicateGroup-Request"></a>

### ReplicationServiceReplicateGroup.Request

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| group | [Group](#weshnet-protocol-v1-Group) |  |  |

<a name="weshnet-protocol-v1-ServiceExportData"></a>

### ServiceExportData

<a name="weshnet-protocol-v1-ServiceExportData-Reply"></a>

### ServiceExportData.Reply

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| exported_data | [bytes](#bytes) |  |  |

<a name="weshnet-protocol-v1-ServiceExportData-Request"></a>

### ServiceExportData.Request

<a name="weshnet-protocol-v1-ServiceGetConfiguration"></a>

### ServiceGetConfiguration

<a name="weshnet-protocol-v1-ServiceGetConfiguration-Reply"></a>

### ServiceGetConfiguration.Reply

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| account_pk | [bytes](#bytes) |  | account_pk is the public key of the current account |
| device_pk | [bytes](#bytes) |  | device_pk is the public key of the current device |
| account_group_pk | [bytes](#bytes) |  | account_group_pk is the public key of the account group |
| peer_id | [string](#string) |  | peer_id is the peer ID of the current IPFS node |
| listeners | [string](#string) | repeated | listeners is the list of swarm listening addresses of the current IPFS node |
| ble_enabled | [ServiceGetConfiguration.SettingState](#weshnet-protocol-v1-ServiceGetConfiguration-SettingState) |  |  |
| wifi_p2p_enabled | [ServiceGetConfiguration.SettingState](#weshnet-protocol-v1-ServiceGetConfiguration-SettingState) |  | MultiPeerConnectivity for Darwin and Nearby for Android |
| mdns_enabled | [ServiceGetConfiguration.SettingState](#weshnet-protocol-v1-ServiceGetConfiguration-SettingState) |  |  |
| relay_enabled | [ServiceGetConfiguration.SettingState](#weshnet-protocol-v1-ServiceGetConfiguration-SettingState) |  |  |
| device_push_token | [PushServiceReceiver](#weshnet-protocol-v1-PushServiceReceiver) |  |  |
| device_push_server | [PushServer](#weshnet-protocol-v1-PushServer) |  |  |

<a name="weshnet-protocol-v1-ServiceGetConfiguration-Request"></a>

### ServiceGetConfiguration.Request

<a name="weshnet-protocol-v1-ServiceToken"></a>

### ServiceToken

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| token | [string](#string) |  |  |
| authentication_url | [string](#string) |  |  |
| supported_services | [ServiceTokenSupportedService](#weshnet-protocol-v1-ServiceTokenSupportedService) | repeated |  |
| expiration | [int64](#int64) |  |  |

<a name="weshnet-protocol-v1-ServiceTokenSupportedService"></a>

### ServiceTokenSupportedService

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| service_type | [string](#string) |  |  |
| service_endpoint | [string](#string) |  |  |

<a name="weshnet-protocol-v1-ServicesTokenCode"></a>

### ServicesTokenCode

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| services | [string](#string) | repeated |  |
| code_challenge | [string](#string) |  |  |
| token_id | [string](#string) |  |  |

<a name="weshnet-protocol-v1-ServicesTokenList"></a>

### ServicesTokenList

<a name="weshnet-protocol-v1-ServicesTokenList-Reply"></a>

### ServicesTokenList.Reply

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| token_id | [string](#string) |  |  |
| service | [ServiceToken](#weshnet-protocol-v1-ServiceToken) |  |  |

<a name="weshnet-protocol-v1-ServicesTokenList-Request"></a>

### ServicesTokenList.Request

<a name="weshnet-protocol-v1-ShareContact"></a>

### ShareContact

<a name="weshnet-protocol-v1-ShareContact-Reply"></a>

### ShareContact.Reply

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| encoded_contact | [bytes](#bytes) |  | encoded_contact is the Protobuf encoding of the ShareableContact. You can further encode the bytes for sharing, such as base58 or QR code. |

<a name="weshnet-protocol-v1-ShareContact-Request"></a>

### ShareContact.Request

<a name="weshnet-protocol-v1-ShareableContact"></a>

### ShareableContact

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pk | [bytes](#bytes) |  | pk is the account to send a contact request to |
| public_rendezvous_seed | [bytes](#bytes) |  | public_rendezvous_seed is the rendezvous seed used by the account to send a contact request to |
| metadata | [bytes](#bytes) |  | metadata is the metadata specific to the app to identify the contact for the request |

<a name="weshnet-protocol-v1-SystemInfo"></a>

### SystemInfo

<a name="weshnet-protocol-v1-SystemInfo-OrbitDB"></a>

### SystemInfo.OrbitDB

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| account_metadata | [SystemInfo.OrbitDB.ReplicationStatus](#weshnet-protocol-v1-SystemInfo-OrbitDB-ReplicationStatus) |  |  |

<a name="weshnet-protocol-v1-SystemInfo-OrbitDB-ReplicationStatus"></a>

### SystemInfo.OrbitDB.ReplicationStatus

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| progress | [int64](#int64) |  |  |
| maximum | [int64](#int64) |  |  |
| buffered | [int64](#int64) |  |  |
| queued | [int64](#int64) |  |  |

<a name="weshnet-protocol-v1-SystemInfo-P2P"></a>

### SystemInfo.P2P

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| connected_peers | [int64](#int64) |  |  |

<a name="weshnet-protocol-v1-SystemInfo-Process"></a>

### SystemInfo.Process

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| version | [string](#string) |  |  |
| vcs_ref | [string](#string) |  |  |
| uptime_ms | [int64](#int64) |  |  |
| user_cpu_time_ms | [int64](#int64) |  |  |
| system_cpu_time_ms | [int64](#int64) |  |  |
| started_at | [int64](#int64) |  |  |
| rlimit_cur | [uint64](#uint64) |  |  |
| num_goroutine | [int64](#int64) |  |  |
| nofile | [int64](#int64) |  |  |
| too_many_open_files | [bool](#bool) |  |  |
| num_cpu | [int64](#int64) |  |  |
| go_version | [string](#string) |  |  |
| operating_system | [string](#string) |  |  |
| host_name | [string](#string) |  |  |
| arch | [string](#string) |  |  |
| rlimit_max | [uint64](#uint64) |  |  |
| pid | [int64](#int64) |  |  |
| ppid | [int64](#int64) |  |  |
| priority | [int64](#int64) |  |  |
| uid | [int64](#int64) |  |  |
| working_dir | [string](#string) |  |  |
| system_username | [string](#string) |  |  |

<a name="weshnet-protocol-v1-SystemInfo-Reply"></a>

### SystemInfo.Reply

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| process | [SystemInfo.Process](#weshnet-protocol-v1-SystemInfo-Process) |  |  |
| p2p | [SystemInfo.P2P](#weshnet-protocol-v1-SystemInfo-P2P) |  |  |
| orbitdb | [SystemInfo.OrbitDB](#weshnet-protocol-v1-SystemInfo-OrbitDB) |  |  |
| warns | [string](#string) | repeated |  |

<a name="weshnet-protocol-v1-SystemInfo-Request"></a>

### SystemInfo.Request

<a name="weshnet-protocol-v1-VerifiedCredentialsList"></a>

### VerifiedCredentialsList

<a name="weshnet-protocol-v1-VerifiedCredentialsList-Reply"></a>

### VerifiedCredentialsList.Reply

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| credential | [AccountVerifiedCredentialRegistered](#weshnet-protocol-v1-AccountVerifiedCredentialRegistered) |  |  |

<a name="weshnet-protocol-v1-VerifiedCredentialsList-Request"></a>

### VerifiedCredentialsList.Request

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| filter_identifier | [string](#string) |  |  |
| filter_issuer | [string](#string) |  |  |
| exclude_expired | [bool](#bool) |  |  |

 

<a name="weshnet-protocol-v1-ContactState"></a>

### ContactState

| Name | Number | Description |
| ---- | ------ | ----------- |
| ContactStateUndefined | 0 |  |
| ContactStateToRequest | 1 |  |
| ContactStateReceived | 2 |  |
| ContactStateAdded | 3 |  |
| ContactStateRemoved | 4 |  |
| ContactStateDiscarded | 5 |  |
| ContactStateBlocked | 6 |  |

<a name="weshnet-protocol-v1-DebugInspectGroupLogType"></a>

### DebugInspectGroupLogType

| Name | Number | Description |
| ---- | ------ | ----------- |
| DebugInspectGroupLogTypeUndefined | 0 |  |
| DebugInspectGroupLogTypeMessage | 1 |  |
| DebugInspectGroupLogTypeMetadata | 2 |  |

<a name="weshnet-protocol-v1-Direction"></a>

### Direction

| Name | Number | Description |
| ---- | ------ | ----------- |
| UnknownDir | 0 |  |
| InboundDir | 1 |  |
| OutboundDir | 2 |  |
| BiDir | 3 |  |

<a name="weshnet-protocol-v1-EventType"></a>

### EventType

| Name | Number | Description |
| ---- | ------ | ----------- |
| EventTypeUndefined | 0 | EventTypeUndefined indicates that the value has not been set. Should not happen. |
| EventTypeGroupMemberDeviceAdded | 1 | EventTypeGroupMemberDeviceAdded indicates the payload includes that a member has added their device to the group |
| EventTypeGroupDeviceChainKeyAdded | 2 | EventTypeGroupDeviceChainKeyAdded indicates the payload includes that a member has sent their device chain key to another member |
| EventTypeAccountGroupJoined | 101 | EventTypeAccountGroupJoined indicates the payload includes that the account has joined a group |
| EventTypeAccountGroupLeft | 102 | EventTypeAccountGroupLeft indicates the payload includes that the account has left a group |
| EventTypeAccountContactRequestDisabled | 103 | EventTypeAccountContactRequestDisabled indicates the payload includes that the account has disabled incoming contact requests |
| EventTypeAccountContactRequestEnabled | 104 | EventTypeAccountContactRequestEnabled indicates the payload includes that the account has enabled incoming contact requests |
| EventTypeAccountContactRequestReferenceReset | 105 | EventTypeAccountContactRequestReferenceReset indicates the payload includes that the account has a new contact request rendezvous seed |
| EventTypeAccountContactRequestOutgoingEnqueued | 106 | EventTypeAccountContactRequestOutgoingEnqueued indicates the payload includes that the account will attempt to send a new contact request |
| EventTypeAccountContactRequestOutgoingSent | 107 | EventTypeAccountContactRequestOutgoingSent indicates the payload includes that the account has sent a contact request |
| EventTypeAccountContactRequestIncomingReceived | 108 | EventTypeAccountContactRequestIncomingReceived indicates the payload includes that the account has received a contact request |
| EventTypeAccountContactRequestIncomingDiscarded | 109 | EventTypeAccountContactRequestIncomingDiscarded indicates the payload includes that the account has ignored a contact request |
| EventTypeAccountContactRequestIncomingAccepted | 110 | EventTypeAccountContactRequestIncomingAccepted indicates the payload includes that the account has accepted a contact request |
| EventTypeAccountContactBlocked | 111 | EventTypeAccountContactBlocked indicates the payload includes that the account has blocked a contact |
| EventTypeAccountContactUnblocked | 112 | EventTypeAccountContactUnblocked indicates the payload includes that the account has unblocked a contact |
| EventTypeContactAliasKeyAdded | 201 | EventTypeContactAliasKeyAdded indicates the payload includes that the contact group has received an alias key |
| EventTypeMultiMemberGroupAliasResolverAdded | 301 | EventTypeMultiMemberGroupAliasResolverAdded indicates the payload includes that a member of the group sent their alias proof |
| EventTypeMultiMemberGroupInitialMemberAnnounced | 302 | EventTypeMultiMemberGroupInitialMemberAnnounced indicates the payload includes that a member has authenticated themselves as the group owner |
| EventTypeMultiMemberGroupAdminRoleGranted | 303 | EventTypeMultiMemberGroupAdminRoleGranted indicates the payload includes that an admin of the group granted another member as an admin |
| EventTypeAccountServiceTokenAdded | 401 | EventTypeAccountServiceTokenAdded indicates that a new service provider has been registered for this account |
| EventTypeAccountServiceTokenRemoved | 402 | EventTypeAccountServiceTokenRemoved indicates that a service provider is not available anymore |
| EventTypeGroupReplicating | 403 | EventTypeGroupReplicating indicates that the group has been registered for replication on a server |
| EventTypePushMemberTokenUpdate | 404 | EventTypePushMemberTokenUpdate |
| EventTypePushDeviceTokenRegistered | 405 | EventTypePushDeviceTokenRegistered |
| EventTypePushDeviceServerRegistered | 406 | EventTypePushDeviceServerRegistered |
| EventTypeAccountVerifiedCredentialRegistered | 500 | EventTypeAccountVerifiedCredentialRegistered |
| EventTypeGroupMetadataPayloadSent | 1001 | EventTypeGroupMetadataPayloadSent indicates the payload includes an app specific event, unlike messages stored on the message store it is encrypted using a static key |

<a name="weshnet-protocol-v1-GroupDeviceStatus-Transport"></a>

### GroupDeviceStatus.Transport

| Name | Number | Description |
| ---- | ------ | ----------- |
| TptUnknown | 0 |  |
| TptLAN | 1 |  |
| TptWAN | 2 |  |
| TptProximity | 3 |  |

<a name="weshnet-protocol-v1-GroupDeviceStatus-Type"></a>

### GroupDeviceStatus.Type

| Name | Number | Description |
| ---- | ------ | ----------- |
| TypeUnknown | 0 |  |
| TypePeerDisconnected | 1 |  |
| TypePeerConnected | 2 |  |
| TypePeerReconnecting | 3 |  |

<a name="weshnet-protocol-v1-GroupType"></a>

### GroupType

| Name | Number | Description |
| ---- | ------ | ----------- |
| GroupTypeUndefined | 0 | GroupTypeUndefined indicates that the value has not been set. For example, happens if group is replicated. |
| GroupTypeAccount | 1 | GroupTypeAccount is the group managing an account, available to all its devices. |
| GroupTypeContact | 2 | GroupTypeContact is the group created between two accounts, available to all their devices. |
| GroupTypeMultiMember | 3 | GroupTypeMultiMember is a group containing an undefined number of members. |

<a name="weshnet-protocol-v1-PeerList-Feature"></a>

### PeerList.Feature

| Name | Number | Description |
| ---- | ------ | ----------- |
| UnknownFeature | 0 |  |
| WeshFeature | 1 |  |
| BLEFeature | 2 |  |
| LocalFeature | 3 |  |
| TorFeature | 4 |  |
| QuicFeature | 5 |  |

<a name="weshnet-protocol-v1-ServiceGetConfiguration-SettingState"></a>

### ServiceGetConfiguration.SettingState

| Name | Number | Description |
| ---- | ------ | ----------- |
| Unknown | 0 |  |
| Enabled | 1 |  |
| Disabled | 2 |  |
| Unavailable | 3 |  |

 

 

<a name="weshnet-protocol-v1-ProtocolService"></a>

### ProtocolService
ProtocolService is the top-level API to manage the Wesh protocol service.
Each active Wesh protocol service is considered as a Wesh device and is associated with a Wesh user.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| ServiceExportData | [ServiceExportData.Request](#weshnet-protocol-v1-ServiceExportData-Request) | [ServiceExportData.Reply](#weshnet-protocol-v1-ServiceExportData-Reply) stream | ServiceExportData exports the current data of the protocol service |
| ServiceGetConfiguration | [ServiceGetConfiguration.Request](#weshnet-protocol-v1-ServiceGetConfiguration-Request) | [ServiceGetConfiguration.Reply](#weshnet-protocol-v1-ServiceGetConfiguration-Reply) | ServiceGetConfiguration gets the current configuration of the protocol service |
| ContactRequestReference | [ContactRequestReference.Request](#weshnet-protocol-v1-ContactRequestReference-Request) | [ContactRequestReference.Reply](#weshnet-protocol-v1-ContactRequestReference-Reply) | ContactRequestReference retrieves the information required to create a reference (ie. included in a shareable link) to the current account |
| ContactRequestDisable | [ContactRequestDisable.Request](#weshnet-protocol-v1-ContactRequestDisable-Request) | [ContactRequestDisable.Reply](#weshnet-protocol-v1-ContactRequestDisable-Reply) | ContactRequestDisable disables incoming contact requests |
| ContactRequestEnable | [ContactRequestEnable.Request](#weshnet-protocol-v1-ContactRequestEnable-Request) | [ContactRequestEnable.Reply](#weshnet-protocol-v1-ContactRequestEnable-Reply) | ContactRequestEnable enables incoming contact requests |
| ContactRequestResetReference | [ContactRequestResetReference.Request](#weshnet-protocol-v1-ContactRequestResetReference-Request) | [ContactRequestResetReference.Reply](#weshnet-protocol-v1-ContactRequestResetReference-Reply) | ContactRequestResetReference changes the contact request reference |
| ContactRequestSend | [ContactRequestSend.Request](#weshnet-protocol-v1-ContactRequestSend-Request) | [ContactRequestSend.Reply](#weshnet-protocol-v1-ContactRequestSend-Reply) | ContactRequestSend attempt to send a contact request |
| ContactRequestAccept | [ContactRequestAccept.Request](#weshnet-protocol-v1-ContactRequestAccept-Request) | [ContactRequestAccept.Reply](#weshnet-protocol-v1-ContactRequestAccept-Reply) | ContactRequestAccept accepts a contact request |
| ContactRequestDiscard | [ContactRequestDiscard.Request](#weshnet-protocol-v1-ContactRequestDiscard-Request) | [ContactRequestDiscard.Reply](#weshnet-protocol-v1-ContactRequestDiscard-Reply) | ContactRequestDiscard ignores a contact request, without informing the other user |
| ShareContact | [ShareContact.Request](#weshnet-protocol-v1-ShareContact-Request) | [ShareContact.Reply](#weshnet-protocol-v1-ShareContact-Reply) | ShareContact uses ContactRequestReference to get the contact information for the current account and returns the Protobuf encoding of a shareable contact which you can further encode and share. If needed, this will reset the contact request reference and enable contact requests. To decode the result, see DecodeContact. |
| DecodeContact | [DecodeContact.Request](#weshnet-protocol-v1-DecodeContact-Request) | [DecodeContact.Reply](#weshnet-protocol-v1-DecodeContact-Reply) | DecodeContact decodes the Protobuf encoding of a shareable contact which was returned by ShareContact. |
| ContactBlock | [ContactBlock.Request](#weshnet-protocol-v1-ContactBlock-Request) | [ContactBlock.Reply](#weshnet-protocol-v1-ContactBlock-Reply) | ContactBlock blocks a contact from sending requests |
| ContactUnblock | [ContactUnblock.Request](#weshnet-protocol-v1-ContactUnblock-Request) | [ContactUnblock.Reply](#weshnet-protocol-v1-ContactUnblock-Reply) | ContactUnblock unblocks a contact from sending requests |
| ContactAliasKeySend | [ContactAliasKeySend.Request](#weshnet-protocol-v1-ContactAliasKeySend-Request) | [ContactAliasKeySend.Reply](#weshnet-protocol-v1-ContactAliasKeySend-Reply) | ContactAliasKeySend send an alias key to a contact, the contact will be able to assert that your account is being present on a multi-member group |
| MultiMemberGroupCreate | [MultiMemberGroupCreate.Request](#weshnet-protocol-v1-MultiMemberGroupCreate-Request) | [MultiMemberGroupCreate.Reply](#weshnet-protocol-v1-MultiMemberGroupCreate-Reply) | MultiMemberGroupCreate creates a new multi-member group |
| MultiMemberGroupJoin | [MultiMemberGroupJoin.Request](#weshnet-protocol-v1-MultiMemberGroupJoin-Request) | [MultiMemberGroupJoin.Reply](#weshnet-protocol-v1-MultiMemberGroupJoin-Reply) | MultiMemberGroupJoin joins a multi-member group |
| MultiMemberGroupLeave | [MultiMemberGroupLeave.Request](#weshnet-protocol-v1-MultiMemberGroupLeave-Request) | [MultiMemberGroupLeave.Reply](#weshnet-protocol-v1-MultiMemberGroupLeave-Reply) | MultiMemberGroupLeave leaves a multi-member group |
| MultiMemberGroupAliasResolverDisclose | [MultiMemberGroupAliasResolverDisclose.Request](#weshnet-protocol-v1-MultiMemberGroupAliasResolverDisclose-Request) | [MultiMemberGroupAliasResolverDisclose.Reply](#weshnet-protocol-v1-MultiMemberGroupAliasResolverDisclose-Reply) | MultiMemberGroupAliasResolverDisclose discloses your alias resolver key |
| MultiMemberGroupAdminRoleGrant | [MultiMemberGroupAdminRoleGrant.Request](#weshnet-protocol-v1-MultiMemberGroupAdminRoleGrant-Request) | [MultiMemberGroupAdminRoleGrant.Reply](#weshnet-protocol-v1-MultiMemberGroupAdminRoleGrant-Reply) | MultiMemberGroupAdminRoleGrant grants an admin role to a group member |
| MultiMemberGroupInvitationCreate | [MultiMemberGroupInvitationCreate.Request](#weshnet-protocol-v1-MultiMemberGroupInvitationCreate-Request) | [MultiMemberGroupInvitationCreate.Reply](#weshnet-protocol-v1-MultiMemberGroupInvitationCreate-Reply) | MultiMemberGroupInvitationCreate creates an invitation to a multi-member group |
| AppMetadataSend | [AppMetadataSend.Request](#weshnet-protocol-v1-AppMetadataSend-Request) | [AppMetadataSend.Reply](#weshnet-protocol-v1-AppMetadataSend-Reply) | AppMetadataSend adds an app event to the metadata store, the message is encrypted using a symmetric key and readable by future group members |
| AppMessageSend | [AppMessageSend.Request](#weshnet-protocol-v1-AppMessageSend-Request) | [AppMessageSend.Reply](#weshnet-protocol-v1-AppMessageSend-Reply) | AppMessageSend adds an app event to the message store, the message is encrypted using a derived key and readable by current group members |
| GroupMetadataList | [GroupMetadataList.Request](#weshnet-protocol-v1-GroupMetadataList-Request) | [GroupMetadataEvent](#weshnet-protocol-v1-GroupMetadataEvent) stream | GroupMetadataList replays previous and subscribes to new metadata events from the group |
| GroupMessageList | [GroupMessageList.Request](#weshnet-protocol-v1-GroupMessageList-Request) | [GroupMessageEvent](#weshnet-protocol-v1-GroupMessageEvent) stream | GroupMessageList replays previous and subscribes to new message events from the group |
| GroupInfo | [GroupInfo.Request](#weshnet-protocol-v1-GroupInfo-Request) | [GroupInfo.Reply](#weshnet-protocol-v1-GroupInfo-Reply) | GroupInfo retrieves information about a group |
| ActivateGroup | [ActivateGroup.Request](#weshnet-protocol-v1-ActivateGroup-Request) | [ActivateGroup.Reply](#weshnet-protocol-v1-ActivateGroup-Reply) | ActivateGroup explicitly opens a group |
| DeactivateGroup | [DeactivateGroup.Request](#weshnet-protocol-v1-DeactivateGroup-Request) | [DeactivateGroup.Reply](#weshnet-protocol-v1-DeactivateGroup-Reply) | DeactivateGroup closes a group |
| GroupDeviceStatus | [GroupDeviceStatus.Request](#weshnet-protocol-v1-GroupDeviceStatus-Request) | [GroupDeviceStatus.Reply](#weshnet-protocol-v1-GroupDeviceStatus-Reply) stream | GroupDeviceStatus monitor device status |
| DebugListGroups | [DebugListGroups.Request](#weshnet-protocol-v1-DebugListGroups-Request) | [DebugListGroups.Reply](#weshnet-protocol-v1-DebugListGroups-Reply) stream |  |
| DebugInspectGroupStore | [DebugInspectGroupStore.Request](#weshnet-protocol-v1-DebugInspectGroupStore-Request) | [DebugInspectGroupStore.Reply](#weshnet-protocol-v1-DebugInspectGroupStore-Reply) stream |  |
| DebugGroup | [DebugGroup.Request](#weshnet-protocol-v1-DebugGroup-Request) | [DebugGroup.Reply](#weshnet-protocol-v1-DebugGroup-Reply) |  |
| DebugAuthServiceSetToken | [DebugAuthServiceSetToken.Request](#weshnet-protocol-v1-DebugAuthServiceSetToken-Request) | [DebugAuthServiceSetToken.Reply](#weshnet-protocol-v1-DebugAuthServiceSetToken-Reply) |  |
| SystemInfo | [SystemInfo.Request](#weshnet-protocol-v1-SystemInfo-Request) | [SystemInfo.Reply](#weshnet-protocol-v1-SystemInfo-Reply) |  |
| AuthServiceInitFlow | [AuthServiceInitFlow.Request](#weshnet-protocol-v1-AuthServiceInitFlow-Request) | [AuthServiceInitFlow.Reply](#weshnet-protocol-v1-AuthServiceInitFlow-Reply) | AuthServiceInitFlow Initialize an authentication flow |
| AuthServiceCompleteFlow | [AuthServiceCompleteFlow.Request](#weshnet-protocol-v1-AuthServiceCompleteFlow-Request) | [AuthServiceCompleteFlow.Reply](#weshnet-protocol-v1-AuthServiceCompleteFlow-Reply) | AuthServiceCompleteFlow Completes an authentication flow |
| CredentialVerificationServiceInitFlow | [CredentialVerificationServiceInitFlow.Request](#weshnet-protocol-v1-CredentialVerificationServiceInitFlow-Request) | [CredentialVerificationServiceInitFlow.Reply](#weshnet-protocol-v1-CredentialVerificationServiceInitFlow-Reply) | CredentialVerificationServiceInitFlow Initialize a credential verification flow |
| CredentialVerificationServiceCompleteFlow | [CredentialVerificationServiceCompleteFlow.Request](#weshnet-protocol-v1-CredentialVerificationServiceCompleteFlow-Request) | [CredentialVerificationServiceCompleteFlow.Reply](#weshnet-protocol-v1-CredentialVerificationServiceCompleteFlow-Reply) | CredentialVerificationServiceCompleteFlow Completes a credential verification flow |
| VerifiedCredentialsList | [VerifiedCredentialsList.Request](#weshnet-protocol-v1-VerifiedCredentialsList-Request) | [VerifiedCredentialsList.Reply](#weshnet-protocol-v1-VerifiedCredentialsList-Reply) stream | VerifiedCredentialsList Retrieves the list of verified credentials |
| ServicesTokenList | [ServicesTokenList.Request](#weshnet-protocol-v1-ServicesTokenList-Request) | [ServicesTokenList.Reply](#weshnet-protocol-v1-ServicesTokenList-Reply) stream | ServicesTokenList Retrieves the list of services tokens |
| ReplicationServiceRegisterGroup | [ReplicationServiceRegisterGroup.Request](#weshnet-protocol-v1-ReplicationServiceRegisterGroup-Request) | [ReplicationServiceRegisterGroup.Reply](#weshnet-protocol-v1-ReplicationServiceRegisterGroup-Reply) | ReplicationServiceRegisterGroup Asks a replication service to distribute a group contents |
| PeerList | [PeerList.Request](#weshnet-protocol-v1-PeerList-Request) | [PeerList.Reply](#weshnet-protocol-v1-PeerList-Reply) | PeerList returns a list of P2P peers |
| OutOfStoreReceive | [OutOfStoreReceive.Request](#weshnet-protocol-v1-OutOfStoreReceive-Request) | [OutOfStoreReceive.Reply](#weshnet-protocol-v1-OutOfStoreReceive-Reply) | OutOfStoreReceive parses a payload received outside a synchronized store |
| OutOfStoreSeal | [OutOfStoreSeal.Request](#weshnet-protocol-v1-OutOfStoreSeal-Request) | [OutOfStoreSeal.Reply](#weshnet-protocol-v1-OutOfStoreSeal-Reply) | OutOfStoreSeal creates a payload of a message present in store to be sent outside a synchronized store |
| RefreshContactRequest | [RefreshContactRequest.Request](#weshnet-protocol-v1-RefreshContactRequest-Request) | [RefreshContactRequest.Reply](#weshnet-protocol-v1-RefreshContactRequest-Reply) | RefreshContactRequest try to refresh the contact request for the given contact |

 

## Scalar Value Types

| .proto Type | Notes | C++ | Java | Python | Go | C# | PHP | Ruby |
| ----------- | ----- | --- | ---- | ------ | -- | -- | --- | ---- |
| <a name="double" /> double |  | double | double | float | float64 | double | float | Float |
| <a name="float" /> float |  | float | float | float | float32 | float | float | Float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum or Fixnum (as required) |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="bool" /> bool |  | bool | boolean | boolean | bool | bool | boolean | TrueClass/FalseClass |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode | string | string | string | String (UTF-8) |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str | []byte | ByteString | string | String (ASCII-8BIT) |

