package models

import eth "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"

type SignRequest_SyncAggregatorSelectionData struct {
	SyncAggregatorSelectionData *eth.SyncAggregatorSelectionData
}

func (m *SignRequest_SyncAggregatorSelectionData) isSignRequest_Object() {}
