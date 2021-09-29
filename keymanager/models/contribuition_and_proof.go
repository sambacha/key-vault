package models

import eth "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"

type SignRequest_ContributionAndProof struct {
	ContributionAndProof *eth.ContributionAndProof
}

func (m *SignRequest_ContributionAndProof) isSignRequest_Object() {}
