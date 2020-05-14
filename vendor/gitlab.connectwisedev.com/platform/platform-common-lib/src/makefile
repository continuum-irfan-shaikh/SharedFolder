# global command
all: dependencies build-all test cover

build-all: api-windows-pdh-build app-build cassandra-build checksum-build cherwell-build circuit-build clar-build communication-build config-build db-build db-sqlite-build distributed-build downloader-build encoding-build entitlement-build env-build exception-build exec-build freecache-build json-build kafka-build cryptomgr-build messaging-build metric-build namedpipes-build plugin-build pluginUtils-build procParser-build redis-client-build retry-build runtime-build servicemanager-build sanitize-build services-build setupUtils-build snowflake-build sync-build testApps-build timeutil-build user-build utils-build uuid-build wal-build web-build webClient-build

dependencies:
	rm -f glide.lock
	rm -rf vendor
	glide cc
	glide install

api-windows-pdh-build:
	GOOS=windows go install ./api/win/pdh

app-build: 
	cp ./testApps/app/versioninfo.json $(GOPATH)/src/gitlab.connectwisedev.com/platform/platform-common-lib/src/app/generate/versioninfo.json
	go generate gitlab.connectwisedev.com/platform/platform-common-lib/src/app/generate
	go build -o $(GOPATH)/src/gitlab.connectwisedev.com/platform/platform-common-lib/src/testApps/app/application ./testApps/app/.
	rm $(GOPATH)/src/gitlab.connectwisedev.com/platform/platform-common-lib/src/app/generate/versioninfo.json
	rm $(GOPATH)/src/gitlab.connectwisedev.com/platform/platform-common-lib/src/testApps/app/application

cassandra-build:
	go install ./cassandra
	go install ./cassandra/mock

checksum-build:
	go install ./checksum
	go install ./checksum/mock

cherwell-build:
	go install ./cherwell
	go install ./cherwell/mapping

circuit-build:
	go install ./circuit
	go install ./circuit/example

clar-build:
	go install ./clar
	go install ./clar/mock

config-build:
	go install ./config
	go install ./config/mock

communication-build:
	go install ./communication/udp
	go install ./communication/udp/mock
	go install ./communication/udp/example
	go install ./communication/http
	go install ./communication/http/client

db-build:
	go install ./db
	go install ./db/mock
	go install ./db/example	
	go install ./db/mssql

db-sqlite-build:
#	go install ./db/sqlite

distributed-build:
	go install ./distributed
	go install ./distributed/leader-election
	go install ./distributed/lock
	go install ./distributed/queue
	go install ./distributed/scheduler
	go install ./distributed/zookeeper

downloader-build:
	go install ./downloader
	go install ./downloader/http
	go install ./downloader/http/example
	go install ./downloader/grab
	go install ./downloader/grab/example
	go install ./downloader/mock

encoding-build:
	go install ./encoding/aes

entitlement-build:
	go install ./entitlement

env-build:
	go install ./env
	go install ./env/mock

exception-build:
	go install ./exception

exec-build:
	go install ./exec
	go install ./exec/mock
	go install ./exec/with

freecache-build:
	go install ./freecache

json-build:
	go install ./json
	go install ./json/mock

kafka-build:
	go install ./kafka
	go install ./kafka/encode
	go install ./kafka/mock

cryptomgr-build:
	go install ./cryptomgr
	go install ./cryptomgr/example
	go install ./cryptomgr/mock

messaging-build:
	go install ./messaging
	go install ./messaging/mock

metric-build:
	go install ./metric
	go install ./metric/example

namedpipes-build:
	go install ./namedpipes
	go install ./namedpipes/mock
	GOOS=windows go install ./namedpipes/npipe
	GOOS=windows go install ./namedpipes/winio

plugin-build:
	go install ./plugin/protocol 
	go install ./plugin/protocol/mock
	go install ./plugin/protocol/http
	GOOS=windows go install ./plugin/wmi 
	GOOS=windows go install ./plugin/wmi/wmiMock 

pluginUtils-build:
	go install ./pluginUtils
	go install ./pluginUtils/mock

procParser-build:
	go install ./procParser
	go install ./procParser/mock

