package sign_request

import types "github.com/prysmaticlabs/eth2-types"

type SignRequest_Slot struct {
	Slot types.Slot
}

func (m *SignRequest_Slot) isSignRequest_Object() {}
