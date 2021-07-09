package vo

type Image struct {
	Id    string `json:"id" bson:"id"`
	Name  string `json:"name" bson:"name"`
	Size  int64  `json:"size" bson:"size"`
	Usage int64  `json:"usage" bson:"usage"`
}