redis-client-build:
	go install ./redis
	go install ./redis/redismock

retry-build:
	go install ./retry

runtime-build:
	go install ./runtime/logger
	go install ./runtime/logger/example
	go install ./runtime/util

sanitize-build:
	go install ./sanitize

servicemanager-build:
	go install ./servicemanager
	GOOS=darwin go install ./servicemanager
	GOOS=windows go install ./servicemanager

services-build:
	go install ./services
	go install ./services/mock
	go install ./services/model

setupUtils-build:
	go install ./setupUtils
	go install ./setupUtils/appManagers

snowflake-build:
	go install ./snowflake

sync-build:
	go install ./sync/zookeeper
	go install ./sync/zookeeper/mock

testApps-build:
#	for file in ./testApps/* ; do go install $$file; done

timeutil-build:
	go install ./timeutil

user-build:
	go install ./user

utils-build:
	go install ./utils

uuid-build:
	go install ./uuid

wal-build:
	go install ./wal
	go install ./wal/example

web-build:
	go install ./web
	go install ./web/mock
	go install ./web/microService
	go install ./web/rest

webClient-build:
	go install ./webClient
	go install ./webClient/mock

cassandra-mockgen:
	mockgen -package mock gitlab.connectwisedev.com/platform/platform-common-lib/src/cassandra DbConnector,Factory,BatchQueryExecutor > cassandra/mock/mocks.go
	mockgen -package mock gitlab.connectwisedev.com/platform/platform-common-lib/src/cassandra/cql Session,Query > cassandra/cql/mock/mocks.go

checksum-mockgen:
	mockgen -package mock gitlab.connectwisedev.com/platform/platform-common-lib/src/checksum Service > checksum/mock/mocks.go

clar-mockgen:
	mockgen	-package mock gitlab.connectwisedev.com/platform/platform-common-lib/src/clar ServiceInit,ServiceInitFactory > clar/mock/mocks.go

communication-mockgen:
	mockgen	-package mock gitlab.connectwisedev.com/platform/platform-common-lib/src/communication/udp Server > communication/udp/mock/mocks.go

config-mockgen:
	mockgen -package mock gitlab.connectwisedev.com/platform/platform-common-lib/src/config ConfigurationService > config/mock/mocks.go

db-mockgen:
	mockgen -package mock gitlab.connectwisedev.com/platform/platform-common-lib/src/db DatabaseProvider > db/mock/mocks.go
db-sqlite-mockgen:
	mockgen -package mock gitlab.connectwisedev.com/platform/platform-common-lib/src/db/sqlite Service > db/sqlite/mock/mocks.go

downloader-mockgen:
	mockgen -package mock gitlab.connectwisedev.com/platform/platform-common-lib/src/downloader Service > downloader/mock/mocks.go

env-mockgen:
	mockgen -package mock gitlab.connectwisedev.com/platform/platform-common-lib/src/env FactoryEnv,Env > env/mock/mocks.go

exec-mockgen:
	mockgen -package mock gitlab.connectwisedev.com/platform/platform-common-lib/src/exec Command > exec/mock/mocks.go

json-mockgen:
	mockgen -package mock gitlab.connectwisedev.com/platform/platform-common-lib/src/json FactoryJSON,DeserializerJSON,SerializerJSON > json/mock/mocks.go

cryptomgr-mockgen:
	mockgen -package mock gitlab.connectwisedev.com/platform/platform-common-lib/src/cryptomgr CryptoMgr > cryptomgr/mock/mocks.go

kafka-mockgen:
	mockgen -package mock gitlab.connectwisedev.com/platform/platform-common-lib/src/kafka ProducerFactory,ConsumerFactory,ProducerService,ConsumerService,Limiter > kafka/mock/mocks.go

messaging-mockgen:
	mockgen -package mock gitlab.connectwisedev.com/platform/platform-common-lib/src/messaging Service > messaging/mock/mocks.go

namedpipes-mockgen:
	mockgen -package mock gitlab.connectwisedev.com/platform/platform-common-lib/src/namedpipes ServerPipe,ClientPipe > namedpipes/mock/mocks.go
	mockgen -package mock net Listener,Conn > namedpipes/mock/nMocks.go

