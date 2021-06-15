package es

import (
	"context"
	"encoding/json"
	"github.com/labstack/gommon/log"
	"github.com/olivere/elastic/v7"
	"l6p.io/kun/api/pkg/core"
	"l6p.io/kun/api/pkg/core/cve/vo"
)

const IndexName = "cve"

const IndexMappings = `
{
  "mappings": {
	"properties": {
      "matches": {
        "type": "nested",
        "properties": {
		  "artifact": {
			"type": "nested",
			"properties": {
			  "name": {"type": "keyword"},
			  "version": {"type": "keyword"},
			  "type": {"type": "keyword"},
			  "locations": {
				"type": "nested",
				"properties": {
				  "path": {"type": "text"},
				  "layerID": {"type": "keyword"}
				}
              },
			  "language": {"type": "keyword"},
			  "licenses": {"type": "keyword"},
			  "cpes": {"type": "keyword"},
			  "purl": {"type": "text"},
			  "metadata": {
                "type": "nested",
                "properties": {
                  "VirtualPath": {"type": "text"},
                  "PomArtifactID": {"type": "keyword"},
                  "PomGroupID": {"type": "keyword"}
                }
			  }
			}
		  },
		  "vulnerability": {
			"type": "nested",
			"properties": {
              "id": {"type": "keyword"},
              "dataSource": {"type": "text"},
              "namespace": {"type": "keyword"},
              "severity": {"type": "keyword"},
              "urls": {"type": "text"},
              "description": {"type": "text"},
              "cvss": {
			    "type": "nested",
                "properties": {
                  "version": {"type": "keyword"},
                  "metrics": {
                    "type": "nested",
                    "properties": {
                      "baseScore": {"type": "double"},
                      "exploitabilityScore": {"type": "double"},
                      "impactScore": {"type": "double"}
                    }
                  }
                }
              },
			  "fix": {
			    "type": "nested",
                "properties": {
				  "versions": {"type": "keyword"},
				  "state": {"type": "keyword"}
				}                
              }
			}
		  }
		}
      },
	  "source": {
		"type": "nested",
		"properties": {
		  "target": {
            "type": "nested",
            "properties": {
			  "userInput": {"type": "keyword"},
			  "imageID": {"type": "keyword"},
			  "imageSize": {"type": "long"}
			}
		  }
		}
	  }
	}
  }
}
`

func CreateIndex(conf *core.Config) {
	ctx := context.Background()

	exists, err := conf.EsClient.IndexExists(IndexName).Do(ctx)
	if err != nil {
		log.Fatalf("Failed to check Elasticsearch index: %v", err)
	}

	if exists {
		return
	}

	_, err = conf.EsClient.CreateIndex(IndexName).BodyString(IndexMappings).Do(ctx)
	if err != nil {
		log.Fatalf("Failed to create ES index: %v", err)
	}
	log.Infof("Elasticsearch index '%v' is created", IndexName)
	return
}

func Index(conf *core.Config, report vo.Report) (string, error) {
	id := report.Source.Target.ImageID
	res, err := conf.EsClient.Index().Index(IndexName).
		Id(id).OpType("create").BodyJson(report).
		Do(context.Background())

	if elastic.IsConflict(err) {
		return id, nil
	}

	if err != nil {
		return "", err
	}

	log.Infof("Docker image with ID '%v' has been indexed", res.Id)
	return id, nil
}

func SearchByImageID(conf *core.Config, imageID string) ([]*vo.Report, error) {
	query := elastic.NewNestedQuery(
		"source", elastic.NewNestedQuery(
			"source.target", elastic.NewTermQuery("source.target.imageID", imageID),
		),
	)
	return Search(conf, query)
}

func Search(conf *core.Config, query elastic.Query) ([]*vo.Report, error) {
	result, err := conf.EsClient.Search(IndexName).
		Query(query).
		Do(context.Background())
	if err != nil {
		return nil, err
	}

	var ret = make([]*vo.Report, 0)
	for _, hits := range result.Hits.Hits {
		report := vo.Report{}
		if err := json.Unmarshal(hits.Source, &report); err != nil {
			return nil, err
		}
		ret = append(ret, &report)
	}
	return ret, nil
}
