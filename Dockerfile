ARG BASE_IMAGE

FROM ${BASE_IMAGE} AS build

# Update/Upgrade/Add packages for building

RUN apk add --no-cache bash git go build-base

# Build happening

WORKDIR /build/happening

ADD . .

ENV GOPATH=/build/happening/gospace

RUN make clobber

RUN go get -u github.com/betterplace/go-init

RUN make fetch all

FROM ${BASE_IMAGE}

# Update/Upgrade/Add packages

RUN apk add --no-cache bash ca-certificates

ARG APP_DIR=/app

RUN adduser -h ${APP_DIR} -s /bin/bash -D appuser

RUN mkdir -p /opt/bin

COPY --from=0 --chown=appuser:appuser /build/happening/gospace/bin/go-init /build/happening/happening /build/happening/happening-server /opt/bin/

ENV PATH /opt/bin:${PATH}

EXPOSE 8080

CMD [ "/opt/bin/go-init", "-pre", "/bin/sleep 3", "-main", "/opt/bin/happening-server" ]
