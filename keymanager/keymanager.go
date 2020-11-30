package keymanager

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	ethpb "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
	validatorpb "github.com/prysmaticlabs/prysm/proto/validator/accounts/v2"
	"github.com/prysmaticlabs/prysm/shared/bls"
	"github.com/prysmaticlabs/prysm/validator/keymanager"
	"github.com/sirupsen/logrus"

	"github.com/bloxapp/key-vault/utils/bytex"
	"github.com/bloxapp/key-vault/utils/endpoint"
	"github.com/bloxapp/key-vault/utils/httpex"
)

// To make sure V2 implements keymanager.IKeymanager interface
var _ keymanager.IKeymanager = &KeyManager{}

// Predefined errors
var (
	ErrLocationMissing    = NewGenericErrorMessage("wallet location is required")
	ErrTokenMissing       = NewGenericErrorMessage("wallet access token is required")
	ErrPubKeyMissing      = NewGenericErrorMessage("wallet public key is required")
	ErrUnsupportedSigning = NewGenericErrorWithMessage("remote HTTP key manager does not support such signing method")
	ErrNoSuchKey          = NewGenericErrorWithMessage("no such key")
)

// KeyManager is a key manager that accesses a remote vault wallet daemon through HTTP connection.
type KeyManager struct {
	remoteAddress string
	accessToken   string
	originPubKey  string
	pubKey        [48]byte
	network       string
	httpClient    *http.Client

	log *logrus.Entry
}

// NewKeyManager is the constructor of KeyManager.
func NewKeyManager(log *logrus.Entry, opts *Config) (*KeyManager, error) {
	if len(opts.Location) == 0 {
		return nil, ErrLocationMissing
	}
	if len(opts.AccessToken) == 0 {
		return nil, ErrTokenMissing
	}
	if len(opts.PubKey) == 0 {
		return nil, ErrPubKeyMissing
	}

	// Decode public key
	decodedPubKey, err := hex.DecodeString(opts.PubKey)
	if err != nil {
		return nil, NewGenericError(err, "failed to hex decode public key '%s'", opts.PubKey)
	}

	return &KeyManager{
		remoteAddress: opts.Location,
		accessToken:   opts.AccessToken,
		originPubKey:  opts.PubKey,
		pubKey:        bytex.ToBytes48(decodedPubKey),
		network:       opts.Network,
		httpClient: httpex.CreateClient(log, func(resp *http.Response, err error, numTries int) (*http.Response, error) {
			if err == nil {
				return resp, nil
			}

			fields := logrus.Fields{}
			if resp != nil {
				fields["status_code"] = resp.StatusCode

				if resp.Body != nil {
					defer resp.Body.Close()

					respBody, err := ioutil.ReadAll(resp.Body)
					if err != nil {
						return resp, err
					}
					fields["response_body"] = string(respBody)
				}
			}

			log.WithError(err).WithFields(fields).Error("failed to send request to key manager")

			return resp, fmt.Errorf("giving up after %d attempt(s): %s", numTries, err)
		}),
		log: log,
	}, nil
}

// FetchValidatingPublicKeys implements KeyManager-v2 interface.
func (km *KeyManager) FetchValidatingPublicKeys(_ context.Context) ([][48]byte, error) {
	return [][48]byte{km.pubKey}, nil
}

// FetchAllValidatingPublicKeys implements KeyManager-v2 interface.
func (km *KeyManager) FetchAllValidatingPublicKeys(_ context.Context) ([][48]byte, error) {
	return [][48]byte{km.pubKey}, nil
}

// Sign implements IKeymanager interface.
func (km *KeyManager) Sign(_ context.Context, req *validatorpb.SignRequest) (bls.Signature, error) {
	if bytex.ToBytes48(req.GetPublicKey()) != km.pubKey {
		return nil, ErrNoSuchKey
	}

	domain := bytex.ToBytes32(req.GetSignatureDomain())
	switch data := req.GetObject().(type) {
	case *validatorpb.SignRequest_Block:
		return km.SignProposal(req.GetSignatureDomain(), &ethpb.BeaconBlockHeader{
			Slot:          data.Block.GetSlot(),
			ProposerIndex: data.Block.GetProposerIndex(),
			StateRoot:     data.Block.GetStateRoot(),
			ParentRoot:    data.Block.GetParentRoot(),
			BodyRoot:      req.GetSigningRoot(),
		})
	case *validatorpb.SignRequest_AttestationData:
		return km.SignAttestation(req.GetSignatureDomain(), data.AttestationData)
	case *validatorpb.SignRequest_AggregateAttestationAndProof:
		return km.SignGeneric(req.GetSigningRoot(), domain)
	case *validatorpb.SignRequest_Slot:
		return km.SignGeneric(req.GetSigningRoot(), domain)
	case *validatorpb.SignRequest_Epoch:
		return km.SignGeneric(req.GetSigningRoot(), domain)
	default:
		return nil, ErrUnsupportedSigning
	}
}

// sendRequest implements the logic to work with HTTP requests.
func (km *KeyManager) sendRequest(method, path string, reqBody []byte, respBody interface{}) error {
	endpoint := km.remoteAddress + endpoint.Build(km.network, path)

	// Prepare a new request
	req, err := http.NewRequest(method, endpoint, bytes.NewBuffer(reqBody))
	if err != nil {
		return NewGenericError(err, "failed to create HTTP request")
	}

	// Pass auth token.
	req.Header.Set("Authorization", "Bearer "+km.accessToken)
	req.Header.Set("Content-Type", "application/json")

	// Send request.
	resp, err := km.httpClient.Do(req)
	if err != nil {
		return NewGenericError(err, "failed to send HTTP request")
	}
	defer resp.Body.Close()

	// Check status code. Must be 200.
	if resp.StatusCode != http.StatusOK {
		responseBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			km.log.WithError(err).Error("failed to read error response body")
		}

		return NewHTTPRequestError(endpoint, resp.StatusCode, responseBody, "unexpected status code")
	}

	// Retrieve response body
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return NewGenericError(err, "failed to read response body")
	}

	// Read response body into the given object.
	if err := json.Unmarshal(responseBody, &respBody); err != nil {
		return NewGenericError(err, "failed to decode response body")
	}

	return nil
}
