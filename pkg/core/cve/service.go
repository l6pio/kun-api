package cve

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/labstack/gommon/log"
	"l6p.io/kun/api/pkg/core"
	"l6p.io/kun/api/pkg/core/cve/vo"
	"l6p.io/kun/api/pkg/core/es"
	"os/exec"
)

func Scan(conf *core.Config) {
	go func() {
		for {
			key := <-conf.ScanRequests
			img := fmt.Sprintf("%s:%s", key.ImageRepo, key.ImageTag)

			var buildOut bytes.Buffer
			buildCmd := exec.Command("grype", img, "-o=json")
			buildCmd.Dir = "/usr/local/bin/"
			buildCmd.Stdout = &buildOut
			buildCmd.Stderr = &buildOut

			if err := buildCmd.Run(); err != nil {
				log.Errorf("CVE scanning failed for '%s': %v\n", img, err)
			}

			var report vo.Report
			if err := json.Unmarshal(buildOut.Bytes(), &report); err != nil {
				log.Errorf("Scan result parsing failed: %v\n", err)
			}

			err := es.Index(conf, report)
			if err != nil {
				log.Errorf("Index scan result of '%s' failed: %v\n", img, err)
			}
		}
	}()
}
