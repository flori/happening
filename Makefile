DOCKER_IMAGE = happening:$(REVISION_SHORT)
PROJECT_ID = betterplace-183212
POSTGRES_URL ?= "postgresql://happening:happening@localhost:5432/postgres"
REMOTE_TAG = eu.gcr.io/$(PROJECT_ID)/$(DOCKER_IMAGE)
REVISION := $(shell git rev-parse HEAD)
REVISION_SHORT := $(shell echo $(REVISION) | head -c 10)

all: happening happening-server

happening: cmd/happening/main.go *.go
	go build -o happening cmd/happening/main.go

happening-server: cmd/happening-server/main.go *.go
	go build -o happening-server cmd/happening-server/main.go

local: happening-server
	POSTGRES_URL=$(POSTGRES_URL) ./happening-server

fetch:
	go get -u github.com/dgrijalva/jwt-go
	go get -u github.com/labstack/echo
	go get -u github.com/kelseyhightower/envconfig
	go get -u github.com/go-pg/pg

test:
	@go test

clean:
	rm -f happening happening-server

clobber: clean
	rm -rf gospace/src/*

build:
	time docker build -t $(DOCKER_IMAGE) .
	@echo DOCKER_IMAGE="$(DOCKER_IMAGE)"

build-force:
	time docker build -t $(DOCKER_IMAGE) --build-arg FORCE=$(shell date +%s) .
	@echo DOCKER_IMAGE="$(DOCKER_IMAGE)"

debug:
	docker run --rm -it -p 8080:8080 $(DOCKER_IMAGE) bash

server:
	docker run -e POSTGRES_URL=$(POSTGRES_URL) --rm -it -p 8080:8080 $(DOCKER_IMAGE)

push: build
	gcloud auth configure-docker
	docker tag $(DOCKER_IMAGE) $(REMOTE_TAG)
	docker push $(REMOTE_TAG)

pull:
	gcloud auth configure-docker
	docker pull $(REMOTE_TAG)
