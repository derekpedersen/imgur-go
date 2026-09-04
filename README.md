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

This project uses Go modules.

To download and tidy dependencies:

```bash
go mod tidy
```

There is also a `makefile` target for updating project dependencies:

```bash
make dependencies
```

## usage

Create one shared client, then construct package services from that client.

```go
client, err := imgur.NewClient(imgur.Config{
	ClientID: os.Getenv("IMGUR_CLIENT_ID"),
	Mode:     imgur.AuthModeAnonymous,
})
if err != nil {
	return err
}

albumService, err := album.NewService(client)
if err != nil {
	return err
}

imageService, err := images.NewService(client)
if err != nil {
	return err
}

_ = albumService
_ = imageService
```

For OAuth endpoints, create the client with `Mode: imgur.AuthModeOAuth` and a valid `AccessToken`.

## migration notes

Recent refactors removed Java-style service constructors and implementation naming.

- Removed constructors:
	- `album.NewAlbumService(auth, apiURL)`
	- `album.NewAlbumServiceWithClient(client)`
	- `images.NewImageService(auth, apiURL)`
	- `images.NewImageServiceWithClient(client)`
- New constructor pattern:
	- `album.NewService(client)`
	- `images.NewService(client)`
	- `account.NewService(client)`
	- `gallery.NewService(client)`
- Constructor behavior:
	- Service constructors now return `(*Service, error)` and validate that a non-nil client is provided.

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

### table-driven style enforcement

`make test` now runs a style gate before unit/integration tests.

- The gate checks changed `*_test.go` files and requires table-driven structure.
- A changed test file with `Test*` functions must include both `[]struct` test cases and `t.Run(...)` subtests.
- The script is `scripts/check-table-tests.sh` and can be run directly.

By default, the check compares against `origin/master` when available, then falls back to `HEAD~1`.
You can override the comparison base with `TABLE_TEST_BASE_REF`:

```bash
TABLE_TEST_BASE_REF=origin/main bash ./scripts/check-table-tests.sh
```