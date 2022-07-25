package models

import types "github.com/prysmaticlabs/prysm/consensus-types/primitives"

// SignRequestEpoch struct
type SignRequestEpoch struct {
	Epoch types.Epoch
}

// isSignRequestObject implement func
func (m *SignRequestEpoch) isSignRequestObject() {}
