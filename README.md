# Simple Match API

- [api document](/routes.md)

- [system design document](/system-design.md)

## Run the server

```shell
go run main.go
```

## Run the test

```shell
go test
```

## Generate API documentation

This command will print out the api document:

```shell
go main.go -routes
```

## Docker

### Build

```shell
docker build -f Dockerfile -t simple-match .
```

### Run

```shell
docker run -p 8080:8080 simple-match
```

