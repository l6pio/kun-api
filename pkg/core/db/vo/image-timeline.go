package vo

type ImageStatus int

const (
	ImageUp   ImageStatus = 1
	ImageDown ImageStatus = 0
)

type ImageTimeline struct {
	Timestamp int64       `json:"timestamp" bson:"timestamp"`
	ImageId   string      `json:"imageId" bson:"imageId"`
	Image     string      `json:"image" bson:"image"`
	Status    ImageStatus `json:"status" bson:"status"`
}
