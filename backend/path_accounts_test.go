package backend

import (
	"context"
	"sort"
	"testing"
	"time"

	"github.com/bloxapp/eth2-key-manager/core"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func getBackend(t *testing.T) (logical.Backend, logical.Storage) {
	config := &logical.BackendConfig{
		Logger:      logging.NewVaultLogger(log.Trace),
		System:      &logical.StaticSystemView{},
		StorageView: &logical.InmemStorage{},
		BackendUUID: "test",
	}

	b, err := Factory("test", logrus.New())(context.Background(), config)
	if err != nil {
		t.Fatalf("unable to create backend: %v", err)
	}

	// Wait for the upgrade to finish
	time.Sleep(time.Second)

	return b, config.StorageView
}

func setupBaseStorage(t *testing.T, req *logical.Request, configModifiers ...func(*Config)) {
	cfg := Config{
		Network:       core.PraterNetwork,
		FeeRecipients: FeeRecipients{"0x95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf": "0x6a3f3ee924a940ce0d795c5a41a817607e520520"},
	}
	for _, mod := range configModifiers {
		mod(&cfg)
	}
	entry, err := logical.StorageEntryJSON("config", cfg)
	require.NoError(t, err)
	require.NoError(t, req.Storage.Put(context.Background(), entry))
}

func TestAccountsList(t *testing.T) {
	b, _ := getBackend(t)

	t.Run("Successfully List Accounts", func(t *testing.T) {
		req := logical.TestRequest(t, logical.ListOperation, "accounts/")
		setupBaseStorage(t, req)

		// setup logical storage
		_, err := baseHashicorpStorage(context.Background(), req.Storage)
		require.NoError(t, err)

		res, err := b.HandleRequest(context.Background(), req)
		require.NoError(t, err)
		require.NotNil(t, res.Data["accounts"])
		require.Len(t, res.Data["accounts"], 1)
		require.Contains(t, res.Data["accounts"].([]map[string]string)[0]["name"], "account-")

		// make sure only the following fields are present to prevent accidental secret sharing
		keys := make([]string, 0)
		for k := range res.Data["accounts"].([]map[string]string)[0] {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		require.Equal(t, keys, []string{"id", "name", "validationPubKey", "withdrawalPubKey"})
	})
}