plugin-protocol-mockgen:
	mockgen -package mock gitlab.connectwisedev.com/platform/platform-common-lib/src/plugin/protocol Client,Server,ClientFactory,ServerFactory,ResponseSerializer > plugin/protocol/mock/mocks.go

plugin-wmi-mockgen:
	mockgen -package wmiMock gitlab.connectwisedev.com/platform/platform-common-lib/src/plugin/wmi Wrapper > ./plugin/wmi/wmiMock/mocks_windows.go

pluginUtils-mockgen:
	mockgen -package mock gitlab.connectwisedev.com/platform/platform-common-lib/src/pluginUtils IOReaderFactory,IOWriterFactory > ./pluginUtils/mock/mocks.go

procParser-mockgen:
	mockgen -package mock gitlab.connectwisedev.com/platform/platform-common-lib/src/procParser Parser,ParserFactory > ./procParser/mock/mocks.go

services-mockgen:
	mockgen -package mock gitlab.connectwisedev.com/platform/platform-common-lib/src/services/model VersionFactory,VersionService,VersionDependencies,HealthCheckServiceFactory,HealthCheckService,HealthCheckDalFactory,HealthCheckDal,HealthCheckDependencies > services/mock/mocks.go

sync-mockgen:
	mockgen -package mock gitlab.connectwisedev.com/platform/platform-common-lib/src/sync Service > sync/mock/mocks.go
	mockgen -package mock gitlab.connectwisedev.com/platform/platform-common-lib/src/sync/zookeeper Connection > sync/zookeeper/mock/mocks.go

wal-mockgen:
	mockgen -package mock gitlab.connectwisedev.com/platform/platform-common-lib/src/wal WAL > wal/mock/mocks.go

web-mockgen:
	mockgen -package mock gitlab.connectwisedev.com/platform/platform-common-lib/src/web Server,ServerFactory,Resource,RequestContext,HTTPServer,Router > web/mock/mocks.go

webClient-mockgen:
	mockgen -package mock gitlab.connectwisedev.com/platform/platform-common-lib/src/webClient ClientFactory,ClientService,HTTPClientFactory,HTTPClientService > webClient/mock/mocks.go

redis-mockgen:
	mockgen -package redismock gitlab.connectwisedev.com/platform/platform-common-lib/src/redis Client,Pipeliner > redis/redismock/mocks.go


integration-test:
	go test -tags=integration -coverprofile ./cassandra/cover.out -covermode=count ./cassandra

test:
#	go test -coverprofile ./api/win/pdh/cover.out -covermode=count ./api/win/pdh
	go test -coverprofile ./app/cover.out -covermode=count ./app
	go test -coverprofile ./cassandra/cover.out -covermode=count ./cassandra
	go test -coverprofile ./cassandra/cql/cover.out -covermode=count ./cassandra/cql
	go test -coverprofile ./checksum/cover.out -covermode=count ./checksum
	go test -coverprofile ./cherwell/cover.out -covermode=count ./cherwell
	go test -coverprofile ./cherwell/mapping/cover.out -covermode=count ./cherwell/mapping
	go test -coverprofile ./circuit/cover.out -covermode=count ./circuit
	go test -coverprofile ./clar/cover.out -covermode=count ./clar
	go test -coverprofile ./config/cover.out -covermode=count ./config
