//go:build tools
// +build tools

// Package tools ensures that `go mod` detect some required dependencies.
//
// This package should not be imported directly.
package tools

import (
	// required by Makefile
	_ "github.com/bufbuild/buf/cmd/buf"
	// build tool
	_ "github.com/buicongtan1997/protoc-gen-swagger-config"
	// required by Makefile
	_ "github.com/daixiang0/gci"
	// required by protoc
	_ "github.com/gogo/protobuf/gogoproto"
	// required by protoc
	_ "github.com/gogo/protobuf/protoc-gen-gogo"
	// required by protoc
	_ "github.com/gogo/protobuf/types"
	// required by protoc
	_ "github.com/golang/protobuf/proto"
	// required by protoc
	_ "github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway"
	// required by Makefile
	_ "github.com/mdomke/git-semver/v5"
	// required by protoc
	_ "github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc"
	// required by protoc
	_ "golang.org/x/tools/cmd/goimports"
	// required by protoc
	_ "google.golang.org/grpc/cmd/protoc-gen-go-grpc"
	// required by protoc
	_ "google.golang.org/protobuf/cmd/protoc-gen-go"
	// required by Makefile
	_ "moul.io/testman"
	// required by Makefile
	_ "mvdan.cc/gofumpt"
)
