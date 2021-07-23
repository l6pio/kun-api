package vo

type Image struct {
	Id   string `json:"id" bson:"id"`
	Name string `json:"name" bson:"name"`
	Size int64  `json:"size" bson:"size"`
	Pods int64  `json:"pods" bson:"pods"`
}
