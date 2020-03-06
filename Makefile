ifndef GOPATH
	GOPATH := $(HOME)/go
endif

ifndef GOOS
	GOOS := linux
endif

ifndef GO111MODULE
	GO111MODULE := on
endif

all: swagger build

build:
	GOOS=$(GOOS) GO111MODULE=$(GO111MODULE) go build -o bin/filer-gateway internal/main.go

swagger:
	swagger validate pkg/swagger/swagger.yaml
	go generate github.com/Donders-Institute/filer-gateway/internal github.com/Donders-Institute/filer-gateway/pkg/swagger

doc: swagger
	swagger serve pkg/swagger/swagger.yaml

clean:
	rm -rf bin
