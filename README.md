# Gateway service for interfacing the DCCN filers

## Build

It requires [Golang](https://golang.org/) version >= 1.13 (i.e. with [Go Modules](https://blog.golang.org/using-go-modules) support) to build the source code.

```bash
$ git clone https://github.com/Donders-Institute/filer-gateway
$ cd filer-gateway
$ make
```

## Run

```bash
$ ./bin/filer-gateway [-v]
```

The HTTP service runs on port `8080`.

## API document

The API document is embedded with the service.  The URL is http://localhost:8080/docs
