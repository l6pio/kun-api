package vo

import v1 "k8s.io/api/core/v1"

type Pod struct {
	Name         string      `json:"name" bson:"name"`
	Namespace    string      `json:"namespace" bson:"namespace"`
	Phase        v1.PodPhase `json:"phase" bson:"phase"`
	Status       string      `json:"status" bson:"status"`
	Ready        int64       `json:"ready" bson:"ready"`
	Finished     int64       `json:"finished" bson:"finished"`
	RestartCount int32       `json:"restartCount" bson:"restartCount"`
}
