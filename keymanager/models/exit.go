package models

import eth "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"

// SignRequestExit struct
type SignRequestExit struct {
	Exit *eth.VoluntaryExit
}

// isSignRequestObject implement func
func (m *SignRequestExit) isSignRequestObject() {}
