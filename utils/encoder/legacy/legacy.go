package legacy

import (
	"encoding/json"
	"errors"

	"github.com/bloxapp/key-vault/utils/encoder/legacy/ethereum_validator_accounts_v2"

	oldPrysm "github.com/bloxapp/key-vault/utils/encoder/legacy/eth"

	validatorpb "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1/validator-client"

	newPrysm "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"
)

// Legacy is an old encoder type used by Prysm before they moved to SSZ completely.
// Older versions of KeyVault use this encoding.
type Legacy struct {
}

func New() *Legacy {
	return &Legacy{}
}

// Encode takes a new/ old prysm object, encodes it to the legacy encoding
func (l *Legacy) Encode(obj interface{}) ([]byte, error) {
	switch t := obj.(type) {
	case *newPrysm.AttestationData:
		return oldPrysm.NewAttestationDataFromNewPrysm(t).Marshal()
	case oldPrysm.AttestationData:
		return t.Marshal()
	case *newPrysm.BeaconBlock:
		return oldPrysm.NewBeaconBlockFromNewPrysm(t).Marshal()
	case oldPrysm.BeaconBlock:
		return t.Marshal()
	case *validatorpb.SignRequest:
		toEncode, err := ethereum_validator_accounts_v2.NewSignRequestFromNewPrysm(t)
		if err != nil {
			return nil, err
		}
		return json.Marshal(toEncode)
	case ethereum_validator_accounts_v2.SignRequest:
		return json.Marshal(t)
	}
	return nil, errors.New("type not supported")
}

// Decode takes an old legacy encoding and populates a old/ new prysm object
func (l *Legacy) Decode(data []byte, v interface{}) error {
	switch t := v.(type) {
	case *newPrysm.AttestationData:
		return oldPrysm.NewAttestationDataFromNewPrysm(t).Unmarshal(data)
	case oldPrysm.AttestationData:
		return t.Unmarshal(data)
	case *newPrysm.BeaconBlock:
		return oldPrysm.NewBeaconBlockFromNewPrysm(t).Unmarshal(data)
	case *oldPrysm.BeaconBlock:
		return t.Unmarshal(data)
	case *validatorpb.SignRequest:
		toDecode, err := ethereum_validator_accounts_v2.NewSignRequestFromNewPrysm(t)
		if err != nil {
			return err
		}
		return toDecode.Unmarshal(data)
	case ethereum_validator_accounts_v2.SignRequest:
		return json.Unmarshal(data, t)
	}
	return errors.New("type not supported")
}
