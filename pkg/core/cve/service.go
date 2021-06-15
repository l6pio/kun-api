package cve

import (
	"bytes"
	"encoding/json"
	"fmt"
	"l6p.io/kun/api/pkg/core"
	"l6p.io/kun/api/pkg/core/cve/vo"
	"l6p.io/kun/api/pkg/core/es"
	"os/exec"
)

func Scan(conf *core.Config, imageRepo string, imageTag string) (string, error) {
	var buildOut bytes.Buffer
	buildCmd := exec.Command("grype", fmt.Sprintf("%s:%s", imageRepo, imageTag), "-o=json")
	buildCmd.Dir = "/usr/local/bin/"
	buildCmd.Stdout = &buildOut
	buildCmd.Stderr = &buildOut

	if err := buildCmd.Run(); err != nil {
		return "", err
	}

	var report vo.Report
	if err := json.Unmarshal(buildOut.Bytes(), &report); err != nil {
		return "", err
	}

	id, err := es.Index(conf, report)
	if err != nil {
		return "", err
	}
	return id, nil
}
