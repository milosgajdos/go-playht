# go-playht

[![Build Status](https://github.com/milosgajdos/go-playht/actions/workflows/ci.yaml/badge.svg?branch=main)](https://github.com/milosgajdos/go-playht/actions?query=workflow%3ACI)
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/milosgajdos/go-playht)
[![License: Apache-2.0](https://img.shields.io/badge/License-Apache--2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

Unofficial Go module for [play.ht](https://play.ht/) API client.

The official [play.ht](https://play.ht) API documentation, upon which this Go module has been built, can be found [here](https://docs.play.ht/reference).

In order to use this Go module you must create an account with [play.ht](https://play.ht) and generate API secret and retrieve your User ID. See the official docs [here](https://docs.play.ht/reference/api-authentication).

# Get started

Get the module
```shell
go get ./...
```

Run tests:
```shell
go test -v ./...
```

There are a few code samples available in the [examples](./examples) directory so please do have a look. They could give you some idea about how to use this Go module.

> [!IMPORTANT]
> Before you attempt to run the samples you must set a couple of environment variables
> These are automatically read by the client when it gets created; you can override them in your own code.

* `PLAYHT_SECRET_KEY`: API secret key
* `PLAYHT_USER_ID`: Play.HT User ID

## Nix

There is a [Nix flake](https://nixos.wiki/wiki/Flakes) file vailable which lets you work on the Go module using nix.

Just run the following command and you are in the business:
```shell
nix develop
```

# Basics

There are two ways to create audio/speech from the text using the API:
* Job: audio generation is done in async; when you create a job you can monitor its progress via [SSE](https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events)
* Stream: a real-time audio stream available immediately as soon as the stream has been created via the API

The API also allows you to clone a voice using a small sample of limited size. See the [docs](https://docs.play.ht/reference/api-create-instant-voice-clone).
