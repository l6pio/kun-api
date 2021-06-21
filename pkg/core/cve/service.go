package cve

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/labstack/gommon/log"
	"l6p.io/kun/api/pkg/core"
	"l6p.io/kun/api/pkg/core/cve/vo/api"
	"l6p.io/kun/api/pkg/core/cve/vo/raw"
	"l6p.io/kun/api/pkg/core/db"
	"l6p.io/kun/api/pkg/core/db/query/cve"
	"os/exec"
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
				continue
			}

			log.Info("Scanning completed")

			var report raw.Report
			if err := json.Unmarshal(buildOut.Bytes(), &report); err != nil {
				log.Errorf("Scan result parsing failed: %v\n", err)
				continue
			}

			if len(report.Matches) == 0 {
				log.Info("No vulnerabilities found")
			}

			log.Info("Start saving scan results")

			err := Insert(conf.DbConn, &report)
			if err != nil {
				log.Error(err)
			}

			log.Infof("Scan results of '%v' has been saved", report.Source.Target.UserInput)
		}
	}()
}

func Insert(conn *sql.DB, report *raw.Report) error {
	return db.RunTx(conn, cve.InsertReportIntoCveTableSQL(), func(stmt *sql.Stmt) error {
		for _, m := range report.Matches {
			imgName := report.Source.Target.UserInput
			imgId := report.Source.Target.ImageID
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

			if len(m.Vulnerability.Cvss) == 0 {
				if _, err := stmt.Exec(
					imgName, imgId, imgSize,
					artName, artVersion, artType, artLang, artLicenses, artCpes, artPurl,
					vulId, vulDataSource, vulNamespace, vulSeverity, vulUrls, vulDesc,
					"", 0.0, 0.0, 0.0,
					vulFixVersions, vulFixState, time.Now(),
				); err != nil {
					return err
				}
			} else {
				for _, cvss := range m.Vulnerability.Cvss {
					vulCvssVersion := cvss.Version
					vulCvssBaseScore := cvss.Metrics.BaseScore
					vulCvssExploitScore := cvss.Metrics.ExploitabilityScore
					vulCvssImpactScore := cvss.Metrics.ImpactScore

					if _, err := stmt.Exec(
						imgName, imgId, imgSize,
						artName, artVersion, artType, artLang, artLicenses, artCpes, artPurl,
						vulId, vulDataSource, vulNamespace, vulSeverity, vulUrls, vulDesc,
						vulCvssVersion, vulCvssBaseScore, vulCvssExploitScore, vulCvssImpactScore,
						vulFixVersions, vulFixState, time.Now(),
					); err != nil {
						return err
					}
				}
			}
		}
		return nil
	})
}

func List(conf *core.Config, page int, order string) (*db.Paging, error) {
	ret, err := (&db.Paging{
		Page: page,
		DoCount: func() (*sql.Rows, error) {
			return conf.DbConn.Query(cve.CountAllCvesSQL())
		},
		DoQuery: func(from int, size int) (*sql.Rows, error) {
			return conf.DbConn.Query(cve.SelectAllCvesSQL(order), from, size)
		},
		Convert: func(rows *sql.Rows) []interface{} {
			ret := make([]interface{}, 0)
			for rows.Next() {
				var vulId string
				var vulSeverity int64
				if err := rows.Scan(&vulId, &vulSeverity); err != nil {
					log.Error(err)
				}
				ret = append(ret, api.Vulnerability{
					Id:       vulId,
					Severity: vulSeverity,
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
