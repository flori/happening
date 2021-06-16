FROM alpine:3.14.0 AS builder

# Update/Upgrade/Add packages for building

RUN apk add --no-cache bash git go build-base yarn

# Build happening

WORKDIR /build

ADD . .

ENV GOPATH=/build/gospace

RUN make clobber

RUN make setup test all

FROM alpine:3.14.0 AS runner

# Update/Upgrade/Add packages

RUN apk add --no-cache bash ca-certificates

ARG APP_DIR=/app

RUN adduser -h ${APP_DIR} -s /bin/bash -D appuser

RUN mkdir -p /opt/bin

COPY --from=builder --chown=appuser:appuser /build/happening /build/happening-server /opt/bin/

COPY --from=builder --chown=appuser:appuser /build/webui/build /webui/build

WORKDIR /

ENV PATH /opt/bin:${PATH}

EXPOSE 8080

CMD [ "/opt/bin/happening-server" ]
