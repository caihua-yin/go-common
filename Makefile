REPO_PATH = github.com/caihua-yin/go-common.git

all: check

deps:
	go get -v github.com/golang/lint/golint

check: deps
	go vet ./...
	golint ./...

docker-test-check:
	docker run --rm \
		-v "$$PWD":/go/src/$(REPO_PATH) \
		-w /go/src/$(REPO_PATH) \
		golang:1.7 \
		make check

.PHONY: deps check docker-test-check
