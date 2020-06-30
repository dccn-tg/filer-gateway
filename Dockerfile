# stage 0: compile go program
FROM golang:1.13.8
RUN mkdir -p /tmp/filer-gateway
WORKDIR /tmp/filer-gateway
ADD internal ./internal
ADD pkg ./pkg
ADD go.mod .
ADD go.sum .
RUN ls -l /tmp/filer-gateway && GOOS=linux go build -a -installsuffix cgo -o bin/filer-gateway-api internal/api-server/main.go
RUN ls -l /tmp/filer-gateway && GOOS=linux go build -a -installsuffix cgo -o bin/filer-gateway-worker internal/worker/main.go

# stage 1: build image for the api-server
FROM centos:7 as api-server
RUN yum install -y nfs4-acl-tools sssd-client && yum clean all && rm -rf /var/cache/yum/*
WORKDIR /root
EXPOSE 8080
VOLUME ["/project", "project_freenas", "/project_cephfs", "/home"]
COPY --from=0 /tmp/filer-gateway/bin/filer-gateway-api .

## entrypoint in shell form so that we can use $PORT environment variable
ENTRYPOINT ["./filer-gateway-api"]

# stage 2: build image for the worker
FROM centos:7 as worker
RUN yum install -y nfs4-acl-tools sssd-client && yum clean all && rm -rf /var/cache/yum/*
WORKDIR /root
EXPOSE 8080
VOLUME ["/project", "project_freenas", "/project_cephfs", "/home"]
COPY --from=0 /tmp/filer-gateway/bin/filer-gateway-worker .

## entrypoint in shell form so that we can use $PORT environment variable
ENTRYPOINT ["./filer-gateway-worker"]
