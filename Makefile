LOCAL_BIN=$(CURDIR)/bin
PROJECT_NAME=buf-deps

export GO111MODULE=on

build:
	$(GOENV) CGO_ENABLED=0 go build -v -ldflags "$(LDFLAGS)" -o $(LOCAL_BIN)/$(PROJECT_NAME) ./cmd/proto-resolver