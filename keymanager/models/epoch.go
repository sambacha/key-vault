package models

import types "github.com/prysmaticlabs/eth2-types"

// SignRequestEpoch struct
type SignRequestEpoch struct {
	Epoch types.Epoch
}

// isSignRequestObject implement func
func (m *SignRequestEpoch) isSignRequestObject() {}
