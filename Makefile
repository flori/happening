DOCKER_IMAGE_LATEST = happening
DOCKER_IMAGE = $(DOCKER_IMAGE_LATEST):$(REVISION_SHORT)
DOCKER_PORT=8080
PROJECT_ID = betterplace-183212
DATABASE_NAME ?= "happening"
POSTGRES_URL ?= "postgresql://flori@dbms:5432/%s?sslmode=disable"
REMOTE_LATEST_TAG := eu.gcr.io/${PROJECT_ID}/$(DOCKER_IMAGE_LATEST)
REMOTE_TAG = eu.gcr.io/$(PROJECT_ID)/$(DOCKER_IMAGE)
REVISION := $(shell git rev-parse HEAD)
REVISION_SHORT := $(shell echo $(REVISION) | head -c 10)
GOPATH ?= gospace

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
	go get -u github.com/go-pg/pg

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
	@rm -rf gospace/src/*

tags: clean
	@gotags -tag-relative=false -silent=true -R=true -f $@ . $(GOPATH)

build:
	time docker build -t $(DOCKER_IMAGE) .
	@echo DOCKER_IMAGE="$(DOCKER_IMAGE)"

build-force:
	time docker build -t $(DOCKER_IMAGE) --no-cache .
	@echo DOCKER_IMAGE="$(DOCKER_IMAGE)"

debug:
	docker run --rm -it $(DOCKER_IMAGE) bash

server:
	docker run -e POSTGRES_URL=$(POSTGRES_URL) --rm -it -p $(DOCKER_PORT):$(DOCKER_PORT) $(DOCKER_IMAGE)

pull:
	gcloud auth configure-docker
	docker pull $(REMOTE_TAG)

push: build
	gcloud auth configure-docker
	docker tag $(DOCKER_IMAGE) $(REMOTE_TAG)
	docker push $(REMOTE_TAG)

push-latest: push
	docker tag ${DOCKER_IMAGE} ${REMOTE_LATEST_TAG}
	docker push ${REMOTE_LATEST_TAG}
