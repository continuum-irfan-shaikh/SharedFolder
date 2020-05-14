# defining environment variables

TARGETNAME=platform-tasking-service
REPONAME=platform-tasking-service
TARGETPATH=$(GOPATH)/bin/$(TARGETNAME)
SERVICEPATH=$(GOPATH)/src/gitlab.connectwisedev.com/platform/$(REPONAME)/services
SERVICEPACKPATH=$(SERVICEPATH)/continuum
SERVICEBINPATH=$(SERVICEPACKPATH)/taskingService
SERVICECONFIGPATH=$(SERVICEPACKPATH)/config
SERVICEPACKNAME=continuum.zip
SERVICEPACKTARGET=continuum
TRANSLATIONPATH=./translation

BUILDNUMBER = $(BUILD_NUMBER)
BUILDCOMMIT = $(GIT_COMMIT)

LDFLAGBUILDVERSION=-X gitlab.connectwisedev.com/platform/platform-tasking-service/src/config.BuildNumber=$(BUILDNUMBER)
LDFLAGBUILDCOMMIT=-X gitlab.connectwisedev.com/platform/platform-tasking-service/src/config.BuildCommitSHA=$(BUILDCOMMIT)

LINTERCOMMAND=gometalinter

# global command
.PHONY: all
all: dependencies platform-tasking-service cover pack-service

.PHONY: fast-pack
fast-pack: pack-service

.PHONY: local
local: dependencies localbuild test code-quality

.PHONY: localbuild
localbuild:
	go build -ldflags "${LDFLAGBUILDVERSION} ${LDFLAGBUILDCOMMIT}" -o $(TARGETPATH)

.PHONY: platform-tasking-service
platform-tasking-service: bindata
	GOOS=linux GOARCH=amd64 go build -ldflags "${LDFLAGBUILDVERSION} ${LDFLAGBUILDCOMMIT}" -o $(GOPATH)/bin/linux_amd64/$(TARGETNAME)

.PHONY: dependencies
dependencies:
	rm -f glide.lock
	rm -rf vendor
	glide cc
	echo "Installing dependencies"
	glide install
	go get -u github.com/jteeuwen/go-bindata/...

packages = \
    ./usecases/user \
    ./usecases/tasks \
    ./usecases/trigger \
    ./usecases/trigger/handlers \
    ./usecases/scheduler \
    ./usecases/legacy \
    ./handlers/api \
    ./handlers/kafka \
    ./persistency/cassandra \
    ./persistency/memcached \
    ./persistency/kafka \
	./services/common \
	./services/common/errorcode \
	./services \
	./entities \
	./services/tasks \
	./services/encryption \
	./services/task-counters \
	./services/templates \
	./services/task-definitions \
	./services/execution-results \
	./services/execution-results-update \
	./middlewares/transaction-id \
	./middlewares/partner-id \
	./access-control \
	./models \
	./models/statuses \
	./models/triggers \
	./config \
	./app-loader \
	./validator \
	./translation \
	./integration \
    ./types \
    ./repository/asset \
    ./repository/site \
    ./repository/cassandra \
    ./repository/kafka \
    ./repository/triggers \
    ./repository/task-execution \
    ./repository/task-counter-cassandra \
    ./integration \
	./integration/asset \
	./integration/sites \
	./integration/agent \
	./integration/agent-config \
	./integration/dynamic-groups \
	./integration/automation-engine \
	./integration/task-result-webhook \
	./integration/managed-endpoints \
	./logger \
	./middlewares/user \
    ./middlewares/permission \
    ./middlewares/permission_temp \

.PHONY: test
test:
	@$(foreach package,$(packages), \
		set -e; \
		go test -coverprofile $(package)/cover.out -covermode=count $(package);)

.PHONY: generate-mock
generate-mock:
	@echo "Launching mocks code generating..."
	@go generate ./models
	@go generate ./persistency/memcached
	@go generate ./logger
	@go generate ./integration
	@go generate ./repository
	@echo "Successfully finished"

.PHONY: cover-xml
cover-xml:
	@$(foreach package,$(packages), \
	    gocov convert $(package)/cover.out | gocov-xml > $(package)/coverage.xml;)

.PHONY: install-linters
install-linters:
	go get -u github.com/alecthomas/gometalinter
	gometalinter --install --update
	@echo "Installing additional Go linters"
	go get -u github.com/m3db/build-tools/linters/badtime
	go get -u github.com/hypnoglow/durcheck
	go get -u github.com/alexkohler/prealloc

# config file is accessible here https://github.com/ContinuumLLC/platform-gometalinter-config
.PHONY: code-quality-sonar
code-quality-sonar:
	$(LINTERCOMMAND) --checkstyle  ./... > static-analysis.xml

.PHONY: cover
cover: test
	echo "mode: count" > cover-all.out
	@$(foreach package,$(packages), \
		tail -n +2 $(package)/cover.out >> cover-all.out;)
	gocover-cobertura < cover-all.out > cover-cobertura.xml

.PHONY: show
show:
	echo "Launching web browser to show overall coverage..."
	go tool cover -html=cover-all.out

.PHONY: code-quality
code-quality:
	gometalinter --vendor --tests --skip=mock \
		--exclude='_gen.go' --exclude='docs.go' --exclude='vendor/*'\
		--disable=gotype --disable=errcheck --disable=gas --disable=dupl \
		--deadline=1500s --checkstyle --sort=linter ./... > static-analysis.xml

.PHONY: code-quality-print
code-quality-print:
	gometalinter --config=.gometalinter_new.json ./...

.PHONY: pack-service
pack-service:
	mkdir -p $(SERVICEBINPATH)
	mkdir -p $(SERVICECONFIGPATH)
	cp $(GOPATH)/bin/linux_amd64/$(TARGETNAME) $(SERVICEBINPATH)
	cp ../config.json $(SERVICECONFIGPATH)/ctm_tasking_cfg.json
	cd $(SERVICEPATH); rm -f $(SERVICEPACKNAME); zip -r $(SERVICEPACKNAME) $(SERVICEPACKTARGET)

.PHONY: bindata
bindata:
	${GOPATH}/bin/go-bindata -pkg genlang -o $(TRANSLATIONPATH)/genlang/bindata_gen.go $(TRANSLATIONPATH)/languages/

# !!!ATTENTION!!!
# Only for the local usage!
migrate-up:
	migrate -path ../cql/migrations/ -verbose -database cassandra://localhost:9042/platform_tasking_db?consistency=ONE up $(N)

# !!!ATTENTION!!!
# Only for the local usage!
migrate-down:
	migrate -path ../cql/migrations/ -verbose -database cassandra://localhost:9042/platform_tasking_db?consistency=ONE down $(N)

# !!!ATTENTION!!!
# Only for the local usage!
migrate-create:
	migrate -verbose create -ext cql -dir ../cql/migrations/ $(NAME)

# !!!ATTENTION!!!
# Only for the local usage!
migrate-force:
	migrate -verbose -database cassandra://localhost:9042/platform_tasking_db?consistency=ONE -path ../cql/migrations/ force $(V)
