package models

import eth "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"

// SignRequestBlockV2 struct
type SignRequestBlockV2 struct {
	BlockV2 *eth.BeaconBlockAltair
}

// isSignRequestObject implement func
func (m *SignRequestBlockV2) isSignRequestObject() {}
