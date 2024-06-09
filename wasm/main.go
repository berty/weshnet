//go:build js

package main

import (
	"context"
	"encoding/base64"
	"errors"
	"os"
	"syscall/js"

	"berty.tech/weshnet"
	"berty.tech/weshnet/pkg/protocoltypes"
)

var svc weshnet.ServiceClient

func main() {
	c := make(chan struct{}, 0)
	js.Global().Set("weshnet_initService", js.FuncOf(initService))
	js.Global().Set("weshnet_contactRequestReference", js.FuncOf(contactRequestReference))
	js.Global().Set("weshnet_serviceGetConfiguration", js.FuncOf(serviceGetConfiguration))
	js.Global().Set("weshnet_groupMessageList", js.FuncOf(groupMessageList))
	js.Global().Set("weshnet_groupMetadataList", js.FuncOf(groupMetadataList))
	js.Global().Set("weshnet_multiMemberGroupCreate", js.FuncOf(multiMemberGroupCreate))
	js.Global().Set("weshnet_multiMemberGroupInvitationCreate", js.FuncOf(multiMemberGroupInvitationCreate))
	js.Global().Set("weshnet_multiMemberGroupJoin", js.FuncOf(multiMemberGroupJoin))
	js.Global().Set("weshnet_activateGroup", js.FuncOf(activateGroup))
	js.Global().Set("weshnet_peerList", js.FuncOf(peerList))
	js.Global().Set("weshnet_appMessageSend", js.FuncOf(appMessageSend))
	js.Global().Set("weshnet_appMetadataSend", js.FuncOf(appMetadataSend))
	<-c
}

func initService(this js.Value, args []js.Value) any {
	return promisify(func() ([]any, error) {
		if svc != nil {
			return nil, errors.New("weshnet already initialized")
		}
		if len(args) != 1 {
			return nil, errors.New("expected 1 arg")
		}
		helia := args[0]
		ipfsCore := &coreAPIFromJS{helia: helia}
		os.Setenv("WESHNET_LOG_FILTER", "*")
		os.Setenv("IPFS_LOGGING", "debug")
		var err error
		svc, err = weshnet.NewServiceClient(weshnet.Opts{
			DatastoreDir: weshnet.InMemoryDirectory,
			IpfsCoreAPI:  ipfsCore,
		})
		if err != nil {
			return nil, err
		}
		return nil, nil
	})
}

func serviceGetConfiguration(this js.Value, args []js.Value) any {
	return promisify(func() ([]any, error) {
		res, err := svc.ServiceGetConfiguration(context.Background(), &protocoltypes.ServiceGetConfiguration_Request{})
		if err != nil {
			return nil, err
		}
		return []any{map[string]any{
			"AccountPK":      bytesToString(res.AccountPK),
			"DevicePK":       bytesToString(res.DevicePK),
			"AccountGroupPK": bytesToString(res.AccountGroupPK),
			"PeerID":         res.PeerID,
			"Listeners":      jsArray(res.Listeners),
		}}, nil
	})
}

func groupMessageList(this js.Value, args []js.Value) any {
	return promisify(func() ([]any, error) {
		assertInitialized()

		if len(args) != 2 {
			return nil, errors.New("expected 2 args")
		}
		if args[0].Type() != js.TypeString {
			return nil, errors.New("expected first arg to be a string")
		}
		if args[1].Type() != js.TypeFunction {
			return nil, errors.New("expected second arg to be a function")
		}

		stream, err := svc.GroupMessageList(context.Background(), &protocoltypes.GroupMessageList_Request{
			GroupPK: mustStringToBytes(args[0].String()),
		})
		if err != nil {
			return nil, err
		}
		for {
			msg, err := stream.Recv()
			if err != nil {
				panic(err)
			}
			args[1].Invoke(map[string]any{
				"EventContext": jsEventContext(msg.EventContext),
				"Message":      bytesToString(msg.Message),
			})
		}
	})
}

