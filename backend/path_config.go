package backend

import (
	"context"

	"github.com/bloxapp/eth2-key-manager/core"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/pkg/errors"
	field_params "github.com/prysmaticlabs/prysm/config/fieldparams"
)

// Endpoints patterns
const (
	// ConfigPattern is the path pattern for config endpoint
	ConfigPattern = "config"
)

// Config contains the configuration for each mount
type Config struct {
	Network       core.Network      `json:"network"`
	FeeRecipients map[string]string `json:"fee_recipients"`
}

func (c Config) Map() map[string]interface{} {
	return map[string]interface{}{
		"network":        c.Network,
		"fee_recipients": c.FeeRecipients,
	}
}

func configPaths(b *backend) []*framework.Path {
	return []*framework.Path{
		{
			Pattern: ConfigPattern,
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.CreateOperation: &framework.PathOperation{
					Callback: b.pathWriteConfig,
				},
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.pathWriteConfig,
				},
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.pathReadConfig,
				},
			},
			HelpSynopsis:    "Configure the Vault Ethereum plugin.",
			HelpDescription: "Configure the Vault Ethereum plugin.",
			Fields: map[string]*framework.FieldSchema{
				"network": {
					Type: framework.TypeString,
					Description: `Ethereum network - can be one of the following values:
					mainnet - MainNet Network
					prater - Prater Test Network`,
					AllowedValues: []interface{}{
						string(core.PraterNetwork),
						string(core.MainNetwork),
					},
				},
				"fee_recipients": {
					Type:        framework.TypeMap,
					Description: `Validator pubic keys and their associated fee recipient addresses.`,
				},
			},
		},
	}
}

// pathWriteConfig is the write config path handler
func (b *backend) pathWriteConfig(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	network := core.NetworkFromString(data.Get("network").(string))
	if network == "" {
		return nil, errors.New("invalid network provided")
	}

	normalizedRecipients := make(map[string]string)
	recipients, _ := data.Get("fee_recipients").(map[string]interface{})
	for key, value := range recipients {
		var normalizedKey string
		switch key {
		case "default":
			normalizedKey = "default"
		default:
			validatorPubkey, err := hexutil.Decode(key)
			if err != nil || len(validatorPubkey) != field_params.BLSPubkeyLength {
				return nil, errors.Wrap(err, "invalid fee_recipients provided")
			}
			normalizedKey = hexutil.Encode(validatorPubkey)
		}

		recipientAddrStr, _ := value.(string)
		recipientAddr, err := hexutil.Decode(recipientAddrStr)
		if err != nil || len(recipientAddr) != field_params.FeeRecipientLength {
			return nil, errors.Wrap(err, "invalid fee_recipients provided")
		}

		normalizedRecipients[normalizedKey] = hexutil.Encode(recipientAddr)
	}

	configBundle := Config{
		Network:       network,
		FeeRecipients: normalizedRecipients,
	}

	// Create storage entry
	entry, err := logical.StorageEntryJSON("config", configBundle)
	if err != nil {
		return nil, err
	}

	// Store config
	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}

	// Return the secret
	return &logical.Response{
		Data: configBundle.Map(),
	}, nil
}

// pathReadConfig is the read config path handler
func (b *backend) pathReadConfig(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	configBundle, err := b.readConfig(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	if configBundle == nil {
		return nil, nil
	}

	return &logical.Response{
		Data: configBundle.Map(),
	}, nil
}

// readConfig returns the configuration for this PluginBackend.
func (b *backend) readConfig(ctx context.Context, s logical.Storage) (*Config, error) {
	entry, err := s.Get(ctx, "config")
	if err != nil {
		return nil, err
	}

	if entry == nil {
		return nil, errors.Errorf("the plugin has not been configured yet")
	}

	var result Config
	if err := entry.DecodeJSON(&result); err != nil {
		return nil, errors.Wrap(err, "error reading configuration")
	}

	return &result, nil
}
