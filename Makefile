GO_PKGS=$(shell go list ./...)
GOPATH=$(shell echo $$GOPATH)

all: docker-test

test:
ifndef ZOOKEEPER_PATH
	$(info Using default ZOOKEEPER_PATH: '/usr/share/java/')
	export ZOOKEEPER_PATH="/usr/share/java/"
endif
	TESTING="true" go test -cover ${GO_PKGS}

docker-killall:
	docker kill $$(docker ps -aq -f 'ancestor=tevino/zoo')

docker-image:
	docker build --build-arg http_proxy=${http_proxy} --build-arg https_proxy=${https_proxy} -t tevino/zoo .

docker-test: docker-image
	docker run --rm -v '${GOPATH}:/go' tevino/zoo make test

docker-shell: docker-image
	docker run -it --rm -v '${GOPATH}:/go' tevino/zoo /bin/bash

.PHONY: test docker-test docker-image docker-shell docker-killall
