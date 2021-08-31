package vo

type SettingType int32

const (
	RegistryAuthType = 1
)

type RegistryAuth struct {
	Type      SettingType `json:"type" bson:"type"`
	Authority string      `json:"authority" bson:"authority"`
	Username  string      `json:"username" bson:"username"`
	Password  string      `json:"password" bson:"password"`
}
