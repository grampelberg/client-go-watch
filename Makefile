
GIT_COMMIT      := $(shell git rev-parse --short HEAD)
CHECKSUM        := $(shell find -s . -type f -not -path "./vendor/*" -not -path "./.git/*" -exec md5sum {} \; | md5sum | awk '{ print $$1 }')
TAG             := ${GIT_COMMIT}-${CHECKSUM}
export IMAGE           := thomasr/client-go-watch:${TAG}

.PHONY: run
run:
	go run main.go

.PHONY: build
build:
	docker build --rm -t ${IMAGE} -f Dockerfile ./

.PHONY: push
push:
	docker push ${IMAGE}

.PHONY: test
test:
	cat deploy.yaml | TYPE=clean  envsubst | kubectl apply -f -
	cat deploy.yaml \
		| TYPE=proxy envsubst \
		| linkerd inject --proxy-log-level=linkerd2_proxy=trace - \
		| kubectl apply -f -
