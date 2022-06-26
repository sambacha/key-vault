package tests

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/bloxapp/eth2-key-manager/core"
	"github.com/stretchr/testify/require"

	"github.com/bloxapp/key-vault/backend"
	"github.com/bloxapp/key-vault/e2e"
)

// ConfigUpdate tests update config endpoint.
type ConfigUpdate struct {
}

// Name returns the name of the test.
func (test *ConfigUpdate) Name() string {
	return "Test update config endpoint"
}

// Run runs the test.
func (test *ConfigUpdate) Run(t *testing.T) {
	setup := e2e.Setup(t)

	test.testValidUpdate(t, setup, makeConfig(map[string]string{
		"0x803fdf24c80ab6d1f78e93ae29791efe70280f35e05eebb4d8d6c381674346c6a346ec89790bdd8ea5caa4b18bebcf2c": "0x96759b16ab543551c06b5e6b0e4f887c9401b654",
		"default": "0x96759b16ab543551c06b5e6b0e4f887c9401b654",
	}))

	test.testInvalidUpdate(t, setup, makeConfig(map[string]string{
		"0803fdf24c80ab6d1f78e93ae29791efe70280f35e05eebb4d8d6c381674346c6a346ec89790bdd8ea5caa4b18bebcf2c": "0x96759b16ab543551c06b5e6b0e4f887c9401b654",
		"default": "0x96759b16ab543551c06b5e6b0e4f887c9401b654",
	}))

	test.testInvalidUpdate(t, setup, makeConfig(map[string]string{
		"0x803fdf24c80ab6d1f78e93ae29791efe70280f35e05eebb4d8d6c381674346c6a346ec89790bdd8ea5caa4b18bebcf2c": "096759b16ab543551c06b5e6b0e4f887c9401b654",
		"default": "0x96759b16ab543551c06b5e6b0e4f887c9401b654",
	}))

	test.testInvalidUpdate(t, setup, makeConfig(map[string]string{
		"0x803fdf24c80ab6d1f78e93ae29791efe70280f35e05eebb4d8d6c381674346c6a346ec89790bdd8ea5caa4b18bebcf2c": "0x96759b16ab543551c06b5e6b0e4f887c9401b654",
		"default": "096759b16ab543551c06b5e6b0e4f887c9401b654",
	}))
}

func (test *ConfigUpdate) testValidUpdate(t *testing.T, setup *e2e.BaseSetup, configUpdate backend.Config) {
	// update the config
	_, statusCode := setup.UpdateConfig(t, core.PraterNetwork, configUpdate)
	require.Equal(t, http.StatusOK, statusCode)

	// read the updated config
	configActual, err := test.readAndParseConfig(t, setup)
	require.NoError(t, err)

	// compare the config
	require.EqualValues(t, configUpdate.FeeRecipients, configActual.FeeRecipients, "config didn't update")
}

func (test *ConfigUpdate) testInvalidUpdate(t *testing.T, setup *e2e.BaseSetup, configUpdate backend.Config) {
	// read the current config
	// to ensure it didn't change after the failed update
	configOld, err := test.readAndParseConfig(t, setup)
	require.NoError(t, err)

	// fail if the config update is identical to the old one
	require.NotEqualValues(t, configOld, configUpdate)

	// update the config
	_, statusCode := setup.UpdateConfig(t, core.PraterNetwork, configUpdate)
	require.Equal(t, statusCode, 500)

	// read the updated config
	configActual, err := test.readAndParseConfig(t, setup)
	require.NoError(t, err)

	// ensure the config didn't change
	require.EqualValues(t, configOld, configActual, "config updated even though it should not")
}

func (test *ConfigUpdate) readAndParseConfig(t *testing.T, setup *e2e.BaseSetup) (backend.Config, error) {
	// read the config
	configBytes, statusCode := setup.ReadConfig(t, core.PraterNetwork)
	require.Equal(t, http.StatusOK, statusCode)
	var config configModel
	err := json.Unmarshal(configBytes, &config)
	return config.Data, err
}

func makeConfig(feeRecipients map[string]string) backend.Config {
	return backend.Config{
		Network:       "prater",
		FeeRecipients: feeRecipients,
	}
}
