TARGETS    ?= darwin/amd64 linux/amd64 windows/amd64
BINARY		 ?= eventing-upgrade-check
GOFILES 	 ?= ./cmd

GO         ?= go
DOCKER     ?= docker
IMAGE_REPO ?= gcr.io/triggermesh
IMAGE      ?= $(IMAGE_REPO)/$(BINARY)
IMAGE_TAG	 ?= v0.1.0



.PHONY: mod-download build test release image push all-in-one

all: build

mod-download:
	$(GO) mod download

build:
	$(GO) build -o $(BINARY) $(GOFILES)

test:
	$(GO) test ./...

release:
	@set -e ; \
	for platform in $(TARGETS); do \
		GOOS=$${platform%/*} ; \
		GOARCH=$${platform#*/} ; \
		RELEASE_BINARY=$(BINARY)-$${GOOS}-$${GOARCH} ; \
		[ $${GOOS} = "windows" ] && RELEASE_BINARY=$${RELEASE_BINARY}.exe ; \
		echo "GOOS=$${GOOS} GOARCH=$${GOARCH} $(GO) build -ldflags "$(LDFLAGS)" -o $(DIST_DIR)$${RELEASE_BINARY} -installsuffix cgo" ; \
		GOOS=$${GOOS} GOARCH=$${GOARCH} $(GO) build -o $(DIST_DIR)$${RELEASE_BINARY} $(GOFILES) ; \
	done

image:
	$(DOCKER) build -t $(IMAGE):$(IMAGE_TAG) .

push: image
	$(DOCKER) push $(IMAGE):$(IMAGE_TAG)

all-in-one:
	mkdir -p ./deploy/
	rm -rf ./deploy/all-in-one.yaml

	for f in ./config/*.yaml; do \
		cat $$f >> ./deploy/all-in-one.yaml ; \
		echo "\r\n---\r\n" >> ./deploy/all-in-one.yaml ; \
	done

	sed -i 's#ko://github.com/triggermesh/eventing-upgrade-check/cmd#$(IMAGE):$(IMAGE_TAG)#' ./deploy/all-in-one.yaml
