package models

import types "github.com/prysmaticlabs/eth2-types"

// SignRequestSlot struct fir sign req slot
type SignRequestSlot struct {
	Slot types.Slot
}

// isSignRequestObject implementation interface
func (m *SignRequestSlot) isSignRequestObject() {}
