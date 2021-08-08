package vo

type ImageStatus int

const (
	ImageUp   ImageStatus = 1
	ImageDown ImageStatus = 0
)

type Image struct {
	Id   string `json:"id" bson:"id"`
	Name string `json:"name" bson:"name"`
	Size int64  `json:"size" bson:"size"`
	Pods int64  `json:"pods" bson:"pods"`
}

type ImageTimeline struct {
	Timestamp int64       `json:"timestamp" bson:"timestamp"`
	ImageId   string      `json:"imageId" bson:"imageId"`
	Image     string      `json:"image" bson:"image"`
	Status    ImageStatus `json:"status" bson:"status"`
}
