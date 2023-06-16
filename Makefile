DOCKER_IMAGE_LATEST = happening
DOCKER_IMAGE = $(DOCKER_IMAGE_LATEST):$(REVISION_SHORT)
BASE_IMAGE = $(shell awk '/^FROM .+ AS runner/ { print $$2 }' Dockerfile)
DOCKER_PORT=8080
POSTGRES_URL ?= postgresql://postgres:secret@localhost:6666/happening?sslmode=disable
REMOTE_LATEST_TAG := flori303/$(DOCKER_IMAGE_LATEST)
REMOTE_TAG = flori303/$(DOCKER_IMAGE)
REVISION := $(shell git rev-parse HEAD)
REVISION_SHORT := $(shell echo $(REVISION) | head -c 7)
GOPATH := $(shell pwd)/gospace
GOBIN = $(GOPATH)/bin
WEBUI_DIR := $(shell pwd)/webui
HAPPENING_SERVER_URL ?= http://localhost:8080
HTTP_AUTH ?= admin:test1234
SIGNING_SECRET ?= secret
DOCKER_BUILDKIT = 1

.EXPORT_ALL_VARIABLES:

all: webui-build happening happening-server

happening: cmd/happening/main.go *.go
	go build -o $@ $<

happening-server: cmd/happening-server/main.go *.go
	go build -o $@ $<

local: happening-server
	./happening-server

psql:
	PGPASSWORD=secret psql -U postgres -h localhost -p 6666 -d happening

webui-build:
	cd webui && yarn --network-timeout 1000000 --network-concurrency 4 && yarn build

webui-start:
	cd webui && REACT_APP_HAPPENING_SERVER_URL=$(HAPPENING_SERVER_URL) yarn start

webui-browserslist-update:
	cd webui && npx browserslist@latest --update-db

fetch:
	go mod download

setup: fake-package fetch

fake-package:
	rm -rf $(GOPATH)/src/github.com/flori/happening
	mkdir -p $(GOPATH)/src/github.com/flori
	ln -s $(shell pwd) $(GOPATH)/src/github.com/flori/happening

test:
	@go test

coverage:
	@go test -coverprofile=coverage.out

coverage-display: coverage
	@go tool cover -html=coverage.out

clean:
	@rm -f happening happening-server coverage.out tags

clobber: clean
	@rm -rf $(GOPATH)/*

tags: clean
	@gotags -tag-relative=false -silent=true -R=true -f $@ . $(GOPATH)

build-info:
	@echo $(REMOTE_TAG)

build:
	docker build --progress=plain --pull -t $(DOCKER_IMAGE) -t $(DOCKER_IMAGE_LATEST) .
	$(MAKE) build-info

build-force:
	docker build --progress=plain --pull -t $(DOCKER_IMAGE) -t $(DOCKER_IMAGE_LATEST) --no-cache .
	$(MAKE) build-info

server:
	docker run --network=host -e SIGNING_SECRET -e HTTP_AUTH -e POSTGRES_URL=$(POSTGRES_URL) -e HAPPENING_SERVER_URL=$(HAPPENING_SERVER_URL) -it -p $(DOCKER_PORT):$(DOCKER_PORT) $(DOCKER_IMAGE)

push: build
	docker tag $(DOCKER_IMAGE) $(REMOTE_TAG)
	docker push $(REMOTE_TAG)

push-latest: push
	docker tag ${DOCKER_IMAGE} ${REMOTE_LATEST_TAG}
	docker push ${REMOTE_LATEST_TAG}

git-tag:
	git tag tag-$(REVISION_SHORT) $(REVISION)
	git push origin
	git push origin tag-$(REVISION_SHORT)
