# Make file to build all Infrastructure Components

## defining environment variables
TARGETPATH=$(GOPATH)/bin

all: dependencies build test-case cover

dependencies:
	glide cc
	rm -f glide.lock
	glide install

build:
	go build -o $(TARGETPATH)/parallel_consumer "./messaging/consumer/example"
	go build -o $(TARGETPATH)/consumer "./messaging/publisher/example"
	go install ./messaging/publisher
	go install ./messaging/publisher/mock
	go install ./messaging/consumer
	go install ./messaging/consumer/mock
	go install ./db
	go install ./db/generator/repo
	go install ./db/goc
	go install ./db/reflect
	go install ./db/template
	go build -o $(TARGETPATH)/db "./db/example/src" 

test-case:
	go test -coverprofile ./messaging/consumer/cover.out -covermode=count ./messaging/consumer
	go test -coverprofile ./messaging/publisher/cover.out -covermode=count ./messaging/publisher
	go test -coverprofile ./db/cover.out -covermode=count ./db

cover:
	echo "mode: count" > cover-all.out
	tail -n +2 ./messaging/consumer/cover.out >> cover-all.out
	tail -n +2 ./messaging/publisher/cover.out >> cover-all.out
	gocover-cobertura < cover-all.out > cover-cobertura.xml

code-quality:
	golangci-lint run -c quality.yaml > static-analysis.xml

mock : messaging-consumer-mock messaging-publisher-mock

messaging-consumer-mock:
	mockgen -package mock gitlab.connectwisedev.com/platform/Platform-Infrastructure-lib/messaging/consumer  Service > messaging/consumer/mock/mocks.go
	mockgen -package mock github.com/bsm/sarama-cluster PartitionConsumer > messaging/consumer/extmock/mock_sarama_cluster.go

messaging-publisher-mock:
	mockgen -package mock gitlab.connectwisedev.com/platform/Platform-Infrastructure-lib/messaging/publisher Producer > messaging/publisher/mock/mocks.go

elastic-mocks:
	mockgen -mock_names=BulkProcessor=MockBulkProcessor -package=esMock gitlab.connectwisedev.com/platform/Platform-Infrastructure-lib/elasticsearch BulkProcessor > elasticsearch/mocks/bulk-processor.go
	mockgen -mock_names=Elastic=MockElastic -package=esMock gitlab.connectwisedev.com/platform/Platform-Infrastructure-lib/elasticsearch Elastic > elasticsearch/mocks/elastic.go
