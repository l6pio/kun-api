package vo

import v1 "k8s.io/api/core/v1"

type Pod struct {
	Timestamp int64       `json:"timestamp" bson:"timestamp"`
	Namespace string      `json:"namespace" bson:"namespace"`
	Name      string      `json:"name" bson:"name"`
	Phase     v1.PodPhase `json:"phase" bson:"phase"`
	Status    string      `json:"status" bson:"status"`
}

type PodTimeline struct {
	Timestamp int64       `json:"timestamp" bson:"timestamp"`
	Name      string      `json:"name" bson:"name"`
	Status    ImageStatus `json:"status" bson:"status"`
}
