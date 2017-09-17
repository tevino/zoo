FROM zookeeper:3.4.9

### Golang 1.7.4
### https://hub.docker.com/r/library/golang/
RUN apk add --no-cache ca-certificates

ENV GOLANG_VERSION 1.7.4
ENV GOLANG_SRC_URL https://golang.org/dl/go$GOLANG_VERSION.src.tar.gz
ENV GOLANG_SRC_SHA256 4c189111e9ba651a2bb3ee868aa881fab36b2f2da3409e80885ca758a6b614cc

# https://golang.org/issue/14851
# COPY no-pic.patch /
# https://golang.org/issue/17847
# COPY 17847.patch /

RUN set -ex \
	&& apk add --no-cache --virtual .build-deps \
    wget \
		bash \
		gcc \
		musl-dev \
		openssl \
		go \
	\
	&& export GOROOT_BOOTSTRAP="$(go env GOROOT)" \
	\
	&& wget -q "$GOLANG_SRC_URL" -O golang.tar.gz \
	&& echo "$GOLANG_SRC_SHA256  golang.tar.gz" | sha256sum -c - \
	&& tar -C /usr/local -xzf golang.tar.gz \
	&& rm golang.tar.gz \
	&& cd /usr/local/go/src \
	#&& patch -p2 -i /no-pic.patch \
	#&& patch -p2 -i /17847.patch \
	&& ./make.bash \
	\
	#&& rm -rf /*.patch \
	&& apk del .build-deps

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"

### End Golang

RUN apk update
RUN apk add make

### go-zookeeper scans for ZooKeeper here.
RUN mkdir -p /usr/share/java/ \
    && ln -s /zookeeper-*/contrib/fatjar/zookeeper-*-fatjar.jar /usr/share/java/

VOLUME [$GOPATH]
WORKDIR $GOPATH/src/github.com/tevino/zoo