func groupMetadataList(this js.Value, args []js.Value) any {
	return promisify(func() ([]any, error) {
		assertInitialized()

		if len(args) != 2 {
			return nil, errors.New("expected 2 args")
		}
		if args[0].Type() != js.TypeString {
			return nil, errors.New("expected first arg to be a string")
		}
		if args[1].Type() != js.TypeFunction {
			return nil, errors.New("expected second arg to be a function")
		}

		stream, err := svc.GroupMetadataList(context.Background(), &protocoltypes.GroupMetadataList_Request{
			GroupPK: mustStringToBytes(args[0].String()),
		})
		if err != nil {
			return nil, err
		}

		for {
			msg, err := stream.Recv()
			if err != nil {
				return nil, err
			}
			args[1].Invoke(map[string]any{
				"EventContext": jsEventContext(msg.EventContext),
				"Metadata":     jsGroupMetadata(msg.Metadata),
				"Event":        bytesToString(msg.Event),
			})
		}
	})
}

func contactRequestReference(this js.Value, args []js.Value) any {
	return promisify(func() ([]any, error) {
		assertInitialized()

		res, err := svc.ContactRequestEnable(context.Background(), &protocoltypes.ContactRequestEnable_Request{})
		if err != nil {
			return nil, err
		}

		ret := res.PublicRendezvousSeed
		if len(ret) == 0 {
			res, err := svc.ContactRequestResetReference(context.Background(), &protocoltypes.ContactRequestResetReference_Request{})
			if err != nil {
				return nil, err
			}
			ret = res.PublicRendezvousSeed
		}

		return []any{js.ValueOf(bytesToString(ret))}, nil
	})
}

func multiMemberGroupCreate(this js.Value, args []js.Value) any {
	return promisify(func() ([]any, error) {
		res, err := svc.MultiMemberGroupCreate(context.Background(), &protocoltypes.MultiMemberGroupCreate_Request{})
		if err != nil {
			return nil, err
		}
		return []any{bytesToString(res.GroupPK)}, nil
	})
}

func multiMemberGroupInvitationCreate(this js.Value, args []js.Value) any {
	return promisify(func() ([]any, error) {
		if len(args) != 1 {
			return nil, errors.New("expected 1 arg")
		}
		if args[0].Type() != js.TypeString {
			return nil, errors.New("expected first arg to be a string")
		}
		groupPK := mustStringToBytes(args[0].String())
		res, err := svc.MultiMemberGroupInvitationCreate(context.Background(), &protocoltypes.MultiMemberGroupInvitationCreate_Request{
			GroupPK: groupPK,
		})
		if err != nil {
			return nil, err
		}
		invit, err := res.Group.Marshal()
		if err != nil {
			return nil, err
		}
		return []any{bytesToString(invit)}, nil
	})
}

func multiMemberGroupJoin(this js.Value, args []js.Value) any {
	return promisify(func() ([]any, error) {
		if len(args) != 1 {
			return nil, errors.New("expected 1 arg")
		}
		if args[0].Type() != js.TypeString {
			return nil, errors.New("expected first arg to be a string")
		}
		invit := mustStringToBytes(args[0].String())
		group := &protocoltypes.Group{}
		err := group.Unmarshal(invit)
		if err != nil {
			panic(err)
		}
		_, err = svc.MultiMemberGroupJoin(context.Background(), &protocoltypes.MultiMemberGroupJoin_Request{
			Group: group,
		})
		if err != nil {
			return nil, err
		}
		return []any{bytesToString(group.PublicKey)}, nil
	})
}

func activateGroup(this js.Value, args []js.Value) any {
	return promisify(func() ([]any, error) {
		if len(args) != 1 {
			return nil, errors.New("expected 1 arg")
		}
		if args[0].Type() != js.TypeString {
			return nil, errors.New("expected first arg to be a string")
		}
		groupPK := mustStringToBytes(args[0].String())
		_, err := svc.ActivateGroup(context.Background(), &protocoltypes.ActivateGroup_Request{GroupPK: groupPK})
		if err != nil {
			return nil, err
		}
		return nil, nil
	})
}

func peerList(this js.Value, args []js.Value) any {
	return promisify(func() ([]any, error) {
		res, err := svc.PeerList(context.Background(), &protocoltypes.PeerList_Request{})
		if err != nil {
			return nil, err
		}
		return []any{jsArrayTransform(res.Peers, jsPeer)}, nil
	})
}

