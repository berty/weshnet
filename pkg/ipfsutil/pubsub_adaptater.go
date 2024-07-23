package ipfsutil

import (
	coreiface "github.com/ipfs/kubo/core/coreiface"
)

type pubsubCoreAPIAdapter struct {
	coreiface.PubSubAPI

	coreiface.CoreAPI
}

func (ps *pubsubCoreAPIAdapter) PubSub() coreiface.PubSubAPI {
	return ps.PubSubAPI
}

func InjectPubSubAPI(api coreiface.CoreAPI, ps coreiface.PubSubAPI) coreiface.CoreAPI {
	return &pubsubCoreAPIAdapter{
		PubSubAPI: ps,
		CoreAPI:   api,
	}
}

type pubsubExtendedCoreAPIAdapter struct {
	coreiface.PubSubAPI

	ExtendedCoreAPI
}

func (ps *pubsubExtendedCoreAPIAdapter) PubSub() coreiface.PubSubAPI {
	return ps.PubSubAPI
}

func InjectPubSubCoreAPIExtendedAdapter(exapi ExtendedCoreAPI, ps coreiface.PubSubAPI) ExtendedCoreAPI {
	return &pubsubExtendedCoreAPIAdapter{
		PubSubAPI:       ps,
		ExtendedCoreAPI: exapi,
	}
}
