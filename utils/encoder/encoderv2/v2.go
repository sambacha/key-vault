package encoderv2

import (
	"errors"

	validatorpb "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1/validator-client"

	newPrysm "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"
)

// V2 encoder
type V2 struct {
}

func New() *V2 {
	return &V2{}
}

// Encode
func (l *V2) Encode(obj interface{}) ([]byte, error) {
	switch t := obj.(type) {
	case *newPrysm.AttestationData:
		return t.MarshalSSZ()
	case *newPrysm.BeaconBlock:
		return t.MarshalSSZ()
	case *validatorpb.SignRequest:
		return encodeSignReuqest(t)
	}
	return nil, errors.New("type not supported")
}

// Decode
func (l *V2) Decode(data []byte, v interface{}) error {
	switch t := v.(type) {
	case *newPrysm.AttestationData:
		return t.UnmarshalSSZ(data)
	case *newPrysm.BeaconBlock:
		return t.UnmarshalSSZ(data)
	case *validatorpb.SignRequest:
		return decodeSignRequest(data, t)
	}
	return errors.New("type not supported")
}
