ifndef GOPATH
	GOPATH := $(HOME)/go
endif

ifndef GOOS
	GOOS := linux
endif

ifndef GO111MODULE
	GO111MODULE := on
endif

all: build

build: api-server worker

worker:
	GOOS=$(GOOS) GO111MODULE=$(GO111MODULE) go build -o bin/filer-gateway-worker internal/worker/main.go

api-server:
	GOOS=$(GOOS) GO111MODULE=$(GO111MODULE) go build -o bin/filer-gateway-api internal/api-server/main.go

swagger:
	swagger validate pkg/swagger/swagger.yaml
	go generate github.com/dccn-tg/filer-gateway/internal/api-server github.com/dccn-tg/filer-gateway/pkg/swagger

doc: swagger
	swagger serve pkg/swagger/swagger.yaml

test_worker: build
	GOOS=$(GOOS) GO111MODULE=$(GO111MODULE) go test -v github.com/dccn-tg/filer-gateway/internal/worker/... 

test_api-server: build
	GOOS=$(GOOS) GO111MODULE=$(GO111MODULE) go test -v github.com/dccn-tg/filer-gateway/internal/api-server/... 

clean:
	rm -rf bin
