package models

import eth "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"

// SignRequestSyncAggregatorSelectionData struct for sign req committiee msg
type SignRequestSyncAggregatorSelectionData struct {
	SyncAggregatorSelectionData *eth.SyncAggregatorSelectionData
}

// isSignRequestObject implement interface func
func (m *SignRequestSyncAggregatorSelectionData) isSignRequestObject() {}
