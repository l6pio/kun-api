package cve

import (
	source2 "l6p.io/kun/api/pkg/core/service/vo/cve/source"
)

type Source struct {
	Target *source2.Target `json:"target"`
}
