package es

import (
	"context"
	_ "embed"
	"encoding/json"
	"github.com/labstack/gommon/log"
	"github.com/olivere/elastic/v7"
	"l6p.io/kun/api/pkg/core"
	"l6p.io/kun/api/pkg/core/cve/vo"
)

const IndexName = "cve"

//go:embed mappings.json
var IndexMappings string

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
		return "", nil
	}

	if err != nil {
		return "", err
	}
	return res.Id, nil
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
