package encoderv2

import (
	"encoding/json"
	"errors"
	"reflect"

	"github.com/bloxapp/key-vault/keymanager/models"

	types "github.com/prysmaticlabs/eth2-types"

	eth "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"
)

type signRequestEncoded struct {
	PublicKey       []byte `protobuf:"bytes,1,opt,name=public_key,json=publicKey,proto3" json:"public_key,omitempty"`
	SigningRoot     []byte `protobuf:"bytes,2,opt,name=signing_root,json=signingRoot,proto3" json:"signing_root,omitempty"`
	SignatureDomain []byte `protobuf:"bytes,3,opt,name=signature_domain,json=signatureDomain,proto3" json:"signature_domain,omitempty"`
	Data            []byte
	ObjectType      string
}

func encodeSignReuqest(sr *models.SignRequest) ([]byte, error) {
	toEncode := signRequestEncoded{
		PublicKey:       sr.PublicKey,
		SigningRoot:     sr.SigningRoot,
		SignatureDomain: sr.SignatureDomain,
	}

	if sr.Object == nil {
		return json.Marshal(toEncode)
	}

	switch t := sr.Object.(type) {
	case *models.SignRequest_AttestationData:
		byts, err := t.AttestationData.MarshalSSZ()
		if err != nil {
			return nil, err
		}
		toEncode.Data = byts
		toEncode.ObjectType = reflect.TypeOf(t).String()
		break
	case *models.SignRequest_Block:
		byts, err := t.Block.MarshalSSZ()
		if err != nil {
			return nil, err
		}
		toEncode.Data = byts
		toEncode.ObjectType = reflect.TypeOf(t).String()
		break
	case *models.SignRequest_BlockV2:
		byts, err := t.BlockV2.MarshalSSZ()
		if err != nil {
			return nil, err
		}
		toEncode.Data = byts
		toEncode.ObjectType = reflect.TypeOf(t).String()
		break
	case *models.SignRequest_AggregateAttestationAndProof:
		byts, err := t.AggregateAttestationAndProof.MarshalSSZ()
		if err != nil {
			return nil, err
		}
		toEncode.Data = byts
		toEncode.ObjectType = reflect.TypeOf(t).String()

		ddd := &eth.AggregateAttestationAndProof{}
		if err := ddd.UnmarshalSSZ(byts); err != nil {
			return nil, err
		}

		break
	case *models.SignRequest_Epoch:
		byts, err := t.Epoch.MarshalSSZ()
		if err != nil {
			return nil, err
		}
		toEncode.Data = byts
		toEncode.ObjectType = reflect.TypeOf(t).String()
		break
	case *models.SignRequest_Slot:
		byts, err := t.Slot.MarshalSSZ()
		if err != nil {
			return nil, err
		}
		toEncode.Data = byts
		toEncode.ObjectType = reflect.TypeOf(t).String()
		break
	case *models.SignRequest_SyncCommitteeMessage:
		toEncode.Data = t.Root
		toEncode.ObjectType = reflect.TypeOf(t).String()
		break
	case *models.SignRequest_SyncAggregatorSelectionData:
		byts, err := t.SyncAggregatorSelectionData.MarshalSSZ()
		if err != nil {
			return nil, err
		}
		toEncode.Data = byts
		toEncode.ObjectType = reflect.TypeOf(t).String()
		break
	case *models.SignRequest_ContributionAndProof:
		byts, err := t.ContributionAndProof.MarshalSSZ()
		if err != nil {
			return nil, err
		}
		toEncode.Data = byts
		toEncode.ObjectType = reflect.TypeOf(t).String()
		break
	default:
		return nil, errors.New("sign request unknown object type")
	}

	return json.Marshal(toEncode)
}

func decodeSignRequest(data []byte, sr *models.SignRequest) error {
	toDecode := &signRequestEncoded{}
	if err := json.Unmarshal(data, toDecode); err != nil {
		return err
	}

	sr.PublicKey = toDecode.PublicKey
	sr.SignatureDomain = toDecode.SignatureDomain
	sr.SigningRoot = toDecode.SigningRoot

	if toDecode.Data == nil {
		return nil
	}

	switch toDecode.ObjectType {
	case "*sign_request.SignRequest_AttestationData":
		data := &eth.AttestationData{}
		if err := data.UnmarshalSSZ(toDecode.Data); err != nil {
			return err
		}
		sr.Object = &models.SignRequest_AttestationData{AttestationData: data}
		break
	case "*sign_request.SignRequest_Block":
		data := &eth.BeaconBlock{}
		if err := data.UnmarshalSSZ(toDecode.Data); err != nil {
			return err
		}
		sr.Object = &models.SignRequest_Block{Block: data}
		break
	case "*sign_request.SignRequest_BlockV2":
		data := &eth.BeaconBlockAltair{}
		if err := data.UnmarshalSSZ(toDecode.Data); err != nil {
			return err
		}
		sr.Object = &models.SignRequest_BlockV2{BlockV2: data}
		break
	case "*sign_request.SignRequest_Slot":
		data := types.Slot(1)
		if err := data.UnmarshalSSZ(toDecode.Data); err != nil {
			return err
		}
		sr.Object = &models.SignRequest_Slot{Slot: data}
		break
	case "*sign_request.SignRequest_Epoch":
		data := types.Epoch(1)
		if err := data.UnmarshalSSZ(toDecode.Data); err != nil {
			return err
		}
		sr.Object = &models.SignRequest_Epoch{Epoch: data}
		break
	case "*sign_request.SignRequest_AggregateAttestationAndProof":
		data := &eth.AggregateAttestationAndProof{}
		if err := data.UnmarshalSSZ(toDecode.Data); err != nil {
			return err
		}
		sr.Object = &models.SignRequest_AggregateAttestationAndProof{AggregateAttestationAndProof: data}
		break
	case "*sign_request.SignRequest_SyncCommitteeMessage":
		sr.Object = &models.SignRequest_SyncCommitteeMessage{Root: toDecode.Data}
	case "*sign_request.SignRequest_SyncAggregatorSelectionData":
		data := &eth.SyncAggregatorSelectionData{}
		if err := data.UnmarshalSSZ(toDecode.Data); err != nil {
			return err
		}
		sr.Object = &models.SignRequest_SyncAggregatorSelectionData{SyncAggregatorSelectionData: data}
		break
	case "*sign_request.SignRequest_ContributionAndProof":
		data := &eth.ContributionAndProof{}
		if err := data.UnmarshalSSZ(toDecode.Data); err != nil {
			return err
		}
		sr.Object = &models.SignRequest_ContributionAndProof{ContributionAndProof: data}
		break
	default:
		return errors.New("sign request unknown object type")
	}
	return nil
}
