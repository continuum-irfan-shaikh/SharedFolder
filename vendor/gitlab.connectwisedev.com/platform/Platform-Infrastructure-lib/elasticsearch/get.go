package elasticsearch

import (
	"context"
	"fmt"

	"github.com/olivere/elastic"
)

const deleted = "deleted"

// Search looks up query in Elasticsearch
func (e *elasticClient) Search(ctx context.Context, request *SearchRequest) ([]interface{}, error) {
	srcCtx := elastic.NewFetchSourceContext(true).Include(request.IncludedFields...)
	searchReq := e.client.
		Search(request.ElasticIndex).
		Type(request.ElasticType).
		Query(request.ElasticQuery).
		From(request.StartIndex).
		Size(request.Size).
		FetchSourceContext(srcCtx)

	result, err := searchReq.Do(ctx)
	if err != nil {
		return nil, err
	}

	results := result.Each(request.ResultType)

	if len(results) < 1 {
		return nil, NewErrNotFound("no results found type %s, index %s", request.ElasticType, request.ElasticIndex)
	}

	return results, err
}

// MultiSearch looks up queries in Elasticsearch
func (e *elasticClient) MultiSearch(ctx context.Context, requests []*SearchRequest) ([][]interface{}, error) {
	resp, err := e.client.MultiSearch().Add(convertToElasticRequests(requests...)...).Do(ctx)
	if err != nil {
		return nil, err
	}

	if len(requests) != len(resp.Responses) {
		return nil, fmt.Errorf("was made %d requests but got %d responses", len(requests), len(resp.Responses))
	}

	results := make([][]interface{}, len(resp.Responses))
	for i, res := range resp.Responses {
		nestedResult := make([]interface{}, res.TotalHits())

		for j, item := range res.Each(requests[i].ResultType) {
			nestedResult[j] = item
		}
		results[i] = nestedResult
	}

	return results, err
}

// Delete eliminates items from Elasticsearch
func (e *elasticClient) Delete(ctx context.Context, id string, elasticIndex string, elasticType string) error {
	resp, err := e.client.Delete().Id(id).Index(elasticIndex).Type(elasticType).Do(ctx)
	if err != nil {
		return err
	}
	if resp != nil && resp.Result != deleted {
		return fmt.Errorf("resp.Result - %q", resp.Result)
	}
	return err
}

func convertToElasticRequests(requests ...*SearchRequest) []*elastic.SearchRequest {
	esRequests := make([]*elastic.SearchRequest, len(requests))
	for i, req := range requests {
		r := elastic.NewSearchRequest().
			Index(req.ElasticIndex).
			Type(req.ElasticType).
			SearchSource(
				elastic.NewSearchSource().
					Size(req.Size).
					From(req.StartIndex).
					Query(req.ElasticQuery).
					FetchSource(true).
					FetchSourceContext(elastic.NewFetchSourceContext(true).Include(req.IncludedFields...)))

		esRequests[i] = r
	}
	return esRequests
}
