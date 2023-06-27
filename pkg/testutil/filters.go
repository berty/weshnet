package testutil

import (
	"testing"

	"berty.tech/weshnet/pkg/protocoltypes"
)

func TestFilterGroupMetadataPayloadSent(t *testing.T, events <-chan *protocoltypes.GroupMetadataEvent) []*protocoltypes.GroupMetadataPayloadSent {
	t.Helper()

	out := []*protocoltypes.GroupMetadataPayloadSent(nil)

	for evt := range events {
		if evt == nil {
			continue
		}

		if evt.Metadata.EventType != protocoltypes.EventTypeGroupMetadataPayloadSent {
			continue
		}

		m := &protocoltypes.GroupMetadataPayloadSent{}
		if err := m.Unmarshal(evt.Event); err != nil {
			continue
		}

		out = append(out, m)
	}

	return out
}
