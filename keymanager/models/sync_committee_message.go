package models

import types "github.com/prysmaticlabs/eth2-types"

type SignRequest_SyncCommitteeMessage struct {
	Root types.SSZBytes
}

func (m *SignRequest_SyncCommitteeMessage) isSignRequest_Object() {}
