SOURCE_DIRS = main.go
export GO111MODULE=on
export LDFLAGS="-s -w"

.PHONY: all
all: gofmt test build dist hash

.PHONY: build
build:
	go build

.PHONY: gofmt
gofmt:
	@test -z $(shell gofmt -l -s $(SOURCE_DIRS) ./ | tee /dev/stderr) || (echo "[WARN] Fix formatting issues with 'make gofmt'" && exit 1)

.PHONY: test
test:
	CGO_ENABLED=0 go test $(shell go list ./... | grep -v /vendor/|xargs echo) -cover

.PHONY: e2e
e2e:
	CGO_ENABLED=0 go test github.com/self-actuated/actions-usage/pkg/get -cover --tags e2e -v

.PHONY: dist
dist:
	mkdir -p bin
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags $(LDFLAGS)  -o bin/actions-usage
	CGO_ENABLED=0 GOOS=darwin go build -ldflags $(LDFLAGS)  -o bin/actions-usage-darwin
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -a -ldflags $(LDFLAGS) -installsuffix cgo -o bin/actions-usage-darwin-arm64
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=6 go build -ldflags $(LDFLAGS)  -o bin/actions-usage-armhf
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags $(LDFLAGS)  -o bin/actions-usage-arm64
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags $(LDFLAGS)  -o bin/actions-usage.exe
