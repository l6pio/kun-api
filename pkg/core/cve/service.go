package cve

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"github.com/labstack/gommon/log"
	"l6p.io/kun/api/pkg/core"
	"l6p.io/kun/api/pkg/core/cve/vo/api"
	"l6p.io/kun/api/pkg/core/cve/vo/raw"
	"l6p.io/kun/api/pkg/core/db"
	"l6p.io/kun/api/pkg/core/db/query/cve"
	"os/exec"
	"strconv"
	"time"
)

var VulSeverity = map[string]int64{
	"High":       4,
	"Medium":     3,
	"Low":        2,
	"Negligible": 1,
	"Unknown":    0,
}

var VulFixState = map[string]int64{
	"fixed":     3,
	"not-fixed": 2,
	"wont-fix":  1,
	"unknown":   0,
}

func Scan(image string) *raw.Report {
	log.Infof("preparing to scan the image: %s", image)

	var buildOut bytes.Buffer
	buildCmd := exec.Command("grype", image, "-o=json")
	buildCmd.Dir = "/usr/local/bin/"
	buildCmd.Stdout = &buildOut
	buildCmd.Stderr = &buildOut

	if err := buildCmd.Run(); err != nil {
		log.Errorf("CVE scanning failed for '%s': %v", image, err)
		return nil
	}

	log.Info("scanning completed")

	var report raw.Report
	if err := json.Unmarshal(buildOut.Bytes(), &report); err != nil {
		log.Errorf("scan result parsing failed: %v", err)
		return nil
	}

	return &report
}

func Insert(conn *sql.DB, report *raw.Report) error {
	_, err := db.RunTx(conn, cve.InsertReportSql(), func(stmt *sql.Stmt) (interface{}, error) {
		for _, m := range report.Matches {
			imgId := report.Source.Target.ImageID
			imgName := report.Source.Target.UserInput
			imgSize := report.Source.Target.ImageSize

			artName := m.Artifact.Name
			artVersion := m.Artifact.Version
			artType := m.Artifact.Type
			artLang := m.Artifact.Language
			artLicenses := m.Artifact.Licenses
			artCpes := m.Artifact.Cpes
			artPurl := m.Artifact.Purl

			vulId := m.Vulnerability.ID
			vulDataSource := m.Vulnerability.DataSource
			vulNamespace := m.Vulnerability.Namespace
			vulSeverity := VulSeverity[m.Vulnerability.Severity]
			vulUrls := m.Vulnerability.Urls
			vulDesc := m.Vulnerability.Description
			vulFixVersions := m.Vulnerability.Fix.Versions
			vulFixState := VulFixState[m.Vulnerability.Fix.State]

			maxVersion := 0.0
			baseScore := 0.0
			exploitScore := 0.0
			impactScore := 0.0

			for _, cvss := range m.Vulnerability.Cvss {
				version, err := strconv.ParseFloat(cvss.Version, 64)
				if err != nil {
					continue
				}

				if version > maxVersion {
					maxVersion = version
					baseScore = cvss.Metrics.BaseScore
					exploitScore = cvss.Metrics.ExploitabilityScore
					impactScore = cvss.Metrics.ImpactScore
				}
			}

			if _, err := stmt.Exec(
				imgId, imgName, imgSize,
				artName, artVersion, artType, artLang, artLicenses, artCpes, artPurl,
				vulId, vulDataSource, vulNamespace, vulSeverity, vulUrls, vulDesc,
				maxVersion, baseScore, exploitScore, impactScore,
				vulFixVersions, vulFixState, time.Now(),
			); err != nil {
				return nil, err
			}
		}
		return nil, nil
	})
	return err
}

func List(conf *core.Config, page int, order string) (*db.Paging, error) {
	ret, err := (&db.Paging{
		Page: page,
		DoCount: func() (*sql.Rows, error) {
			return conf.DbConn.Query(cve.CountAllSql())
		},
		DoQuery: func(from int, size int) (*sql.Rows, error) {
			return conf.DbConn.Query(cve.SelectAllSql(order), from, size)
		},
		Convert: func(rows *sql.Rows) []interface{} {
			ret := make([]interface{}, 0)
			for rows.Next() {
				var id string
				var severity int64
				var fixState int64
				var imageCount int64

				if err := rows.Scan(&id, &severity, &fixState, &imageCount); err != nil {
					log.Error(err)
				}

				ret = append(ret, api.Vulnerability{
					Id:         id,
					Severity:   severity,
					FixState:   fixState,
					ImageCount: imageCount,
				})
			}
			return ret
		},
	}).Do()
	if err != nil {
		return nil, err
	}
	return ret, err
}
