# Params
PROJECT=micros
IMPORT_PATH=micros/internal
VERSION= $(shell git rev-parse --abbrev-ref HEAD)
TIME= $(shell date '+%Y-%m-%d %H:%M:%S')
BRANCH= $(shell git symbolic-ref --short -q HEAD)
GOVERSION= $(shell go version)
COMMITID= $(shell git log  -1 --pretty=format:"%h")
COMMITDATE= $(shell git show -s --format=%ci)

# Go parameters
GOCMD=go
GORUN=$(GOCMD) run
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BUILDPATH=build
BINPATH=$(BUILDPATH)

FLAGS =  -tags netgo -ldflags='-X "$(IMPORT_PATH)/common.Version=$(VERSION)" -X "$(IMPORT_PATH)/common.BuildTime=$(TIME)" -X "$(IMPORT_PATH)/common.Branch=$(BRANCH)" -X "$(IMPORT_PATH)/common.CommitId=$(COMMITID)" -X "$(IMPORT_PATH)/common.CommitDate=$(COMMITDATE)" -X "$(IMPORT_PATH)/common.GoVersion=$(GOVERSION)"'
TAG = "unknown"

all: test build tag

.PHONY: api
api:
	swag init -g ./cmd/server/main.go  -o api

.PHONY: build
build: api
	@$(GOBUILD) $(FLAGS) -o $(BINPATH)/$(PROJECT) cmd/server/main.go  

.PHONY: tag
tag:
	@docker build -f Dockerfile -t micros:$(TAG) --platform=linux/amd64 --provenance=false --sbom=false . --load

.PHONY: test   
test:
	cp -f configs/conf_dev.ini configs/conf.ini
	go test ./... -parallel=1  -timeout=30m -v

.PHONY: clean
clean:
	rm -f $(BINPATH)/$(PROJECT)
	rm -f $(BINPATH)/init_$(PROJECT)

.PHONY: run
run:
	cp -f configs/conf_dev.ini configs/conf.ini
	@CGO_ENABLE=0 GOOS=darwin GOARCH=amd64 $(GORUN) $(FLAGS) cmd/server/main.go 