func appMetadataSend(this js.Value, args []js.Value) any {
	return promisify(func() ([]any, error) {
		if len(args) != 2 {
			return nil, errors.New("expected 2 args")
		}
		if args[0].Type() != js.TypeString {
			return nil, errors.New("expected groupPk arg to be a string")
		}
		if args[1].Type() != js.TypeString {
			return nil, errors.New("expected payload arg to be a string")
		}
		groupPK := mustStringToBytes(args[0].String())
		payload := mustStringToBytes(args[1].String())
		res, err := svc.AppMetadataSend(context.Background(), &protocoltypes.AppMetadataSend_Request{
			GroupPK: groupPK,
			Payload: payload,
		})
		if err != nil {
			return nil, err
		}
		return []any{bytesToString(res.CID)}, nil
	})
}

func appMessageSend(this js.Value, args []js.Value) any {
	return promisify(func() ([]any, error) {
		if len(args) != 2 {
			return nil, errors.New("expected 2 args")
		}
		if args[0].Type() != js.TypeString {
			return nil, errors.New("expected groupPk arg to be a string")
		}
		if args[1].Type() != js.TypeString {
			return nil, errors.New("expected payload arg to be a string")
		}
		groupPK := mustStringToBytes(args[0].String())
		payload := mustStringToBytes(args[1].String())
		res, err := svc.AppMessageSend(context.Background(), &protocoltypes.AppMessageSend_Request{
			GroupPK: groupPK,
			Payload: payload,
		})
		if err != nil {
			return nil, err
		}
		return []any{bytesToString(res.CID)}, nil
	})
}

func jsPeer(p *protocoltypes.PeerList_Peer) map[string]any {
	return map[string]any{
		"ID":         p.ID,
		"Routes":     jsArrayTransform(p.Routes, jsPeerRoute),
		"Errors":     jsArray(p.Errors),
		"Features":   jsArrayTransform(p.Features, jsPeerFeature),
		"MinLatency": p.MinLatency,
		"IsActive":   p.IsActive,
		"Direction":  p.Direction.String(),
	}
}

func jsPeerFeature(f protocoltypes.PeerList_Feature) any {
	return f.String()
}

func jsPeerRoute(r *protocoltypes.PeerList_Route) map[string]any {
	return map[string]any{
		"IsActive":  r.IsActive,
		"Address":   r.Address,
		"Direction": r.Direction.String(),
		"Latency":   r.Latency,
		"Streams":   jsArrayTransform(r.Streams, jsPeerStream),
	}
}

func jsPeerStream(s *protocoltypes.PeerList_Stream) map[string]any {
	return map[string]any{
		"ID": s.ID,
	}
}

func jsEventContext(eventContext *protocoltypes.EventContext) map[string]any {
	return map[string]any{
		"ID":         bytesToString(eventContext.ID),
		"ParentsIDs": jsArrayTransform(eventContext.ParentIDs, bytesToString),
		"GroupPK":    bytesToString(eventContext.GroupPK),
	}
}

func jsGroupMetadata(gm *protocoltypes.GroupMetadata) map[string]any {
	return map[string]any{
		"EventType": gm.EventType.String(),
		"Payload":   bytesToString(gm.Payload),
		"Sig":       bytesToString(gm.Sig),
	}
}

func bytesToString(bs []byte) string {
	return base64.RawURLEncoding.EncodeToString(bs)
}

func mustStringToBytes(s string) []byte {
	bs, err := base64.RawURLEncoding.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return bs
}

func assertInitialized() {
	if svc == nil {
		panic("weshnet is not initialized, call `weshnet_initService` first")
	}
}

func jsArray[V any](in []V) []any {
	out := make([]any, 0, len(in))
	for _, elem := range in {
		out = append(out, elem)
	}
	return out
}

func jsArrayTransform[I any, O any](in []I, transform func(I) O) []any {
	out := make([]any, 0, len(in))
	for _, elem := range in {
		out = append(out, transform(elem))
	}
	return out
}
