package service

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/labstack/gommon/log"
	"l6p.io/kun/api/pkg/core"
	"l6p.io/kun/api/pkg/core/db"
	dbvo "l6p.io/kun/api/pkg/core/db/vo"
	"l6p.io/kun/api/pkg/core/service/vo"
	"os/exec"
	"strconv"
	"time"
)

var VulSeverity = map[string]int64{
	"Critical":   5,
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

func Scan(image string) (*vo.Report, error) {
	log.Infof("preparing to scan the image: %s", image)

	var buildOut bytes.Buffer
	buildCmd := exec.Command("grype", image, "-o=json")
	buildCmd.Dir = "/usr/local/bin/"
	buildCmd.Stdout = &buildOut
	buildCmd.Stderr = &buildOut

	if err := buildCmd.Run(); err != nil {
		log.Errorf("CVE scanning failed for '%s': %v", image, err)
		return nil, err
	}
	log.Info("scanning completed")

	var report vo.Report
	if err := json.Unmarshal(buildOut.Bytes(), &report); err != nil {
		log.Errorf("scan result parsing failed: %v", err)
		return nil, err
	}
	return &report, nil
}

func Insert(conf *core.Config, imageId string, report *vo.Report) {
	img := dbvo.Image{
		Id:   imageId,
		Name: report.Source.Target.UserInput,
		Size: report.Source.Target.ImageSize,
	}

	if err := db.SaveImage(conf, &img); err != nil {
		log.Error(err)
		return
	}

	for _, m := range report.Matches {
		// The namespace for the vulnerability is required to come from "nvd".
		vulnerability := m.Vulnerability
		if vulnerability.Namespace != "nvd" {
			for _, rv := range m.RelatedVulnerabilities {
				if rv.ID == m.Vulnerability.ID && rv.Namespace == "nvd" {
					vulnerability = rv
				}
			}
		}
		if vulnerability.Namespace != "nvd" {
			continue
		}

		cvssVersion := 0.0
		cvssBaseScore := 0.0
		cvssExploitScore := 0.0
		cvssImpactScore := 0.0

		for _, cvss := range vulnerability.Cvss {
			version, err := strconv.ParseFloat(cvss.Version, 64)
			if err != nil {
				continue
			}

			if version > cvssVersion {
				cvssVersion = version
				cvssBaseScore = cvss.Metrics.BaseScore
				cvssExploitScore = cvss.Metrics.ExploitabilityScore
				cvssImpactScore = cvss.Metrics.ImpactScore
			}
		}

		art := dbvo.Artifact{
			Id:       convertToId(fmt.Sprintf("%s:%s", m.Artifact.Name, m.Artifact.Version)),
			Name:     m.Artifact.Name,
			Version:  m.Artifact.Version,
			Type:     m.Artifact.Type,
			Language: m.Artifact.Language,
			Licenses: m.Artifact.Licenses,
			Cpes:     m.Artifact.Cpes,
			Purl:     m.Artifact.Purl,
		}

		vul := dbvo.Vulnerability{
			Id:               vulnerability.ID,
			DataSource:       vulnerability.DataSource,
			Namespace:        vulnerability.Namespace,
			Severity:         VulSeverity[vulnerability.Severity],
			Urls:             vulnerability.Urls,
			Description:      vulnerability.Description,
			FixVersions:      m.Vulnerability.Fix.Versions,
			FixState:         VulFixState[m.Vulnerability.Fix.State],
			CvssVersion:      cvssVersion,
			CvssBaseScore:    cvssBaseScore,
			CvssExploitScore: cvssExploitScore,
			CvssImpactScore:  cvssImpactScore,
		}

		if err := db.SaveArtifact(conf, &art); err != nil {
			log.Error(err)
			continue
		}

		if err := db.SaveVulnerability(conf, &vul); err != nil {
			log.Error(err)
			continue
		}

		if err := db.SaveCve(conf, &dbvo.Cve{
			ImgId: img.Id,
			ArtId: art.Id,
			VulId: vul.Id,
		}); err != nil {
			log.Error(err)
		}
	}
}

func UpdateVulnerabilityDatabase() error {
	var buildOut bytes.Buffer
	buildCmd := exec.Command("grype", "db", "update")
	buildCmd.Dir = "/usr/local/bin/"
	buildCmd.Stdout = &buildOut
	buildCmd.Stderr = &buildOut

	if err := buildCmd.Run(); err != nil {
		log.Errorf("update vulnerability database failed: %v", err)
		return err
	}
	log.Info(string(buildOut.Bytes()))
	log.Info("update vulnerability database completed")
	return nil
}

func PeriodicallyUpdateVulnerabilityDatabase() {
	t := time.NewTicker(time.Hour)
	defer t.Stop()
	for {
		if err := UpdateVulnerabilityDatabase(); err != nil {
			log.Error(err)
		}
		<-t.C
	}
}

func convertToId(src string) string {
	sha := sha256.New()
	sha.Write([]byte(src))
	return hex.EncodeToString(sha.Sum(nil))
}
