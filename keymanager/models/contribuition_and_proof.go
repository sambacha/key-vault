package models

import eth "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"

// SignRequestContributionAndProof struct
type SignRequestContributionAndProof struct {
	ContributionAndProof *eth.ContributionAndProof
}

// isSignRequestObject implement func
func (m *SignRequestContributionAndProof) isSignRequestObject() {}
