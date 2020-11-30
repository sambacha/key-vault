package keymanager

import (
	"encoding/hex"
	"encoding/json"
	"net/http"

	ethpb "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
	"github.com/prysmaticlabs/prysm/shared/bls"

	"github.com/bloxapp/key-vault/backend"
)

// SignProposal implements ProtectingKeyManager interface.
func (km *KeyManager) SignProposal(domain []byte, data *ethpb.BeaconBlockHeader) (bls.Signature, error) {
	// Prepare request body.
	req := SignProposalRequest{
		PubKey:        km.originPubKey,
		Domain:        hex.EncodeToString(domain[:]),
		Slot:          data.GetSlot(),
		ProposerIndex: data.GetProposerIndex(),
		ParentRoot:    hex.EncodeToString(data.GetParentRoot()),
		StateRoot:     hex.EncodeToString(data.GetStateRoot()),
		BodyRoot:      hex.EncodeToString(data.GetBodyRoot()),
	}

	// Json encode the request body
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, NewGenericError(err, "failed to marshal request body")
	}

	// Send request.
	var resp SignResponse
	if err := km.sendRequest(http.MethodPost, backend.SignProposalPattern, reqBody, &resp); err != nil {
		km.log.WithError(err).Error("failed to send sign proposal request")
		return nil, NewGenericError(err, "failed to send SignProposal request to remote vault wallet")
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
