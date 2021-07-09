package vo

import (
	"l6p.io/kun/api/pkg/core/service/vo/source"
)

type Source struct {
	Target *source.Target `json:"target"`
}
