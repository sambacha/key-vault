package models

import eth "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"

type SignRequest_AttestationData struct {
	AttestationData *eth.AttestationData
}

func (m *SignRequest_AttestationData) isSignRequest_Object() {}
