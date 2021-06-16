package cve

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/labstack/gommon/log"
	"github.com/olivere/elastic/v7"
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

			log.Infof("Preparing to scan the image: %s", img)

			var buildOut bytes.Buffer
			buildCmd := exec.Command("grype", img, "-o=json")
			buildCmd.Dir = "/usr/local/bin/"
			buildCmd.Stdout = &buildOut
			buildCmd.Stderr = &buildOut

			if err := buildCmd.Run(); err != nil {
				log.Errorf("CVE scanning failed for '%s': %v\n", img, err)
			}

			log.Info("Scanning completed")

			var report vo.Report
			if err := json.Unmarshal(buildOut.Bytes(), &report); err != nil {
				log.Errorf("Scan result parsing failed: %v\n", err)
			}

			log.Info("Start indexing scan results")

			id, err := es.Index(conf, report)
			if err != nil {
				log.Errorf("Index scan result of '%s' failed: %v\n", img, err)
			}
			log.Infof("Image with ID '%v' has been indexed", id)
		}
	}()
}

func ListAll(conf *core.Config) ([]*vo.Report, error) {
	return es.Search(conf, elastic.MatchAllQuery{})
}

func FindByImageID(conf *core.Config, imageID string) ([]*vo.Report, error) {
	query := elastic.NewNestedQuery(
		"source", elastic.NewNestedQuery(
			"source.target", elastic.NewTermQuery("source.target.imageID", imageID),
		),
	)
	return es.Search(conf, query)
}
