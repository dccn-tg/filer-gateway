# API gateway for DCCN filers

[![Go Report Card](https://goreportcard.com/badge/github.com/Donders-Institute/filer-gateway)](https://goreportcard.com/report/github.com/Donders-Institute/filer-gateway)

The goal of this project is to provide an uniform API interface to interact with different filer systems at DCCN for managing storage resource for project and/or user's personal home directory.

There are three targeted filer systems at the moment:

- NetApp
  * [x] Provision and update project space as ONTAP volume.
  * [x] Provision and update project space as ONTAP qtree.
  * [x] Provision and update home space as ONTAP qtree.
- FreeNAS
  * [x] Provision and update project space as ZFS dataset.
- Ceph
  * [x] Provision and update project space as CephFS directory.

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
       get quota usage | get ACL                  | set ACL (also create dir, set quota for CephFS)
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

For NetApp and FreeNAS, two filer interfaces are used by the gateway to interact with the filer: the exported filesystem, and the filer APIs.  For retriving the ACL and quota usage of a filer directory, only the exported filesystem is used. For setting the ACL, the filesystem is used; while setting quota on the filer requires interaction with the filer administration interface and thus the filer API is utilized.

For CephFS, only the exported filesystem is used given that all operations (create project directory, set quota and ACL) can be done via the filesystem as long as it is mounted with sufficient privilege.

### Projects/Users cache of the API server

Getting quota usage and ACL via the filer's filesystem interface is an expensive operation; and too many `GET` request to the filer API can generate high load on the filesystem, causing performance degradation.

To prevent it from happening, two in-memory caches are implemented in the API server to store filer data (quota and ACL) for all projects and users, respectively.  When the API server receives a `GET` request, it simply returns data from the caches.  This also allows the API server to provide the `GET /projects` and `GET /users` interfaces with good performance.

The caches are updated in two ways:

1. all data in the caches are refreshed every 10 minutes.
1. data concerning a project or a user is updated after changes are made on the filer via the `POST` and `PATCH` APIs.

The second way is implemented with the Redis' [Pub/Sub](https://redis.io/topics/pubsub). Thus, any 3rd-party tool can interact directly with the Redis to trigger cache update for a specific project or user.

For updating the cache of a specific project, one sends a message `{"project": "<number>"}` to the channel `api_pcache_update`.
For updating the cache of a specific user, one sends a message `{"user":"<username>"}` to the channel `api_ucache_update`.

### Statistics

Statistics of filer usage, project storage quota/usage, etc. are exported as the [Prometheus](ttps://prometheus.io) metrics via the `/metrics` endpoint.

## Build

It requires [Golang](https://golang.org/) version >= 1.14 to build the source code.

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

- for development on localhost, the URL is http://localhost:8080/v1/docs
- for development on the TG playground, the URL is http://dccn-pl001.dccn.nl:8080/v1/docs
- for production, the URL is https://filer-gateway.dccn.nl/v1/docs

## Using Docker

A [docker-compose.yml](docker-compose.yml) file is provided to build and deploy the whole service stack (i.e. API server, key-value store, worker) as containers.

To build the containers, one does:

```bash
$ docker-compose build --force-rm
```

To run the containers, one first generates the `env.sh` file using the [print_env.sh](print_env.sh) script, e.g.

```bash
$ ./print_env.sh > env.sh
```

Changing the environment variables defined in the `env.sh` file accordingly, and starting the containers with:

```bash
$ docker-compose -p filer-gateway -f docker-compose.yml up -d
```

or use the provided [start.sh](start.sh) script.

__NOTE:__ ~~In theory, the service stack can also be deployed on a swarm cluster; however, due to the lack of support of the `cap_add` option when using the `docker stack deploy` command (see [this issue](https://github.com/moby/moby/issues/25885)), the CephFS interface will not function properly in the setting/getting the project manager role that is stored in the `trusted.managers` xattr.~~ From Docker version [20.10.0](https://docs.docker.com/engine/release-notes/#20100) onward, the system capability is supported by the Swarm.  Use Docker >= 20.10.0 for deploying the services as a docker stack.

## Demo scripts for client implementation

A set of demo scripts is provided in the [demo](demo) folder as a reference for client implementation.  Check the [project.sh](demo/project.sh) for API calls concerning project storage; and [user.sh](demo/user.sh) for the API calls related to user's home storage.
