FROM ubuntu:16.04 AS build

# Update/Upgrade/Add packages for building

ARG FORCE=$FORCE

RUN apt-get update && apt-get upgrade -y && apt-get install -y \
      git \
      gcc-6-base autoconf build-essential \
      curl \
      golang

# Build dumb-init

WORKDIR /build
RUN git clone https://github.com/Yelp/dumb-init
WORKDIR /build/dumb-init
RUN git checkout b1e978e486114797347deefcc03ab12629a13cc3 # Pinned to Version v1.2.2
RUN make

# Build happening

WORKDIR /build/happening

RUN curl https://storage.googleapis.com/golang/go1.10.linux-amd64.tar.gz | tar xz -C /usr/local

ENV PATH /usr/local/go/bin:${PATH}

ADD . .

ENV GOPATH=/build/happening/gospace

RUN make clean fetch all

FROM ubuntu:16.04

# Update/Upgrade/Add packages

RUN apt-get update && apt-get upgrade -y && apt-get install -y \
  bash ca-certificates

ARG APP_DIR=/app

RUN useradd -d ${APP_DIR} -s /bin/bash appuser

RUN mkdir -p /opt/bin

COPY --from=0 --chown=appuser:appuser /build/dumb-init/dumb-init /build/happening/happening /build/happening/happening-server /opt/bin/

ENV PATH /opt/bin:${PATH}

EXPOSE 8080

CMD [ "/opt/bin/dumb-init", "/opt/bin/happening-server" ]
