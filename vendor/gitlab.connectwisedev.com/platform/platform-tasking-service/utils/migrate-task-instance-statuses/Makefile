# defining environment variables

TARGET = migrate-task-instance-statuses

# global command
.PHONY: all
all: dependencies migrate-task-instance-statuses

.PHONY: local
local: dependencies localbuild

.PHONY: localbuild
localbuild:
	go build -o $(TARGET)

.PHONY: migrate-task-instance-statuses
migrate-task-instance-statuses:
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
