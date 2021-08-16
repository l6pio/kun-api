package vo

import v1 "k8s.io/api/core/v1"

type Pod struct {
	Name      string      `json:"name" bson:"name"`
	Namespace string      `json:"namespace" bson:"namespace"`
	Phase     v1.PodPhase `json:"phase" bson:"phase"`
	Status    string      `json:"status" bson:"status"`
	Started   int64       `json:"started" bson:"started"`
	Finished  int64       `json:"finished" bson:"finished"`
}
