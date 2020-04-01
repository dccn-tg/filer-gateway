# Gateway service for interfacing the DCCN filers

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

- an __API server__ providing RESTful APIs for clients to get and set filer resources for project and user. The APIs are defined with [swagger](https://swagger.io/); and the server-side stubs are generated with [go-swagger](https://goswagger.io/).

- a __key-value store__ maintains a list of asynchronous tasks created for POST/PATCH operations that will take too long for API client to wait. It is implemented with [Redis](https://redis.io/).

- __workers__ are distributable and concurrent processes respondible for executing asynchronous tasks in the background.

Task management are implemented using [Bokchoy](https://github.com/thoas/bokchoy) given that the functionality it provides is just sufficient enough (and not too complicated) for the implementation of the filer gateway.

Two filer interfaces are involved for the gateway to interact with the filer: the exported filesystem, and the filer APIs.  For retriving the ACL and quota usage of a filer directory, only the exported filesystem is used. For setting the ACL, the filesystem is used; while setting quota on the filer requires interaction with the filer administration interface and thus the filer API is utilized.

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
$ ./bin/filer-gateway-api [-v] [-p 8080] [-r localhost:6379]
```

The HTTP service runs on port `8080`.

Start a worker as

```bash
$ ./bin/filer-gateway-worker [-v] [-r localhost:6379]
```

## API document

The API document is embedded with the service.  The URL is http://localhost:8080/docs

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