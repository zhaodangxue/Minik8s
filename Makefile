GOBUILD=go build

DEFAULT_TAGS=default
DEV_TAGS=dev
RELEASE_TAGS=release

all: build

build: 
	make _build TAGS=$(RELEASE_TAGS)

build_dev:
	make _build TAGS=$(DEV_TAGS)

_build: prepare build_targets 
	cp scripts/master_run.sh build/

build_targets: kubelet kubectl apiserver scheduler controller

prepare: deps
	mkdir -p build

deps:
	go mod tidy

kubelet:
	$(GOBUILD) -o build/kubelet -v -tags $(TAGS) kubelet/cmd/server.go

kubectl:
	$(GOBUILD) -o build/kubectl -v -tags $(TAGS) kubectl/run/main.go

apiserver:
	$(GOBUILD) -o build/apiserver -v -tags $(TAGS) apiserver/run/main.go

scheduler:
	$(GOBUILD) -o build/scheduler -v -tags $(TAGS) scheduler/run/main.go

controller:
	$(GOBUILD) -o build/ctlmgr -v -tags $(TAGS) controller/cmd/main.go

.PHONY: all build _build prepare deps 
.PHONY: kubelet kubectl apiserver scheduler controller
