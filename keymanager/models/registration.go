package models

import eth "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"

// SignRequestRegistration struct
type SignRequestRegistration struct {
	Registration *eth.ValidatorRegistrationV1
}

// isSignRequestObject implement func
func (m *SignRequestRegistration) isSignRequestObject() {}
