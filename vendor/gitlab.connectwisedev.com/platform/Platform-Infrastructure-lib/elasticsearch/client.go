package elasticsearch

import (
	"gitlab.connectwisedev.com/platform/platform-common-lib/src/communication/http/client"
	"github.com/olivere/elastic"
	"github.com/pkg/errors"
)

// NewClient creates new elasticsearch client
func NewClient(cfg Configuration) (Elastic, error) {
	// setup a http client
	httpClient := client.Basic(&client.Config{
		MaxIdleConns:          cfg.SearchMaxIdleConns,
		IdleConnTimeoutMinute: cfg.SearchIdleConnTimeoutMin,
	}, true)

	esClient, err := elastic.NewSimpleClient(elastic.SetURL(cfg.URLs...), elastic.SetHttpClient(httpClient))
	if err != nil {
		return nil, errors.Wrap(err, "can't establish connection to ElasticSearch")
	}

	return &elasticClient{
		client:                    esClient,
		indicesGetTemplateService: elastic.NewIndicesGetTemplateService(esClient),
	}, nil
}
