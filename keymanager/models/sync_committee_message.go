package models

import types "github.com/prysmaticlabs/eth2-types"

// SignRequestSyncCommitteeMessage struct for sign req committiee msg
type SignRequestSyncCommitteeMessage struct {
	Root types.SSZBytes
}

// isSignRequestObject implement interface func
func (m *SignRequestSyncCommitteeMessage) isSignRequestObject() {}
