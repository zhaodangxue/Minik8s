GOBUILD=go build

DEFAULT_TAGS=default
DEV_TAGS=dev
RELEASE_TAGS=release
BUILDDIR=build
BINDIR=$(BUILDDIR)/bin
YAMLDIR=$(BUILDDIR)/yamls
IMAGEDIR=$(BUILDDIR)/imagebase
JOBDIR=$(BUILDDIR)/job
FUNCTIONDIR=$(BUILDDIR)/functions
DEV_INSTALL_DIR=/tmp/minik8s

all: build

##### Meta Targets #####

build: 
	make _build TAGS=$(RELEASE_TAGS)

build_cgo0:
	make build CGO_ENABLED=0

build_dev:
	make _build TAGS=$(DEV_TAGS)

install_dev:
	mkdir -p $(DEV_INSTALL_DIR)
	rm -rf $(DEV_INSTALL_DIR)/*
	cp -r build/* $(DEV_INSTALL_DIR)

deploy_no_bin:
	echo "Deploying application..."
	./scripts/deploy_to_master.sh no_bin
	./scripts/deploy_to_worker.sh no_bin
	echo "Application successfully deployed."

deploy: clean build_cgo0
	echo "Deploying application..."
	./scripts/deploy_to_master.sh
	./scripts/deploy_to_worker.sh
	echo "Application successfully deployed."

_build: prepare bin_targets scripts yamls serverless job

prepare: deps
	mkdir -p $(BINDIR)
	mkdir -p $(YAMLDIR)
	mkdir -p $(IMAGEDIR)
	mkdir -p $(FUNCTIONDIR)
	mkdir -p $(JOBDIR)

deps:
	go mod tidy

clean:
	rm -rf $(BUILDDIR)

.PHONY: all build _build prepare deps clean install_dev deploy_dev build_dev deploy_no_bin deploy build_cgo0

##### Binaries #####

bin_targets: kubelet kubectl apiserver scheduler controller proxy serverless_gateway jobserver

.PHONY: kubelet kubectl apiserver scheduler controller proxy jobserver

kubelet:
	$(GOBUILD) -o $(BINDIR)/kubelet -v -tags $(TAGS) kubelet/cmd/server.go

kubectl:
	$(GOBUILD) -o $(BINDIR)/kubectl -v -tags $(TAGS) kubectl/run/main.go

apiserver:
	$(GOBUILD) -o $(BINDIR)/apiserver -v -tags $(TAGS) apiserver/run/main.go

scheduler:
	$(GOBUILD) -o $(BINDIR)/scheduler -v -tags $(TAGS) scheduler/run/main.go

proxy:
	$(GOBUILD) -o $(BINDIR)/kubeproxy -v -tags $(TAGS) kubeproxy/run/main.go

controller:
	$(GOBUILD) -o $(BINDIR)/ctlmgr -v -tags $(TAGS) controller/cmd/main.go

serverless_gateway:
	$(GOBUILD) -o $(BINDIR)/sl_gtw -v -tags $(TAGS) serverless/gateway/cmd/server.go

jobserver:
	$(GOBUILD) -o $(BINDIR)/jobserver -v -tags $(TAGS) jobserver/cmd/server.go

##### Scripts #####

scripts: master_run worker_run

.PHONY: scripts master_run worker_run

master_run:
	cp scripts/master_run.sh $(BINDIR)

worker_run:
	cp scripts/worker_run.sh $(BINDIR)

##### Yamls #####

yamls: apiobject_example

.PHONY: yamls apiobject_example

apiobject_example:
	cp -r apiobjects/examples/* $(YAMLDIR)

##### Serverless #####

serverless: serverless_examples image_base

serverless_examples:
	cp -r serverless/examples/* $(FUNCTIONDIR)

image_base:
	cp -r serverless/imagebase/* $(IMAGEDIR)

.PHONY: serverless serverless_examples image_base

##### Job #####

job: job_function job_image_base

job_image_base:
	mkdir -p $(JOBDIR)/imagebase
	cp -r jobserver/imagebase/imagedata $(JOBDIR)/imagebase/

job_function:
	mkdir -p $(JOBDIR)/functions
	cp -r jobserver/examples/* $(JOBDIR)/functions

.PHONY: job_function job job_image_base