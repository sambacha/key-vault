package tests

import (
	"encoding/hex"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"

	eth "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
	validatorpb "github.com/prysmaticlabs/prysm/proto/validator/accounts/v2"

	"github.com/bloxapp/eth2-key-manager/core"
	"github.com/stretchr/testify/require"

	"github.com/bloxapp/key-vault/e2e"
	"github.com/bloxapp/key-vault/e2e/shared"
)

// AttestationConcurrentSigning tests signing method concurrently.
type AttestationConcurrentSigning struct {
}

// Name returns the name of the test.
func (test *AttestationConcurrentSigning) Name() string {
	return "Test attestation concurrent signing"
}

// Run runs the test.
func (test *AttestationConcurrentSigning) Run(t *testing.T) {
	setup := e2e.Setup(t)

	// setup vault with db
	store := setup.UpdateStorage(t, core.PyrmontNetwork, true, core.HDWallet, nil)
	account := shared.RetrieveAccount(t, store)
	pubKey := account.ValidatorPublicKey()

	wg := &sync.WaitGroup{}
	signedCnt := int64(0)
	for i := uint64(0); i < 20; i++ {
		go test.runSlashableAttestation(t, &signedCnt, wg, setup, pubKey)
	}
	wg.Wait()
	require.EqualValues(t, 1, signedCnt)
}

// will return no error if trying to sign a slashable attestation will not work
func (test *AttestationConcurrentSigning) runSlashableAttestation(t *testing.T, cnt *int64, wg *sync.WaitGroup, setup *e2e.BaseSetup, pubKey []byte) {
	wg.Add(1)
	defer wg.Done()

	randomCommittee := func() uint64 {
		max := 1000
		min := 2
		return uint64(rand.Intn(max-min) + min)
	}

	att := &eth.AttestationData{
		Slot:            284115,
		CommitteeIndex:  randomCommittee(),
		BeaconBlockRoot: _byteArray32("7b5679277ca45ea74e1deebc9d3e8c0e7d6c570b3cfaf6884be144a81dac9a0e"),
		Source: &eth.Checkpoint{
			Epoch: 77,
			Root:  _byteArray32("7402fdc1ce16d449d637c34a172b349a12b2bae8d6d77e401006594d8057c33d"),
		},
		Target: &eth.Checkpoint{
			Epoch: 78,
			Root:  _byteArray32("17959acc370274756fa5e9fdd7e7adf17204f49cc8457e49438c42c4883cbfb0"),
		},
	}
	domain := _byteArray32("01000000f071c66c6561d0b939feb15f513a019d99a84bd85635221e3ad42dac")

	req, err := test.serializedReq(pubKey, nil, domain, att)
	require.NoError(t, err)
	_, err = setup.Sign("sign", req, core.PyrmontNetwork)
	if err == nil {
		atomic.AddInt64(cnt, 1)
	}
}

func (test *AttestationConcurrentSigning) serializedReq(pk, root, domain []byte, attestation *eth.AttestationData) (map[string]interface{}, error) {
	req := &validatorpb.SignRequest{
		PublicKey:       pk,
		SigningRoot:     root,
		SignatureDomain: domain,
		Object:          &validatorpb.SignRequest_AttestationData{AttestationData: attestation},
	}

	byts, err := req.Marshal()
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"sign_req": hex.EncodeToString(byts),
	}, nil
}
