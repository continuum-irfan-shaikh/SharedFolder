
# ElasticSearch

ElasticSearch is a wrapper on top of [olivere/elastic](https://github.com/olivere/elastic) → An ElasticSearch client library

- **License** [Apache License 2.0](https://github.com/elastic/elasticsearch/blob/master/LICENSE.txt)

## To run an Example
 - Step 1  
`
curl  -X PUT 0.0.0.0:9200/_template/me_template \
-H 'Content-Type: application/json' \
-d '{  
      "template":"managed_endpoint*",
      "settings":{ 
         "index":{ 
            "analysis":{ 
               "normalizer":{ 
                  "lowerasciinormalizer":{ 
                     "filter":[ 
                        "lowercase",
                        "asciifolding"
                     ],
                     "type":"custom"
                  }
               }
            },
            "max_result_window":"30000"
         }
      },
      "mappings":{
	      		"managed_endpoints":{ 
	            "properties":{ 
	            "name":{
	            	"type":"text"
	            	},
	               "client":{ 
	                  "type":"keyword"
	               },
	               "id":{ 
	                  "type":"text"
	               },
	               "site":{ 
	                  "type":"keyword"
	               },
	               "partner":{ 
	                  "type":"keyword"
	               }
	        	}
	    	}
    	}
	}'
`  
 - Step 2  
`go run example/main.go`

 

## Configuration

```go
// Configuration represents config data for elastic client
type Configuration struct {
	// URLs elastic search urls
	URLs []string `json:"ElasticSearchURL"`
	// SearchMaxIdleConns controls the maximum number of idle (keep-alive)
	// connections across all hosts. Zero means no limit.
	SearchMaxIdleConns int `json:"ElasticSearchMaxIdleConns"           default:"100"`
	// SearchIdleConnTimeoutMin is the maximum amount of time an idle
	// (keep-alive) connection will remain idle before closing
	// itself.
	// Zero means no limit.
	SearchIdleConnTimeoutMin int `json:"ElasticSearchIdleConnTimeoutMin"     default:"1"`
}

// BulkConfig represents a config data used by bulkExecute
type BulkConfig struct {
	// Name is an optional name to identify this bulk processor.
	Name string
	// Workers is the number of concurrent workers allowed to be
	// executed. Defaults to 1 and must be greater or equal to 1.
	Workers int
	// FlushInterval specifies when to flush at the end of the given interval.
	// This is disabled by default. If you want the bulk processor to
	// operate completely asynchronously, set both BulkActions and BulkSize to
	// -1 and set the FlushInterval to a meaningful interval.
	FlushInterval time.Duration
	// BulkSize specifies when to flush based on the size (in bytes) of the actions
	// currently added. Defaults to 5 MB and can be set to -1 to be disabled.
	MaxBulkSize int
	// BulkActions specifies when to flush based on the number of actions
	// currently added. Defaults to 1000 and can be set to -1 to be disabled.
	BulkActions int
	// After specifies a function to be executed when bulk requests have been
	// comitted to Elasticsearch. The After callback executes both when the
	// commit was successful as well as on failures.
	AfterCallback elastic.BulkAfterFunc
}

// SearchRequest represents elastic search requests
type SearchRequest struct {
	ElasticIndex   string       // elastic index needed for search
	ElasticType    string       // mapping type that should be manually created
	StartIndex     int          // start search index
	Size           int          // amount of elements for search
	ElasticQuery   ElasticQuery // query to filter search data
	ResultType     reflect.Type // type of search element
	IncludedFields []string     // can be empty in this case all fields will be fetched
}

```

## Example

Example of creating elastic search client

```go
package main

import (
    "fmt"

    "gitlab.connectwisedev.com/platform/Platform-Infrastructure-lib/elasticsearch"
)

func main() {
    client, err := elasticsearch.NewClient(elasticsearch.Configuration{
        URLs:                     []string{"http://localhost:9200"},
        SearchMaxIdleConns:       100,
        SearchIdleConnTimeoutMin: 1,
    })
    if err != nil {
        fmt.Println(err)
        return
    }
    defer func(){
        err = client.Close()
        if err != nil {
            fmt.Println(err)
        }
    }()
}
```
