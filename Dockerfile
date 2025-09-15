# stage 0: compile go program
FROM golang:1.20.14-bullseye
RUN mkdir -p /tmp/filer-gateway
WORKDIR /tmp/filer-gateway
ADD internal ./internal
ADD pkg ./pkg
ADD go.mod .
ADD go.sum .
RUN ls -l /tmp/filer-gateway && GOOS=linux go build -a -installsuffix cgo -o bin/filer-gateway-api internal/api-server/main.go
RUN ls -l /tmp/filer-gateway && GOOS=linux go build -a -installsuffix cgo -o bin/filer-gateway-worker internal/worker/main.go

# stage 1: build image for the api-server
FROM almalinux:8 as api-server
RUN ulimit -n 1024 && yum install -y nfs4-acl-tools sssd-client attr acl && yum clean all && rm -rf /var/cache/yum/*
WORKDIR /root
EXPOSE 8080
VOLUME ["/project", "rrd", "/project_cephfs", "/home"]
COPY --from=0 /tmp/filer-gateway/bin/filer-gateway-api .

## entrypoint in shell form so that we can use $PORT environment variable
ENTRYPOINT ["./filer-gateway-api"]

# stage 2: build image for the worker
FROM almalinux:8 as worker
RUN ulimit -n 1024 && yum install -y nfs4-acl-tools sssd-client attr acl && yum clean all && rm -rf /var/cache/yum/*
WORKDIR /root
VOLUME ["/project", "/rrd", "/project_cephfs", "/home"]
COPY --from=0 /tmp/filer-gateway/bin/filer-gateway-worker .

## entrypoint in shell form so that we can use $PORT environment variable
ENTRYPOINT ["./filer-gateway-worker"]
