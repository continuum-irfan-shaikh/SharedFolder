package elasticsearch

import (
	"context"
	"fmt"

	"github.com/olivere/elastic"
	"github.com/pkg/errors"
)

// Elastic describes common methods to work with elastic client
type Elastic interface {
	// BulkProcessor allows setting up a concurrent processor of bulk requests.
	StartBulkProcessor(config BulkConfig) (BulkProcessor, error)
	// Aggregate adds an aggreation to perform as part of the search.
	Aggregate(ctx context.Context, req *AggregationRequest) ([]string, error)
	// Search is the entry point for searches.
	Search(ctx context.Context, request *SearchRequest) (results []interface{}, err error)
	// MultiSearch is the entry point for multi searches.
	MultiSearch(ctx context.Context, request []*SearchRequest) (results [][]interface{}, err error)
	// Delete eliminates items from Elasticsearch
	Delete(ctx context.Context, id string, elasticIndex string, elasticType string) error
	// GetTemplate returns an index template.
	GetTemplate(ctx context.Context, templateName, mappingName string) (map[string]interface{}, error)
	// CountByQuery counts documents.
	CountByQuery(ctx context.Context, query ElasticQuery, indices ...string) (int64, error)
	// IsRunning returns true if the background processes of the client are
	// running, false otherwise.
	IsRunning() bool
	// Close stops the background processes that the client is running,
	// i.e. sniffing the cluster periodically and running health checks
	// on the nodes.
	//
	// If the background processes are not running, this is a no-op.
	Close() error
}

// BulkRequest is a wrapper for elastic.BulkableRequest
type BulkRequest interface {
	fmt.Stringer
	Source() ([]string, error)
}

// BulkProcessor describes methods to work with elastic bulk processor
type BulkProcessor interface {
	Add(request BulkRequest)
	Flush() error
	Close() error
}

// ElasticQuery is a wrapper for elastic.Query
type ElasticQuery interface {
	Source() (interface{}, error)
}

type bulkProcessorImpl struct {
	processor *elastic.BulkProcessor
}

func (bp *bulkProcessorImpl) Add(request BulkRequest) {
	bp.processor.Add(request)
}

func (bp *bulkProcessorImpl) Close() error {
	return bp.processor.Close()
}

func (bp *bulkProcessorImpl) Flush() error {
	return bp.processor.Flush()
}

type elasticClient struct {
	client                    *elastic.Client
	indicesGetTemplateService *elastic.IndicesGetTemplateService
}

func (e *elasticClient) Close() error {
	e.client.Stop()
	return nil
}

func (e *elasticClient) StartBulkProcessor(config BulkConfig) (BulkProcessor, error) {
	bp, err := e.client.BulkProcessor().
		Name(config.Name).
		After(config.AfterCallback).
		Workers(config.Workers).
		BulkSize(config.MaxBulkSize).
		BulkActions(config.BulkActions).
		FlushInterval(config.FlushInterval).
		Do(context.Background())

	return &bulkProcessorImpl{bp}, err
}

func (e *elasticClient) GetTemplate(ctx context.Context, templateName, mappingName string) (map[string]interface{}, error) {
	template, err := e.indicesGetTemplateService.Name(templateName).Do(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot get indices template %s from elasticsearch", templateName)
	}

	if len(template) != 1 {
		return nil, errors.Errorf("more than one template found for given name %s", templateName)
	}

	t, ok := template[templateName]
	if !ok {
		return nil, errors.Errorf("template with name %s not found in the response", templateName)
	}

	met, ok := t.Mappings[mappingName].(map[string]interface{})
	if !ok {
		return nil, errors.Errorf("template with name %s with mapping name %s doesn't match with (map[string]interface{})", templateName, mappingName)
	}

	prop, ok := met["properties"].(map[string]interface{})
	if !ok {
		return nil, errors.Errorf("mapping template with name %s cannot convert properties to (map[string]interface{})", mappingName)
	}

	return prop, nil
}

func (e *elasticClient) CountByQuery(ctx context.Context, query ElasticQuery, indices ...string) (int64, error) {
	return e.client.Count(indices...).Query(query).Do(ctx)
}

func (e *elasticClient) IsRunning() bool {
	return e.client.IsRunning()
}
