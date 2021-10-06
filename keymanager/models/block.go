package models

import eth "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"

// SignRequestBlock struct
type SignRequestBlock struct {
	Block *eth.BeaconBlock
}

// isSignRequestObject implement func
func (m *SignRequestBlock) isSignRequestObject() {}
