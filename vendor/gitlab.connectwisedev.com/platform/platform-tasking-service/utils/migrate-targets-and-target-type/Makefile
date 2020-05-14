# defining environment variables

TARGET = platform-tasking-service-migrate-targets-and-target-type

# global command
.PHONY: all
all: dependencies platform-tasking-service-migrate-targets-and-target-type

.PHONY: local
local: dependencies localbuild

.PHONY: localbuild
localbuild:
	go build -o $(TARGET)

.PHONY: platform-tasking-service-migrate-targets-and-target-type
platform-tasking-service-migrate-targets-and-target-type:
	GOOS=linux GOARCH=amd64 go build -o $(TARGET)

.PHONY: dependencies
dependencies:
	echo "Installing dependencies"
	glide install

.PHONY: code-quality
code-quality:
	gometalinter --vendor --tests --skip=mock \
		--exclude='_gen.go' \
		--disable=gotype --disable=errcheck --disable=gas --disable=dupl \
		--deadline=1500s --checkstyle --sort=linter ./... > static-analysis.xml