#	go test -coverprofile ./communication -covermode=count ./communication
	go test -coverprofile ./communication/udp/cover.out -covermode=count ./communication/udp
	go test -coverprofile ./communication/udp/example/cover.out -covermode=count ./communication/udp/example
	go test -coverprofile ./communication/http/cover.out -covermode=count ./communication/http
	go test -coverprofile ./communication/http/client/cover.out -covermode=count ./communication/http/client
	go test -coverprofile ./db/cover.out -covermode=count ./db
	go test -coverprofile ./db/mssql/cover.out -covermode=count ./db/mssql
	go test -coverprofile ./db/sqlite/cover.out -covermode=count ./db/sqlite
	go test -coverprofile ./distributed/cover.out -covermode=count ./distributed
	go test -coverprofile ./distributed/leader-election/cover.out -covermode=count ./distributed/leader-election
	go test -coverprofile ./distributed/lock/cover.out -covermode=count ./distributed/lock
	go test -coverprofile ./distributed/queue/cover.out -covermode=count ./distributed/queue
	go test -coverprofile ./distributed/scheduler/cover.out -covermode=count ./distributed/scheduler
	go test -coverprofile ./distributed/zookeeper/cover.out -covermode=count ./distributed/zookeeper
	go test -coverprofile ./downloader/cover.out -covermode=count ./downloader
	go test -coverprofile ./downloader/http/cover.out -covermode=count ./downloader/http
	go test -coverprofile ./downloader/grab/cover.out -covermode=count ./downloader/grab
	go test -coverprofile ./encoding/aes/cover.out -covermode=count ./encoding/aes
	go test -coverprofile ./entitlement/cover.out -covermode=count ./entitlement
	go test -coverprofile ./env/cover.out -covermode=count ./env
	go test -coverprofile ./exception/cover.out -covermode=count ./exception
	go test -coverprofile ./exec/cover.out -covermode=count ./exec
	go test -coverprofile ./exec/with/cover.out -covermode=count ./exec/with
	go test -coverprofile ./freecache/cover.out -covermode=count ./freecache
	go test -coverprofile ./json/cover.out -covermode=count ./json
	go test -coverprofile ./kafka/cover.out -covermode=count ./kafka
	go test -coverprofile ./kafka/encode/cover.out -covermode=count ./kafka/encode
	go test -coverprofile ./cryptomgr/cover.out -covermode=count ./cryptomgr
	go test -coverprofile ./messaging/cover.out -covermode=count ./messaging
	go test -coverprofile ./metric/cover.out -covermode=count ./metric
	go test -coverprofile ./metric/example/cover.out -covermode=count ./metric/example
	go test -coverprofile ./namedpipes/cover.out -covermode=count ./namedpipes
#	go test -coverprofile ./namedpipes/npipe/cover.out -covermode=count ./namedpipes/npipe
#	go test -coverprofile ./namedpipes/winio/cover.out -covermode=count ./namedpipes/winio
	go test -coverprofile ./notifications/email/cover.out -covermode=count ./notifications/email
	go test -coverprofile ./plugin/protocol/cover.out -covermode=count ./plugin/protocol
	go test -coverprofile ./plugin/protocol/http/cover.out -covermode=count ./plugin/protocol/http
#	go test -coverprofile ./plugin/wmi/cover.out -covermode=count ./plugin/wmi
	go test -coverprofile ./pluginUtils/cover.out -covermode=count ./pluginUtils
	go test -coverprofile ./procParser/cover.out -covermode=count ./procParser
	go test -coverprofile ./redis/cover.out -covermode=count ./redis
	go test -coverprofile ./retry/cover.out -covermode=count ./retry
	go test -coverprofile ./runtime/logger/cover.out -covermode=count ./runtime/logger
	go test -coverprofile ./runtime/logger/example/cover.out -covermode=count ./runtime/logger/example
	go test -coverprofile ./runtime/util/cover.out -covermode=count ./runtime/util
	go test -coverprofile ./sanitize/cover.out -covermode=count ./sanitize
	go test -coverprofile ./servicemanager/cover.out -covermode=count ./servicemanager
#	go test -coverprofile ./services/cover.out -covermode=count ./services
	go test -coverprofile ./services/model/cover.out -covermode=count ./services/model
	go test -coverprofile ./setupUtils/cover.out -covermode=count ./setupUtils
	go test -coverprofile ./setupUtils/appManagers/cover.out -covermode=count ./setupUtils/appManagers
	go test -coverprofile ./snowflake/cover.out -covermode=count ./snowflake
	go test -coverprofile ./sync/zookeeper/cover.out -covermode=count ./sync/zookeeper
	go test -coverprofile ./timeutil/cover.out -covermode=count ./timeutil
	go test -coverprofile ./user/cover.out -covermode=count ./user
	go test -coverprofile ./utils/cover.out -covermode=count ./utils
	go test -coverprofile ./uuid/cover.out -covermode=count ./uuid 
	go test -coverprofile ./validate/is/cover.out -covermode=count ./validate/is 
	go test -coverprofile ./wal/cover.out -covermode=count ./wal
	go test -coverprofile ./wal/example/cover.out -covermode=count ./wal/example
	go test -coverprofile ./web/cover.out -covermode=count ./web
	go test -coverprofile ./web/microService/cover.out -covermode=count ./web/microService
	go test -coverprofile ./web/rest/cover.out -covermode=count ./web/rest
	go test -coverprofile ./webClient/cover.out -covermode=count ./webClient
	go test -coverprofile ./testUtils/cover.out -covermode=count ./testUtils


