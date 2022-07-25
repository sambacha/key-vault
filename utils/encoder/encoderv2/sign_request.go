package encoderv2

import (
	"encoding/json"
	"errors"
	"reflect"

	"github.com/bloxapp/key-vault/keymanager/models"

	types "github.com/prysmaticlabs/prysm/consensus-types/primitives"

	eth "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"
)

type signRequestEncoded struct {
	PublicKey       []byte `protobuf:"bytes,1,opt,name=public_key,json=publicKey,proto3" json:"public_key,omitempty"`
	SigningRoot     []byte `protobuf:"bytes,2,opt,name=signing_root,json=signingRoot,proto3" json:"signing_root,omitempty"`
	SignatureDomain []byte `protobuf:"bytes,3,opt,name=signature_domain,json=signatureDomain,proto3" json:"signature_domain,omitempty"`
	Data            []byte
	ObjectType      string
}

func encodeSignRequest(sr *models.SignRequest) ([]byte, error) {
	toEncode := signRequestEncoded{
		PublicKey:       sr.PublicKey,
		SigningRoot:     sr.SigningRoot,
		SignatureDomain: sr.SignatureDomain,
	}

	if sr.Object == nil {
		return json.Marshal(toEncode)
	}

	switch t := sr.Object.(type) {
	case *models.SignRequestAttestationData:
		byts, err := t.AttestationData.MarshalSSZ()
		if err != nil {
			return nil, err
		}
		toEncode.Data = byts
		toEncode.ObjectType = reflect.TypeOf(t).String()
	case *models.SignRequestBlock:
		byts, err := t.Block.MarshalSSZ()
		if err != nil {
			return nil, err
		}
		toEncode.Data = byts
		toEncode.ObjectType = reflect.TypeOf(t).String()
	case *models.SignRequestBlockV2:
		byts, err := t.BlockV2.MarshalSSZ()
		if err != nil {
			return nil, err
		}
		toEncode.Data = byts
		toEncode.ObjectType = reflect.TypeOf(t).String()
	case *models.SignRequestBlockV3:
		byts, err := t.BlockV3.MarshalSSZ()
		if err != nil {
			return nil, err
		}
		toEncode.Data = byts
		toEncode.ObjectType = reflect.TypeOf(t).String()
	case *models.SignRequestBlindedBlockV3:
		byts, err := t.BlindedBlockV3.MarshalSSZ()
		if err != nil {
			return nil, err
		}
		toEncode.Data = byts
		toEncode.ObjectType = reflect.TypeOf(t).String()
	case *models.SignRequestAggregateAttestationAndProof:
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
	case *models.SignRequestEpoch:
		byts, err := t.Epoch.MarshalSSZ()
		if err != nil {
			return nil, err
		}
		toEncode.Data = byts
		toEncode.ObjectType = reflect.TypeOf(t).String()
	case *models.SignRequestSlot:
		byts, err := t.Slot.MarshalSSZ()
		if err != nil {
			return nil, err
		}
		toEncode.Data = byts
		toEncode.ObjectType = reflect.TypeOf(t).String()
	case *models.SignRequestSyncCommitteeMessage:
		toEncode.Data = t.Root
		toEncode.ObjectType = reflect.TypeOf(t).String()
	case *models.SignRequestSyncAggregatorSelectionData:
		byts, err := t.SyncAggregatorSelectionData.MarshalSSZ()
		if err != nil {
			return nil, err
		}
		toEncode.Data = byts
		toEncode.ObjectType = reflect.TypeOf(t).String()
	case *models.SignRequestContributionAndProof:
		byts, err := t.ContributionAndProof.MarshalSSZ()
		if err != nil {
			return nil, err
		}
		toEncode.Data = byts
		toEncode.ObjectType = reflect.TypeOf(t).String()
	case *models.SignRequestRegistration:
		byts, err := t.Registration.MarshalSSZ()
		if err != nil {
			return nil, err
		}
		toEncode.Data = byts
		toEncode.ObjectType = reflect.TypeOf(t).String()
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
	case "*models.SignRequestAttestationData":
		data := &eth.AttestationData{}
		if err := data.UnmarshalSSZ(toDecode.Data); err != nil {
			return err
		}
		sr.Object = &models.SignRequestAttestationData{AttestationData: data}
	case "*models.SignRequestBlock":
		data := &eth.BeaconBlock{}
		if err := data.UnmarshalSSZ(toDecode.Data); err != nil {
			return err
		}
		sr.Object = &models.SignRequestBlock{Block: data}
	case "*models.SignRequestBlockV2":
		data := &eth.BeaconBlockAltair{}
		if err := data.UnmarshalSSZ(toDecode.Data); err != nil {
			return err
		}
		sr.Object = &models.SignRequestBlockV2{BlockV2: data}
	case "*models.SignRequestBlockV3":
		data := &eth.BeaconBlockBellatrix{}
		if err := data.UnmarshalSSZ(toDecode.Data); err != nil {
			return err
		}
		sr.Object = &models.SignRequestBlockV3{BlockV3: data}
	case "*models.SignRequestBlindedBlockV3":
		data := &eth.BlindedBeaconBlockBellatrix{}
		if err := data.UnmarshalSSZ(toDecode.Data); err != nil {
			return err
		}
		sr.Object = &models.SignRequestBlindedBlockV3{BlindedBlockV3: data}
	case "*models.SignRequestSlot":
		data := types.Slot(1)
		if err := data.UnmarshalSSZ(toDecode.Data); err != nil {
			return err
		}
		sr.Object = &models.SignRequestSlot{Slot: data}
	case "*models.SignRequestEpoch":
		data := types.Epoch(1)
		if err := data.UnmarshalSSZ(toDecode.Data); err != nil {
			return err
		}
		sr.Object = &models.SignRequestEpoch{Epoch: data}
	case "*models.SignRequestAggregateAttestationAndProof":
		data := &eth.AggregateAttestationAndProof{}
		if err := data.UnmarshalSSZ(toDecode.Data); err != nil {
			return err
		}
		sr.Object = &models.SignRequestAggregateAttestationAndProof{AggregateAttestationAndProof: data}
	case "*models.SignRequestSyncCommitteeMessage":
		sr.Object = &models.SignRequestSyncCommitteeMessage{Root: toDecode.Data}
	case "*models.SignRequestSyncAggregatorSelectionData":
		data := &eth.SyncAggregatorSelectionData{}
		if err := data.UnmarshalSSZ(toDecode.Data); err != nil {
			return err
		}
		sr.Object = &models.SignRequestSyncAggregatorSelectionData{SyncAggregatorSelectionData: data}
	case "*models.SignRequestContributionAndProof":
		data := &eth.ContributionAndProof{}
		if err := data.UnmarshalSSZ(toDecode.Data); err != nil {
			return err
		}
		sr.Object = &models.SignRequestContributionAndProof{ContributionAndProof: data}
	case "*models.SignRequestRegistration":
		data := &eth.ValidatorRegistrationV1{}
		if err := data.UnmarshalSSZ(toDecode.Data); err != nil {
			return err
		}
		sr.Object = &models.SignRequestRegistration{Registration: data}
	default:
		return errors.New("sign request unknown object type")
	}
	return nil
}
