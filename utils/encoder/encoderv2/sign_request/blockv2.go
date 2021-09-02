package sign_request

import eth "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"

type SignRequest_BlockV2 struct {
	BlockV2 *eth.BeaconBlockAltair
}

func (m *SignRequest_BlockV2) isSignRequest_Object() {}
