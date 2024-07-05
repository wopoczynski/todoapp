# Todo app

## Name

Todo app

## Description

Todo app backend in go for playing around with base go stuff

Current approach is using [echo](https://echo.labstack.com/) mini framework for handling basing http stuff + [gorm](https://gorm.io/index.html) for database wrapper

## Configuration and verify installation

You may use docker or direct GO runtime (if so go is required locally)

```shell
brew install go
```

for starting the program you may use

```shell
go run ./cmd/main
```

for building executable

```shell
go build ./cmd/main

```

## default settings

| service   | port |
| --------- | ---- |
| webserver | 8123 |
