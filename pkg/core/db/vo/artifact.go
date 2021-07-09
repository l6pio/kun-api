package vo

type Artifact struct {
	Name     string   `json:"name" bson:"name"`
	Version  string   `json:"version" bson:"version"`
	Type     string   `json:"type" bson:"type"`
	Language string   `json:"language" bson:"language"`
	Licenses []string `json:"licenses" bson:"licenses"`
	Cpes     []string `json:"cpes" bson:"cpes"`
	Purl     string   `json:"purl" bson:"purl"`
}
