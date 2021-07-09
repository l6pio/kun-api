package vo

type Cve struct {
	ImgId      string `json:"imgId" bson:"imgId"`
	ArtName    string `json:"artName" bson:"artName"`
	ArtVersion string `json:"artVersion" bson:"artVersion"`
	VulId      string `json:"vulId" bson:"vulId"`
}
