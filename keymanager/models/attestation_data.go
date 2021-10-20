package models

import eth "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"

// SignRequestAttestationData struct
type SignRequestAttestationData struct {
	AttestationData *eth.AttestationData
}

// isSignRequestObject implement func
func (m *SignRequestAttestationData) isSignRequestObject() {}
