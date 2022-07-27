[![Latest release](https://badgen.net/github/release/holyhope/god.js)](https://github.com/holyhope/god/releases)
[![GitHub go.mod Go version of a Go module](https://img.shields.io/github/go-mod/go-version/holyhope/god.svg)](https://github.com/holyhope/god)
[![GoDoc reference example](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/holyhope/god)
[![GitHub license](https://img.shields.io/github/license/holyhope/god.svg)](https://github.com/holyhope/god/blob/master/LICENSE)

# Go-d

Caution: This is a work in progress.

This package is a wrapper for the [`systemd`](https://systemd.io) and [`launchd`](https://www.unix.com/man-page/osx/5/launchd.plist/) services manager.

It uses [golang build contraints](https://pkg.go.dev/cmd/go#hdr-Build_constraints) to automatically determine which manager to use.
