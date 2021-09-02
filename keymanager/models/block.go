package models

import eth "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"

type SignRequest_Block struct {
	Block *eth.BeaconBlock
}

func (m *SignRequest_Block) isSignRequest_Object() {}
