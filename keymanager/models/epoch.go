package models

import types "github.com/prysmaticlabs/eth2-types"

type SignRequest_Epoch struct {
	Epoch types.Epoch
}

func (m *SignRequest_Epoch) isSignRequest_Object() {}
