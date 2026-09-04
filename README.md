# imgur-go

A Go client for the Imgur API built around a shared client and package-specific services.

[![Build Status](https://jenkins.derekpedersen.com/buildStatus/icon?job=derekpedersen/imgur-go/master&style=plastic&.png)](https://jenkins.derekpedersen.com/job/derekpedersen/job/imgur-go/job/master/)
[![Coverage Status](https://coveralls.io/repos/github/derekpedersen/imgur-go/badge.png?branch=master)](https://coveralls.io/github/derekpedersen/imgur-go)

This repository provides a small Go SDK for the Imgur API. The package uses a single shared `imgur.Client` and exposes domain-specific services for albums, images, gallery, account, and auth-related helpers.

## Overview

The main pattern is intentionally simple:

1. Create one shared client with `imgur.NewClient(...)`.
2. Build a service from that client with `X.NewService(client)`.
3. Call typed methods on the service.

This keeps authentication, base URL handling, and request execution centralized in one place.

## Supported packages

- `imgur` — shared client, auth mode selection, base URL, and HTTP request execution
- `album` — album lookup and mutation operations
- `images` — image lookup, upload, update, delete, and favorite operations
- `gallery` — gallery listing, search, album lookup, and voting
- `account` — account metadata and account-scoped collections
- `authorization` — refresh-token based OAuth helpers
- `imgurtypes` — shared response and token types

This is not a full Imgur API wrapper, but the package structure is built to extend cleanly as more endpoints are added.

## Requirements

- Go 1.25+
- An Imgur Client ID for anonymous requests
- An Imgur access token for OAuth-scoped requests

## Installation

```bash
go get github.com/derekpedersen/imgur-go
```

In a Go module, import the package you need:

```go
import "github.com/derekpedersen/imgur-go/imgur"
```

## Quick start

### Anonymous client

Use anonymous mode for public endpoints that only require a Client ID.

```go
package main

import (
    "fmt"
    "os"

    "github.com/derekpedersen/imgur-go/account"
    "github.com/derekpedersen/imgur-go/album"
    "github.com/derekpedersen/imgur-go/gallery"
    "github.com/derekpedersen/imgur-go/images"
    "github.com/derekpedersen/imgur-go/imgur"
)

func main() {
    client, err := imgur.NewClient(imgur.Config{
        ClientID: os.Getenv("IMGUR_CLIENT_ID"),
        Mode:     imgur.AuthModeAnonymous,
    })
    if err != nil {
        panic(err)
    }

    albumSvc, err := album.NewService(client)
    if err != nil {
        panic(err)
    }

    imageSvc, err := images.NewService(client)
    if err != nil {
        panic(err)
    }

    gallerySvc, err := gallery.NewService(client)
    if err != nil {
        panic(err)
    }

    accountSvc, err := account.NewService(client)
    if err != nil {
        panic(err)
    }

    _ = albumSvc
    _ = imageSvc
    _ = gallerySvc
    _ = accountSvc

    fmt.Println("Imgur client and services initialized")
}
```

### OAuth client

Use OAuth mode when an endpoint requires a user access token.

```go
client, err := imgur.NewClient(imgur.Config{
    AccessToken: os.Getenv("IMGUR_ACCESS_TOKEN"),
    Mode:        imgur.AuthModeOAuth,
})
if err != nil {
    panic(err)
}
```

You can also override the default URL or inject your own HTTP client:

- `BaseURL` — custom API base URL
- `HTTPClient` — custom `*http.Client`

## Common usage patterns

### Service constructors

Each service follows the same constructor pattern:

```go
svc, err := album.NewService(client)
if err != nil {
    return err
}
```

Service constructors validate a non-nil client and return an error if required dependencies are missing.

### Typical calls

```go
albumInfo, err := albumSvc.GetAlbum("abc123")
if err != nil {
    return err
}

imageInfo, err := imageSvc.GetImage("def456")
if err != nil {
    return err
}

items, err := gallerySvc.GetGallery(gallery.ListOptions{
    Section: "hot",
    Sort:    "viral",
    Window:  "day",
    Page:    0,
})
if err != nil {
    return err
}

_ = albumInfo
_ = imageInfo
_ = items
```

## Authorization helper

The `authorization` package includes a refresh-token-based helper for generating OAuth tokens.

```go
auth, err := authorization.NewAuthorization()
if err != nil {
    panic(err)
}

fmt.Println(auth.ClientID)
```

This reads environment variables such as:

- `IMGUR_CLIENT_ID`
- `IMGUR_CLIENT_SECRET`
- `IMGUR_REFRESH_TOKEN`

## Repository conventions

This project intentionally keeps a simple, explicit Go API:

- Shared client at the library boundary
- Service constructors for each feature area
- Errors returned to callers instead of being logged internally
- Table-driven tests with `t.Run(...)` subtests

## Migration note

Older Java-style names are not the supported pattern. The current public API uses the following constructors:

- `album.NewService(client)`
- `images.NewService(client)`
- `account.NewService(client)`
- `gallery.NewService(client)`

Use one shared client and build services from it.

## Build and test

```bash
go build
```

```bash
go test ./...
```

This repository also exposes convenience targets in the Makefile:

```bash
make build
make test
```

The test target creates a coverage profile and an HTML coverage report.

### Table-driven test gate

`make test` runs a project check to ensure changed `*_test.go` files remain in the table-driven style expected by the repository.

The gate requires:

- a slice of test cases (`[]struct`)
- `t.Run(...)` for each case
- a Git-based comparison check for changed files

The script lives at `scripts/check-table-tests.sh` and can be run directly:

```bash
bash ./scripts/check-table-tests.sh
```

Override the comparison base when needed:

```bash
TABLE_TEST_BASE_REF=origin/main bash ./scripts/check-table-tests.sh
```

## Notes for consuming developers and AI agents

- Prefer creating one `imgur.Client` and reusing it across service instances.
- Initialize the client once at startup and pass it into service constructors.
- Match the auth mode to the endpoint: anonymous for public reads, OAuth for user-scoped actions.
- Use the package-level services as the public surface for operations; avoid reaching into implementation details.
- The repository is intentionally package-oriented and straightforward, which makes it easy to add new Imgur resources without introducing a new abstraction layer.

## License

This repository does not currently include a root-level license file. If you plan to redistribute or publish this package, confirm the applicable licensing requirements for your environment before release.
