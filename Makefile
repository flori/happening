DOCKER_IMAGE_LATEST = happening
DOCKER_IMAGE = $(DOCKER_IMAGE_LATEST):$(REVISION_SHORT)
BASE_IMAGE = $(shell awk '/^FROM .+ AS runner/ { print $$2 }' Dockerfile)
DOCKER_PORT=8080
DATABASE_NAME ?= happening
POSTGRES_URL ?= postgresql://flori@dbms:5432/%s?sslmode=disable
REMOTE_LATEST_TAG := flori303/$(DOCKER_IMAGE_LATEST)
REMOTE_TAG = flori303/$(DOCKER_IMAGE)
REVISION := $(shell git rev-parse HEAD)
REVISION_SHORT := $(shell echo $(REVISION) | head -c 7)
GOPATH := $(shell pwd)/gospace
GOBIN = $(GOPATH)/bin

.EXPORT_ALL_VARIABLES:

all: happening happening-server

happening: cmd/happening/main.go *.go
	go build -o $@ $<

happening-server: cmd/happening-server/main.go *.go
	go build -o $@ $<

local: happening-server
	POSTGRES_URL=$(POSTGRES_URL) ./happening-server

fetch: fake-package
	go get -u github.com/dgrijalva/jwt-go
	go get -u github.com/labstack/echo
	go get -u github.com/kelseyhightower/envconfig
	go get -u github.com/lib/pq
	go get -u github.com/jinzhu/gorm
	go get -u github.com/stretchr/testify
	go get -u github.com/jasonlvhit/gocron
	go get -u github.com/sendgrid/sendgrid-go
	go get -u github.com/go-playground/validator

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
	@echo $(DOCKER_IMAGE)

pull-base:
	docker pull $(BASE_IMAGE)

build: pull-base
	docker build -t $(DOCKER_IMAGE) .
	$(MAKE) build-info

build-force: pull-base
	docker build -t $(DOCKER_IMAGE) --no-cache .
	$(MAKE) build-info

debug:
	docker run --rm -it $(DOCKER_IMAGE) bash

server:
	docker run -e POSTGRES_URL=$(POSTGRES_URL) --rm -it -p $(DOCKER_PORT):$(DOCKER_PORT) $(DOCKER_IMAGE)

pull:
	docker pull $(REMOTE_TAG)
	docker tag $(REMOTE_TAG) $(DOCKER_IMAGE) 

push: build
	docker tag $(DOCKER_IMAGE) $(REMOTE_TAG)
	docker push $(REMOTE_TAG)

push-latest: push
	docker tag ${DOCKER_IMAGE} ${REMOTE_LATEST_TAG}
	docker push ${REMOTE_LATEST_TAG}
