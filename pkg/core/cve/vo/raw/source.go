package raw

import (
	source2 "l6p.io/kun/api/pkg/core/cve/vo/raw/source"
)

type Source struct {
	Target *source2.Target `json:"target"`
}
