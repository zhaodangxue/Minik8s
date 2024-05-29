GOBUILD=go build

DEFAULT_TAGS=default
DEV_TAGS=dev
RELEASE_TAGS=release
BUILDDIR=build
BINDIR=$(BUILDDIR)/bin
YAMLDIR=$(BUILDDIR)/yamls
IMAGEDIR=$(BUILDDIR)/imagebase
DEV_INSTALL_DIR=/tmp/minik8s

all: build

##### Meta Targets #####

build: 
	make _build TAGS=$(RELEASE_TAGS)

build_dev:
	make _build TAGS=$(DEV_TAGS)

install_dev:
	mkdir -p $(DEV_INSTALL_DIR)
	rm -rf $(DEV_INSTALL_DIR)/*
	cp -r build/* $(DEV_INSTALL_DIR)

_build: prepare bin_targets scripts yamls image_base serverless_gateway

prepare: deps
	mkdir -p $(BINDIR)
	mkdir -p $(YAMLDIR)
	mkdir -p $(IMAGEDIR)

deps:
	go mod tidy

clean:
	rm -rf $(BUILDDIR)

.PHONY: all build _build prepare deps 

##### Binaries #####

bin_targets: kubelet kubectl apiserver scheduler controller proxy

.PHONY: kubelet kubectl apiserver scheduler controller proxy

kubelet:
	$(GOBUILD) -o $(BINDIR)/kubelet -v -tags $(TAGS) kubelet/cmd/server.go

kubectl:
	$(GOBUILD) -o $(BINDIR)/kubectl -v -tags $(TAGS) kubectl/run/main.go

apiserver:
	$(GOBUILD) -o $(BINDIR)/apiserver -v -tags $(TAGS) apiserver/run/main.go

scheduler:
	$(GOBUILD) -o $(BINDIR)/scheduler -v -tags $(TAGS) scheduler/run/main.go

proxy:
	$(GOBUILD) -o $(BINDIR)/proxy -v -tags $(TAGS) kubeproxy/run/main.go

controller:
	$(GOBUILD) -o $(BINDIR)/ctlmgr -v -tags $(TAGS) controller/cmd/main.go

serverless_gateway:
	$(GOBUILD) -o $(BINDIR)/sl_gtw -v -tags $(TAGS) serverless/gateway/cmd/server.go

##### Scripts #####

scripts: master_run

.PHONY: scripts master_run

master_run:
	cp scripts/master_run.sh $(BINDIR)

##### Yamls #####

yamls: apiobject_example

.PHONY: yamls apiobject_example

apiobject_example:
	cp -r apiobjects/examples/* $(YAMLDIR)

##### Image Base #####

image_base:
	cp -r serverless/imagebase/* $(IMAGEDIR)

.PHONY: image_base
