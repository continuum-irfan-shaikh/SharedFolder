package main

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"gitlab.connectwisedev.com/platform/Platform-Infrastructure-lib/elasticsearch"
	"gitlab.connectwisedev.com/platform/platform-common-lib/src/runtime/logger"
	"gitlab.connectwisedev.com/platform/platform-common-lib/src/uuid"
	"github.com/olivere/elastic"
)

func main() {
	log, err := logger.Create(logger.Config{Destination: logger.STDOUT})
	if err != nil {
		fmt.Println(err)
		return
	}

	client, err := elasticsearch.NewClient(elasticsearch.Configuration{
		URLs:                     []string{"http://localhost:9200"},
		SearchMaxIdleConns:       100,
		SearchIdleConnTimeoutMin: 1,
	})
	if err != nil {
		log.Error("", "", "Cannot open elastic search connection: %s", err)
		return
	}
	defer func() {
		err = client.Close()
		if err != nil {
			log.Error("", "", "Cannot close elastic search connection : %s", err)
		}
	}()

	tpl, err := client.GetTemplate(context.Background(), "me_template", "managed_endpoints")
	if err != nil {
		log.Warn("", "Cannot get template: %s", err)
		return
	}

	fmt.Printf("%+v\n", tpl)
	fmt.Println("-----------------------------------------------")

	bulkProc, err := client.StartBulkProcessor(elasticsearch.BulkConfig{
		Name:          "Managed endpoints",
		Workers:       1,
		FlushInterval: time.Millisecond * 250,
		MaxBulkSize:   300,
		AfterCallback: AfterCallback(log),
	})
	if err != nil {
		log.Error("", "", "Cannot start bulk processor: %s", err)
		return
	}
	defer func() {
		if err = bulkProc.Close(); err != nil {
			log.Error("", "", "Cannot stop bulk processor: %s", err)
		}
	}()

	id, err := uuid.NewRandomUUID()
	if err != nil {
		log.Warn("", "Cannot generate uuid: %s", err)
	}

	req := elastic.NewBulkIndexRequest().
		Index("managed_endpoints_1337").
		Type("managed_endpoints").
		Id(id.String()).
		Doc(&ManagedEndpoint{
			Name:    "Some name",
			ID:      id.String(),
			Client:  "client",
			Partner: "partner",
			Site:    "site",
		})

	bulkProc.Add(req)

	if err = bulkProc.Flush(); err != nil {
		log.Warn("", "Cannot Flush bulk processor: %s", err)
		return
	}

	time.Sleep(time.Second)

	multiSearch(client, log)

	if err = client.Delete(context.Background(), id.String(), "managed_endpoints_1337", "managed_endpoints"); err != nil {
		log.Warn("", "Cannot delete by id: %s", err)
		return
	}

	time.Sleep(time.Second)

	search(client, log)

	s, err := client.Aggregate(context.Background(), &elasticsearch.AggregationRequest{
		LookupField:      "client",
		Partner:          "partner",
		Index:            "managed_endpoints_1337",
		MaxNumberOfItems: 3000,
	})
	if err != nil {
		log.Warn("", "Cannot aggregate: %s", err)
		return
	}

	fmt.Println(s)
}

type ManagedEndpoint struct {
	Name    string `json:"name"`
	ID      string `json:"id"`
	Client  string `json:"client"`
	Partner string `json:"partner"`
	Site    string `json:"site"`
}

// AfterCallback is called after each request is processed by elasticsearch
func AfterCallback(log logger.Log) elastic.BulkAfterFunc {
	return func(_ int64, request []elastic.BulkableRequest, resp *elastic.BulkResponse, err error) {
		if err != nil {
			log.Error("", "", "Error persisting message %v to elasticsearch: %v\n", request, err)
		}

		if resp != nil && resp.Errors {
			for _, item := range resp.Items {
				for key, value := range item {
					if value.Error != nil {
						log.Error("", "",
							"Error persisting message to elasticsearch. %s failed for %s/%s/%s. Error status: %d, type: %s, reason: %s",
							key, value.Index, value.Type, value.Id,
							value.Status, value.Error.Type, value.Error.Reason,
						)
					}
				}
			}
		}
	}
}

func search(client elasticsearch.Elastic, log logger.Log) {
	q := elastic.NewBoolQuery().Filter(
		elastic.NewQueryStringQuery(fmt.Sprintf("partner:\"%s\"", "partner")),
		elastic.NewQueryStringQuery(fmt.Sprintf("client:\"%s\"", "client")),
		elastic.NewQueryStringQuery(fmt.Sprintf("site:\"%s\"", "site")),
	)
	results, err := client.Search(context.Background(), &elasticsearch.SearchRequest{
		ElasticIndex: "managed_endpoints_1337",
		ElasticType:  "managed_endpoints",
		ElasticQuery: q,
		Size:         3000,
		ResultType:   reflect.TypeOf(ManagedEndpoint{}),
	})
	if err != nil {
		if elasticsearch.IsErrNotFound(err) {
			log.Warn("", err.Error())
			return
		}
		log.Error("", "", "Unable to search due to %s", err)
		return
	}

	for i, res := range results {
		me, ok := res.(ManagedEndpoint)
		if !ok {
			log.Warn("", "casting failed %v", res)
			continue
		}

		fmt.Printf("%d: %v\n", i, me)
	}
}

func multiSearch(client elasticsearch.Elastic, log logger.Log) {
	results, err := client.MultiSearch(context.Background(), []*elasticsearch.SearchRequest{{
		ElasticIndex: "managed_endpoints_1337",
		ElasticType:  "managed_endpoints",
		Size:         3000,
		ResultType:   reflect.TypeOf(ManagedEndpoint{}),
	}})
	if err != nil {
		log.Error("", "", "MultiSearch is failed %s", err)
		return
	}

	for i, res := range results {
		for j, r := range res {
			me, ok := r.(ManagedEndpoint)
			if !ok {
				log.Warn("", "casting is failed %v", r)
				continue
			}

			fmt.Printf("%d.%d: %v\n", i, j, me)
		}
	}
}
