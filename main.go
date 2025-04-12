// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// ErrorCodePackage main is a tool to automate the creation of code init function.
// Inspired by `github.com/golang/tools/cmd/stringer`.
package main

import (
	"github.com/mingyuans/goerr-gen/cmd"
)

func main() {
	cmd.Run()
}
