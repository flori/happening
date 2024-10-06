FROM golang:1.23-alpine AS builder

RUN apk add --no-cache git go build-base npm ca-certificates

# Create appuser.
ENV USER=appuser
ENV UID=10001

RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/none" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"

WORKDIR /build

ADD . .

ENV GOPATH=/build/gospace

RUN make clobber setup test clean webui-build

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags='-w -s' -o happening cmd/happening/main.go

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags='-w -s' -o happening-server cmd/happening-server/main.go

FROM scratch AS runner

COPY --from=builder /etc/passwd /etc/passwd

COPY --from=builder /etc/group /etc/group

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

WORKDIR /

COPY --from=builder --chown=appuser:appuser /build/happening /build/happening-server /

COPY --from=builder --chown=appuser:appuser /build/webui/build /webui/build

USER appuser:appuser

EXPOSE 8080

CMD [ "/happening-server" ]
