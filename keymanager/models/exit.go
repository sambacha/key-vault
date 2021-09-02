package models

import eth "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"

type SignRequest_Exit struct {
	Exit *eth.VoluntaryExit
}

func (m *SignRequest_Exit) isSignRequest_Object() {}
