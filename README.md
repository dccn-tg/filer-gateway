# API gateway for DCCN filers

The goal of this project is to provide an uniform API interface to interact with different filer systems at DCCN for managing storage resource for project and/or user's personal home directory.  Three systems are targeted at the moment (system is checked if the integration with it is implemented):

- NetApp
  * [x] Provision and update project space as ONTAP volume.
  * [x] Provision and update project space as ONTAP qtree.
  * [x] Provision and update home space as ONTAP qtree.
- FreeNAS
  * [x] Provision and update project space as ZFS dataset.
- Ceph
  * [ ] Provision and update project space as CephFS directory.

## Design and components

```
                  +--------------+          +-----------+
                  |              | new task |           |
       POST/PATCH |              +----------> k/v store |
      +----------->              |          |           |
     <------------+              |          +-----+-----+
async (w/ task id)|              |                | run
                  |  API server  |                | task
                  |              |          +-----v-----+             +-----------+
                  |              |          |           |             |           |
          GET     |              |          |           |  set quota  |           |
      +----------->              |          |  workers  +-------------> filer API |
     <------------+              |          |           |  create vol |           |
 sync (w/ data)   |              |          |           |             |           |
                  +----^---------+          +-----+-----+             +-----------+
                       |                          |
       get quota usage | get ACL                  | set ACL
                       |                          |
                  +----+--------------------------v-------+
                  |                                       |
                  |           filer file system           |
                  |                                       |
                  +---------------------------------------+
```

The filer gateway consists of three components:

1. an __API server__ providing RESTful APIs for clients to get and set filer resources for project and user. The APIs are defined with [swagger](https://swagger.io/); and the server-side stubs are generated with [go-swagger](https://goswagger.io/).

1. a __key-value store__ maintaing a list of asynchronous tasks created for POST/PATCH operations that will take too long for API client to wait. It is implemented with [Redis](https://redis.io/).

1. __workers__ are distributable and concurrent processes respondible for executing asynchronous tasks in the background.

Task management are implemented using [Bokchoy](https://github.com/thoas/bokchoy) given that the functionality it provides is just sufficient enough (and not too complicated) for the implementation of the filer gateway.

For NetApp and FreeNAS, two filer interfaces are involved for the gateway to interact with the filer: the exported filesystem, and the filer APIs.  For retriving the ACL and quota usage of a filer directory, only the exported filesystem is used. For setting the ACL, the filesystem is used; while setting quota on the filer requires interaction with the filer administration interface and thus the filer API is utilized.

For CephFS, only the file system interface is used given that all operations (create project directory, set quota and ACL) can be done via the filesystem as long as it is mounted with sufficient privilege.

## Build

It requires [Golang](https://golang.org/) version >= 1.13 (i.e. with [Go Modules](https://blog.golang.org/using-go-modules) support) to build the source code.

```bash
$ git clone https://github.com/Donders-Institute/filer-gateway
$ cd filer-gateway
$ make
```

## Run

Firstly start a Redis server.

Start the API server as

```bash
$ ./bin/filer-gateway-api [-v] [-p 8080] [-r localhost:6379] [-c api-server.yml]
```

By default, the HTTP service runs on port `8080`. Configurations for API authentication is specified by a configuration YAML file.  An example can be found in [api-server.yml](config/api-server.yml).

Start a worker as

```bash
$ ./bin/filer-gateway-worker [-v] [-r localhost:6379] [-c worker.yml]
```

Details of connecting worker to filer API server are specified in a configuration YAML file. An example can be found in [worker.yml](config/worker.yml).

## API document

The API document is embedded with the service.

- for development on localhost, the URL is http://localhost:8080/docs
- for development on the TG playground, the URL is http://dccn-pl001.dccn.nl:8080/docs
- for production, the URL is https://filer-gateway.dccn.nl/docs

## Using Docker

To build the whole service stack (i.e. API server, key-value store, worker), one does:

```bash
$ docker-compose build --force-rm
```

To run the service stack, use

```bash
$ docker-compose -f docker-compose.yml -f docker-compose.dev.yml up
```

The additional compose file `docker-compose.dev.yml` makes the API server exposed to the host network.

## Demo scripts for client implementation

A set of demo scripts is provided in the [demo](demo) folder as a reference for client implementation.  Check the [project.sh](demo/project.sh) for API calls concerning project storage; and [user.sh](demo/user.sh) for the API calls related to user's home storage.
