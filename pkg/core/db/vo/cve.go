package vo

type Cve struct {
	ImgId string `json:"imgId" bson:"imgId"`
	ArtId string `json:"artId" bson:"artId"`
	VulId string `json:"vulId" bson:"vulId"`
}
