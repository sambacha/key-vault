package keymanager_test

import (
	"context"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"

	"github.com/bloxapp/key-vault/keymanager"
	"github.com/bloxapp/key-vault/utils/bytex"
)

var DefaultAccountPublicKey = "965586b5d05c851873f26cb736ed42de96591674772576e7b43848cd7a5c2827a5c5228034fdd55be0e9dc0f0cbc91d7"
var DefaultAccessToken = "supersecureaccesstoken"

func TestNewKeyManager(t *testing.T) {
	entry := logrus.NewEntry(logrus.New())
	type args struct {
		log  *logrus.Entry
		opts *keymanager.Config
	}
	tests := []struct {
		name    string
		args    args
		want    *keymanager.KeyManager
		wantErr bool
	}{
		{
			name: "empty location",
			args: args{
				log: entry,
				opts: &keymanager.Config{
					AccessToken: "AccessToken",
					PubKey:      DefaultAccountPublicKey,
					Network:     "Network",
				},
			},
			wantErr: true,
		},
		{
			name: "empty access token",
			args: args{
				log: entry,
				opts: &keymanager.Config{
					Location: "Location",
					PubKey:   DefaultAccountPublicKey,
					Network:  "Network",
				},
			},
			wantErr: true,
		},
		{
			name: "empty public key",
			args: args{
				log: entry,
				opts: &keymanager.Config{
					Location:    "Location",
					AccessToken: "AccessToken",
					Network:     "Network",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid public key",
			args: args{
				log: entry,
				opts: &keymanager.Config{
					Location:    "Location",
					AccessToken: "AccessToken",
					PubKey:      "invalid",
					Network:     "Network",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := keymanager.NewKeyManager(tt.args.log, tt.args.opts)
			require.Equal(t, tt.want, got)
			require.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestKeyManager_FetchValidatingPublicKeys(t *testing.T) {
	entry := logrus.NewEntry(logrus.New())

	pubKey, err := hex.DecodeString(DefaultAccountPublicKey)
	require.NoError(t, err)

	type args struct {
		log  *logrus.Entry
		opts *keymanager.Config
	}
	tests := []struct {
		name    string
		args    args
		want    [][48]byte
		wantErr bool
	}{
		{
			name: "fetch all public keys",
			args: args{
				log: entry,
				opts: &keymanager.Config{
					Location:    "Location",
					AccessToken: DefaultAccessToken,
					PubKey:      DefaultAccountPublicKey,
					Network:     "Network",
				},
			},
			want: [][48]byte{bytex.ToBytes48(pubKey)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			km, err := keymanager.NewKeyManager(tt.args.log, tt.args.opts)
			require.NoError(t, err)

			got, err := km.FetchValidatingPublicKeys(context.Background())
			require.Equal(t, tt.want, got)
			require.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestKeyManager_FetchAllValidatingPublicKeys(t *testing.T) {
	entry := logrus.NewEntry(logrus.New())

	pubKey, err := hex.DecodeString(DefaultAccountPublicKey)
	require.NoError(t, err)

	type args struct {
		log  *logrus.Entry
		opts *keymanager.Config
	}
	tests := []struct {
		name    string
		args    args
		want    [][48]byte
		wantErr bool
	}{
		{
			name: "fetch all public keys",
			args: args{
				log: entry,
				opts: &keymanager.Config{
					Location:    "Location",
					AccessToken: DefaultAccessToken,
					PubKey:      DefaultAccountPublicKey,
					Network:     "Network",
				},
			},
			want: [][48]byte{bytex.ToBytes48(pubKey)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			km, err := keymanager.NewKeyManager(tt.args.log, tt.args.opts)
			require.NoError(t, err)

			got, err := km.FetchAllValidatingPublicKeys(context.Background())
			require.Equal(t, tt.want, got)
			require.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func newTestRemoteWallet(handler http.HandlerFunc) *httptest.Server {
	s := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		handler(writer, request)
	}))

	return s
}
