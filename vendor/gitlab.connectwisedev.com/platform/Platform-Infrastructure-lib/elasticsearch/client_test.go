package elasticsearch

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestElasticNegative(t *testing.T) {
	ctx := context.Background()

	el, err := NewClient(Configuration{})
	defer el.Close()

	assert.NotNil(t, el)
	assert.Nil(t, err)

	isRunning := el.IsRunning()
	assert.True(t, isRunning)

	request := SearchRequest{}

	_, err = el.Search(ctx, &request)
	assert.NotNil(t, err)

	_, err = el.MultiSearch(ctx, []*SearchRequest{&request})
	assert.NotNil(t, err)

	_, err = el.GetTemplate(ctx, "template_managed_endpoint", "managed_endpoints")
	assert.NotNil(t, err)

	agrRequest := AggregationRequest{
		Partner: "parnerID",
		Client:  "clientID",
		Site:    "siteID",
	}
	_, err = el.Aggregate(ctx, &agrRequest)
	assert.NotNil(t, err)

	var eq ElasticQuery
	_, err = el.CountByQuery(ctx, eq)

	err = el.Delete(ctx, "testID", "testID", "testType")
	assert.NotNil(t, err)
}

func TestBulkProcessorNegative(t *testing.T) {
	el, err := NewClient(Configuration{})
	defer el.Close()

	assert.NotNil(t, el)
	assert.Nil(t, err)

	bp, err := el.StartBulkProcessor(BulkConfig{})
	assert.NotNil(t, bp)
	assert.Nil(t, err)

	err = bp.Flush()
	assert.Nil(t, err)

	err = bp.Close()
	assert.Nil(t, err)
}
