package replicationtypes

import (
	"encoding/base64"

	"berty.tech/weshnet/v2/pkg/protocoltypes"
)

func (m *ReplicatedGroup) ToGroup() (*protocoltypes.Group, error) {
	pk, err := base64.RawURLEncoding.DecodeString(m.PublicKey)
	if err != nil {
		return nil, err
	}

	signPub, err := base64.RawURLEncoding.DecodeString(m.SignPub)
	if err != nil {
		return nil, err
	}

	linkKey, err := base64.RawURLEncoding.DecodeString(m.LinkKey)
	if err != nil {
		return nil, err
	}

	return &protocoltypes.Group{
		PublicKey: pk,
		SignPub:   signPub,
		LinkKey:   linkKey,
	}, nil
}
