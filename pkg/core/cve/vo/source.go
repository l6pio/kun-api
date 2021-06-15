package vo

import "l6p.io/kun/api/pkg/core/cve/vo/source"

type Source struct {
	Target *source.Target `json:"target"`
}
