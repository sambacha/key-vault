package sign_request

import eth "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"

type SignRequest_AggregateAttestationAndProof struct {
	AggregateAttestationAndProof *eth.AggregateAttestationAndProof
}

func (m *SignRequest_AggregateAttestationAndProof) isSignRequest_Object() {}
