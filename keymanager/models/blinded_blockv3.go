package models

import eth "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"

// SignRequestBlindedBlockV3 struct
type SignRequestBlindedBlockV3 struct {
	BlindedBlockV3 *eth.BlindedBeaconBlockBellatrix
}

// isSignRequestObject implement func
func (m *SignRequestBlindedBlockV3) isSignRequestObject() {}
