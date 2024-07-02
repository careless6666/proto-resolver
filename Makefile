LOCAL_BIN=$(CURDIR)/bin
PROJECT_NAME=proto-resolver

export GO111MODULE=on

build:
	$(GOENV) CGO_ENABLED=0 go build -v -ldflags "$(LDFLAGS)" -o $(LOCAL_BIN)/$(PROJECT_NAME) ./cmd/proto-resolver

build-windows:
	$(GOENV) CGO_ENABLED=0 env GOOS=windows GOARCH=amd64 go build -v -ldflags "$(LDFLAGS)" -o $(LOCAL_BIN)/$(PROJECT_NAME).exe ./cmd/proto-resolver

build-linux:
	$(GOENV) CGO_ENABLED=0 env GOOS=linux GOARCH=amd64 go build -v -ldflags "$(LDFLAGS)" -o $(LOCAL_BIN)/$(PROJECT_NAME) ./cmd/proto-resolver

build-osx:
	$(GOENV) CGO_ENABLED=0 env GOOS=darwin GOARCH=amd64 go build -v -ldflags "$(LDFLAGS)" -o $(LOCAL_BIN)/$(PROJECT_NAME)-osx ./cmd/proto-resolver

build-osx-m1:
	$(GOENV) CGO_ENABLED=0 env GOOS=darwin GOARCH=arm64 go build -v -ldflags "$(LDFLAGS)" -o $(LOCAL_BIN)/$(PROJECT_NAME)-osx-m1 ./cmd/proto-resolver