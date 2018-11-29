# Imgur

A golang package for consuming [imgur](https://imgur.com/) albums.

[![Build Status](https://jenkins.derekpedersen.com/buildStatus/icon?job=derekpedersen/imgur-go/master&style=plastic&.png)](https://jenkins.derekpedersen.com/job/derekpedersen/job/imgur-go/job/master/)
[![Coverage Status](https://coveralls.io/repos/github/derekpedersen/imgur-go/badge.png?branch=master)](https://coveralls.io/github/derekpedersen/imgur-go)

One day I would like to expand to cover the entire [imgur api](https://apidocs.imgur.com/).

## golang

This project is built using golang, if you don't have it installed on your machine you can find the [instructions here](https://golang.org/doc/install).

### main.go

This project is consumed by other projects and isn't an application that is itself deployed, so there is a `main.go` file at the root of the project just to make `golang` happy. I'm sure there is a more elegant solution but for now this is the setup.

## dependencies

Currently the `imgur-go` package relies on `dep` for it's dependency management. If you don't have `dep` installed on your machine just follow the [instructions here](https://github.com/golang/dep#installation).

To use `dep` we must first initialize the project by running the following command:

```bash
dep init
```

After the initialization, when we want to update our depdencies we just run the following command:

```bash
dep ensure
```
Currently there is alreay a `makefile` target for updating the project dependencies:

```bash
make dependencies
```

## build

Since this a golang project if we wanted to build it we could just run the command:

```bash
go build
```

But to make it easier this project has a `makefile` target that handles any additional arguments:

```bash
make build
```

## test

With being a golang project if we just wanted to execute the tests we could run the command:

```bash
go test ./...
```

But to make it easier this project has a `makefile` target that handles the additional arguments and creating a coverage profile:

```bash
make test
```

The coverage profile that is created via `make test` will also include an html webpage that can be used to view the exact lines of code that are covered and not covered. 