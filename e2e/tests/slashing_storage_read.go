package tests

import (
	"encoding/hex"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/bloxapp/key-vault/keymanager/models"

	eth "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"

	"github.com/bloxapp/eth2-key-manager/core"

	"github.com/stretchr/testify/require"

	"github.com/bloxapp/key-vault/e2e"
	"github.com/bloxapp/key-vault/e2e/shared"
)

type slashingHistoryModel struct {
	Data map[string]string `json:"data"`
}

// SlashingStorageRead tests slashing storage reading endpoint.
type SlashingStorageRead struct {
}

// Name returns the name of the test.
func (test *SlashingStorageRead) Name() string {
	return "Test slashing storage read"
}

// Run run the test.
func (test *SlashingStorageRead) Run(t *testing.T) {
	setup := e2e.Setup(t)

	// setup vault with db
	storage := setup.UpdateStorage(t, core.PraterNetwork, true, core.HDWallet, nil)
	account := shared.RetrieveAccount(t, storage)
	require.NotNil(t, account)
	pubKey := account.ValidatorPublicKey()

	// Send sign proposal request
	blk := &eth.BeaconBlock{
		Slot:          78,
		ProposerIndex: 1010,
		ParentRoot:    _byteArray32("7b5679277ca45ea74e1deebc9d3e8c0e7d6c570b3cfaf6884be144a81dac9a0e"),
		StateRoot:     _byteArray32("7402fdc1ce16d449d637c34a172b349a12b2bae8d6d77e401006594d8057c33d"),
		Body:          &eth.BeaconBlockBody{},
	}
	domain := _byteArray32("01000000f071c66c6561d0b939feb15f513a019d99a84bd85635221e3ad42dac")
	req, err := test.serializedReq(pubKey, nil, domain, blk)
	require.NoError(t, err)

	_, err = setup.Sign("sign", req, core.PraterNetwork)

	// Read slashing storage
	storageBytes, statusCode := setup.ReadSlashingStorage(t, core.PraterNetwork)
	require.Equal(t, http.StatusOK, statusCode, string(storageBytes))

	var slashingHistory slashingHistoryModel
	err = json.Unmarshal(storageBytes, &slashingHistory)
	require.NoError(t, err)

	pubKeyHistory, ok := slashingHistory.Data[hex.EncodeToString(pubKey)]
	require.True(t, ok)
	require.NotEmpty(t, pubKeyHistory)
}

func (test *SlashingStorageRead) serializedReq(pk, root, domain []byte, blk *eth.BeaconBlock) (map[string]interface{}, error) {
	req := &models.SignRequest{
		PublicKey:       pk,
		SigningRoot:     root,
		SignatureDomain: domain,
		Object:          &models.SignRequest_Block{Block: blk},
	}

	byts, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"sign_req": hex.EncodeToString(byts),
	}, nil
}
