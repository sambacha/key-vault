package models

import eth "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"

// SignRequestAggregateAttestationAndProof struct
type SignRequestAggregateAttestationAndProof struct {
	AggregateAttestationAndProof *eth.AggregateAttestationAndProof
}

// isSignRequestObject implement func
func (m *SignRequestAggregateAttestationAndProof) isSignRequestObject() {}
