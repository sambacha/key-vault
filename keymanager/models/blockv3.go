package models

import eth "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"

// SignRequestBlockV3 struct
type SignRequestBlockV3 struct {
	BlockV3 *eth.BeaconBlockMerge
}

// isSignRequestObject implement func
func (m *SignRequestBlockV3) isSignRequestObject() {}
