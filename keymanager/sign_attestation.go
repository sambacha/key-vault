package keymanager

import (
	"encoding/hex"
	"encoding/json"
	"net/http"

	ethpb "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
	"github.com/prysmaticlabs/prysm/shared/bls"

	"github.com/bloxapp/key-vault/backend"
)

// SignAttestation implements ProtectingKeyManager interface.
func (km *KeyManager) SignAttestation(domain []byte, data *ethpb.AttestationData) (bls.Signature, error) {
	// Prepare request body.
	req := SignAttestationRequest{
		PubKey:          km.originPubKey,
		Domain:          hex.EncodeToString(domain[:]),
		Slot:            data.GetSlot(),
		CommitteeIndex:  data.GetCommitteeIndex(),
		BeaconBlockRoot: hex.EncodeToString(data.GetBeaconBlockRoot()),
		SourceEpoch:     data.GetSource().GetEpoch(),
		SourceRoot:      hex.EncodeToString(data.GetSource().GetRoot()),
		TargetEpoch:     data.GetTarget().GetEpoch(),
		TargetRoot:      hex.EncodeToString(data.GetTarget().GetRoot()),
	}

	// Json encode the request body
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, NewGenericError(err, "failed to marshal request body")
	}

	// Send request.
	var resp SignResponse
	if err := km.sendRequest(http.MethodPost, backend.SignAttestationPattern, reqBody, &resp); err != nil {
		km.log.WithError(err).Error("failed to send sign attestation request")
		return nil, NewGenericError(err, "failed to send SignAttestation request to remote vault wallet")
	}

	// Signature is base64 encoded, so we have to decode that.
	decodedSignature, err := hex.DecodeString(resp.Data.Signature)
	if err != nil {
		return nil, NewGenericError(err, "failed to base64 decode")
	}

	// Get signature from bytes
	sig, err := bls.SignatureFromBytes(decodedSignature)
	if err != nil {
		return nil, NewGenericError(err, "failed to get BLS signature from bytes")
	}

	return sig, nil
}
