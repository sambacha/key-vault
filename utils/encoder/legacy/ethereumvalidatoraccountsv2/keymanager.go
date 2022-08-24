package ethereumvalidatoraccountsv2

import (
	"errors"

	validatorpb "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1/validator-client"

	"github.com/bloxapp/key-vault/utils/encoder/legacy/eth"
)

// NewSignRequestFromNewPrysm creates new sign request
func NewSignRequestFromNewPrysm(newPrysm *validatorpb.SignRequest) (*SignRequest, error) {
	ret := &SignRequest{
		PublicKey:       newPrysm.PublicKey,
		SigningRoot:     newPrysm.SigningRoot,
		SignatureDomain: newPrysm.SignatureDomain,
	}

	if newPrysm.Object == nil {
		return ret, nil
	}

	switch t := newPrysm.Object.(type) {
	case *validatorpb.SignRequest_AttestationData:
		ret.Object = &SignRequest_AttestationData{AttestationData: eth.NewAttestationDataFromNewPrysm(t.AttestationData)}
	case *validatorpb.SignRequest_Block:
		ret.Object = &SignRequest_Block{Block: eth.NewBeaconBlockFromNewPrysm(t.Block)}
	case *validatorpb.SignRequest_Slot:
		ret.Object = &SignRequest_Slot{Slot: uint64(t.Slot)}
	case *validatorpb.SignRequest_Epoch:
		ret.Object = &SignRequest_Epoch{Epoch: uint64(t.Epoch)}
	case *validatorpb.SignRequest_AggregateAttestationAndProof:
		ret.Object = &SignRequest_AggregateAttestationAndProof{AggregateAttestationAndProof: eth.NewAggregationAndProofFromNewPrysm(t.AggregateAttestationAndProof)}
	default:
		return nil, errors.New("sign request type not supported")
	}

	return ret, nil
}
