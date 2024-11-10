//go:build tools
// +build tools

package main

import (
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/google/wire/cmd/wire"
	_ "github.com/goreleaser/goreleaser"
	_ "golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment"
	_ "honnef.co/go/tools/cmd/staticcheck"
)
