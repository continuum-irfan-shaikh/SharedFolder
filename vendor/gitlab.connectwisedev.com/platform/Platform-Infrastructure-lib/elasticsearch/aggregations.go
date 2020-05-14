package elasticsearch

import (
	"context"
	"fmt"

	"github.com/olivere/elastic"
)

const aggregate = "aggregate"

// Aggregate performs an aggregation with specified parameters on lookupField for provided partner/client/site level
func (e *elasticClient) Aggregate(ctx context.Context, req *AggregationRequest) ([]string, error) {
	termsAggregator := elastic.NewTermsAggregation().
		Field(req.LookupField).
		Size(req.MaxNumberOfItems).
		OrderByTermAsc()

	s, err := e.client.Search().
		Index(req.Index).
		Query(getQueryFromRequest(req)).
		Aggregation(aggregate, termsAggregator).
		Do(ctx)
	if err != nil {
		return nil, err
	}

	var list []string
	if agg, found := s.Aggregations.Terms(aggregate); found {
		for _, bucket := range agg.Buckets {
			b, ok := bucket.Key.(string)
			if !ok {
				continue
			}
			list = append(list, b)
		}
	}
	return list, nil
}

// AggregationRequest describes aggregation request parameters
type AggregationRequest struct {
	LookupField      string // field that should be aggregated
	Index            string // elastic index for search
	Partner          string // can be empty then search will be done for all partners
	Client           string // can be empty then search will be done for all clients (partner must be specified)
	Site             string // can be empty then search will be done for all sites (partner and client must be specified)
	MaxNumberOfItems int    // max number for retrieved items
}

func getQueryFromRequest(req *AggregationRequest) ElasticQuery {
	var query []elastic.Query
	if req.Partner != "" {
		query = append(query, elastic.NewQueryStringQuery(fmt.Sprintf("partner:\"%s\"", req.Partner)))
	}
	if req.Client != "" {
		query = append(query, elastic.NewQueryStringQuery(fmt.Sprintf("client:\"%s\"", req.Client)))
	}
	if req.Site != "" {
		query = append(query, elastic.NewQueryStringQuery(fmt.Sprintf("site:\"%s\"", req.Site)))
	}

	return elastic.NewBoolQuery().Filter(query...)
}
