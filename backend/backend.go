package backend

import (
	"context"
	"encoding/hex"
	"net/http"
	"sync"

	"github.com/bloxapp/key-vault/utils/encoder/encoderv2"

	encoder2 "github.com/bloxapp/key-vault/utils/encoder"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/bloxapp/key-vault/utils/errorex"
)

// Factory returns the backend factory
func Factory(version string, logger *logrus.Logger) logical.Factory {
	return func(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
		b := newBackend(version, logger)
		if err := b.Setup(ctx, conf); err != nil {
			return nil, err
		}
		return b, nil
	}
}

// newBackend returns the backend
func newBackend(version string, logger *logrus.Logger) *backend {
	b := &backend{
		logger:   logger,
		Version:  version,
		signLock: make(map[string]*sync.Mutex),
		encoder:  encoderv2.New(),
	}
	b.Backend = &framework.Backend{
		Help: "",
		Paths: framework.PathAppend(
			versionPaths(b),
			storagePaths(b),
			storageSlashingDataPaths(b),
			accountsPaths(b),
			signsPaths(b),
			configPaths(b),
		),
		PathsSpecial: &logical.Paths{
			SealWrapStorage: []string{
				"wallet/",
			},
		},
		Secrets:     []*framework.Secret{},
		BackendType: logical.TypeLogical,
	}
	return b
}

// backend implements the Backend for this plugin
type backend struct {
	*framework.Backend
	logger   *logrus.Logger
	Version  string
	signLock map[string]*sync.Mutex
	encoder  encoder2.IEncoder
}

// pathExistenceCheck checks if the given path exists
func (b *backend) pathExistenceCheck(ctx context.Context, req *logical.Request, data *framework.FieldData) (bool, error) {
	out, err := req.Storage.Get(ctx, req.Path)
	if err != nil {
		b.logger.WithError(err).Error("Path existence check failed")
		return false, errors.Errorf("existence check failed: %v", err)
	}

	if b.Route(req.Path) == nil {
		b.logger.WithField("path", hex.EncodeToString(out.Value)).Error("Path not found")
		return false, errors.Errorf("existence check failed: %s", hex.EncodeToString(out.Value))
	}

	return out != nil, nil
}

// notFoundResponse returns not found error
func (b *backend) notFoundResponse() (*logical.Response, error) {
	return logical.RespondWithStatusCode(&logical.Response{
		Data: map[string]interface{}{
			"message":     "account not found",
			"status_code": http.StatusNotFound,
		},
	}, nil, http.StatusNotFound)
}

// prepareErrorResponse prepares error response
func (b *backend) prepareErrorResponse(originError error) (*logical.Response, error) {
	switch err := errors.Cause(originError).(type) {
	case *errorex.ErrBadRequest:
		return err.ToLogicalResponse()
	case nil:
		return nil, nil
	default:
		return logical.ErrorResponse(originError.Error()), nil
	}
}