cover:
	echo "mode: count" > cover-all.out
	tail -n +2 ./cassandra/cover.out >> cover-all.out
	tail -n +2 ./cassandra/cql/cover.out >> cover-all.out
	tail -n +2 ./checksum/cover.out >> cover-all.out
	tail -n +2 ./cherwell/cover.out >> cover-all.out
	tail -n +2 ./cherwell/mapping/cover.out >> cover-all.out
	tail -n +2 ./circuit/cover.out >> cover-all.out
	tail -n +2 ./clar/cover.out >> cover-all.out
	tail -n +2 ./communication/udp/cover.out >> cover-all.out
	tail -n +2 ./communication/udp/example/cover.out >> cover-all.out
	tail -n +2 ./communication/http/cover.out >> cover-all.out
	tail -n +2 ./communication/http/client/cover.out >> cover-all.out
	tail -n +2 ./db/cover.out >> cover-all.out
	tail -n +2 ./db/mssql/cover.out >> cover-all.out
	tail -n +2 ./distributed/zookeeper/cover.out >> cover-all.out
	tail -n +2 ./downloader/http/cover.out >> cover-all.out
	tail -n +2 ./downloader/grab/cover.out >> cover-all.out
	tail -n +2 ./entitlement/cover.out >> cover-all.out
	tail -n +2 ./exception/cover.out >> cover-all.out
	tail -n +2 ./exec/with/cover.out >> cover-all.out
	tail -n +2 ./json/cover.out >> cover-all.out
	tail -n +2 ./cryptomgr/cover.out >> cover-all.out
	tail -n +2 ./messaging/cover.out >> cover-all.out	
	tail -n +2 ./metric/cover.out >> cover-all.out	
	tail -n +2 ./metric/example/cover.out >> cover-all.out	
	tail -n +2 ./namedpipes/cover.out >> cover-all.out
	tail -n +2 ./notifications/email/cover.out >> cover-all.out
	tail -n +2 ./plugin/protocol/http/cover.out >> cover-all.out
	tail -n +2 ./procParser/cover.out >> cover-all.out
	tail -n +2 ./redis/cover.out >> cover-all.out
	tail -n +2 ./retry/cover.out >> cover-all.out	
	tail -n +2 ./runtime/logger/cover.out >> cover-all.out
	tail -n +2 ./runtime/logger/example/cover.out >> cover-all.out
	tail -n +2 ./runtime/util/cover.out >> cover-all.out
	tail -n +2 ./sanitize/cover.out >> cover-all.out
	tail -n +2 ./sync/zookeeper/cover.out >> cover-all.out
	tail -n +2 ./timeutil/cover.out >> cover-all.out
	tail -n +2 ./utils/cover.out >> cover-all.out
	tail -n +2 ./uuid/cover.out >> cover-all.out
	tail -n +2 ./validate/is/cover.out >> cover-all.out
	tail -n +2 ./wal/cover.out >> cover-all.out
	tail -n +2 ./wal/example/cover.out >> cover-all.out
	tail -n +2 ./web/cover.out >> cover-all.out
	tail -n +2 ./web/microService/cover.out >> cover-all.out
	tail -n +2 ./web/rest/cover.out >> cover-all.out
	tail -n +2 ./webClient/cover.out >> cover-all.out
	tail -n +2 ./testUtils/cover.out >> cover-all.out
	
	gocover-cobertura < cover-all.out > cover-cobertura.xml

.PHONY: show
show:
	echo "Launching web browser to show overall coverage..."
	go tool cover -html=cover-all.out

.PHONY: code-quality
code-quality:
	golangci-lint run -c quality.yaml > static-analysis.xml