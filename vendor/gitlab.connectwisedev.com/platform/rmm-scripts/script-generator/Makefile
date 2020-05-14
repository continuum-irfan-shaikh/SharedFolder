# defining environment variables
TARGETNAME=script-generator
BUILDCOMMIT = $(GIT_COMMIT)

# global command
.PHONY: all
all: dependencies build cover

.PHONY: build
build:
	go build -o $(TARGETNAME)

.PHONY: dependencies
dependencies:
	echo "Installing dependencies"
	glide install

.PHONY: test
test:
	go test -coverprofile ./cmd/cover.out -covermode=count ./cmd

.PHONY: cover
cover: test
	echo "mode: count" > cover-all.out
	tail -n +2 ./cmd/cover.out >> cover-all.out
	gocover-cobertura < cover-all.out > cover-cobertura.xml


.PHONY: show
show:
	echo "Launching web browser to show overall coverage..."
	go tool cover -html=cover-all.out

.PHONY: code-quality
code-quality:
	gometalinter --vendor --tests --skip=mock \
		--exclude='_gen.go' --exclude='docs.go' \
		--disable=gotype --disable=errcheck --disable=gas --disable=dupl \
		--deadline=1500s --checkstyle --sort=linter ./... > static-analysis.xml
