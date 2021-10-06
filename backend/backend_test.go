package backend

import (
	"os"
	"testing"

	vault "github.com/bloxapp/eth2-key-manager"
)

func TestMain(t *testing.M) {
	vault.InitCrypto()
	t.Run()
	os.Exit(0)
}
