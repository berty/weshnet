package testutil

import (
	"testing"

	"google.golang.org/protobuf/proto"

	"berty.tech/weshnet/v2/pkg/protocoltypes"
)

func TestFilterGroupMetadataPayloadSent(t *testing.T, events <-chan *protocoltypes.GroupMetadataEvent) []*protocoltypes.GroupMetadataPayloadSent {
	t.Helper()

	out := []*protocoltypes.GroupMetadataPayloadSent(nil)

	for evt := range events {
		if evt == nil {
			continue
		}

		if evt.Metadata.EventType != protocoltypes.EventType_EventTypeGroupMetadataPayloadSent {
			continue
		}

		m := &protocoltypes.GroupMetadataPayloadSent{}
		if err := proto.Unmarshal(evt.Event, m); err != nil {
			continue
		}

		out = append(out, m)
	}

	return out
}
