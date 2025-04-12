// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// ErrorCodePackage main is a tool to automate the creation of code init function.
// Inspired by `github.com/golang/tools/cmd/stringer`.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

type Arg struct {
	CodeOutput                  string
	DocOutput                   string
	TrimPrefix                  string
	DocTemplate                 string
	DefinedErrorCodeDirectories []string
}

const (
	preservePackageName = "code"
)

// Usage is a replacement usage function for the flags package.
func Usage() {
	_, _ = fmt.Fprintf(os.Stderr, "Usage of goerr-gen:\n")
	_, _ = fmt.Fprintf(os.Stderr, "\tgoerr-gen [flags] -type T directries...\n")
	_, _ = fmt.Fprintf(os.Stderr, "Flags:\n")
	flag.PrintDefaults()
}

func initArg() Arg {
	arg := Arg{}
	flag.StringVar(&arg.CodeOutput, "output", "", "output file packageName; default srcdir/<type>_gen.go")
	flag.StringVar(&arg.DocOutput, "docOutput", "", "output file packageName; default srcdir/<type>_doc.md")
	flag.StringVar(&arg.TrimPrefix, "trim-prefix", "", "trim the `prefix` from the generated constant names")
	flag.StringVar(&arg.DocTemplate, "doc-template", "", "the template file of doc")

	flag.Usage = Usage
	flag.Parse()

	arg.DefinedErrorCodeDirectories = flag.Args()
	if len(arg.DefinedErrorCodeDirectories) == 0 {
		// Default: process whole package in current directory.
		arg.DefinedErrorCodeDirectories = []string{"."}
	}

	arg.CodeOutput = strings.TrimSpace(arg.CodeOutput)
	if len(arg.CodeOutput) == 0 {
		arg.CodeOutput = "."
	}

	arg.DocOutput = strings.TrimSpace(arg.DocOutput)
	if len(arg.DocOutput) == 0 {
		arg.DocOutput = "."
	}

	return arg
}

func main() {
	log.SetFlags(0)
	log.SetPrefix("codegen: ")

	arg := initArg()

	//codePackagePath := getCodePackagePath()
	var definedErrorCodePackages []ErrorCodePackage
	for _, directory := range arg.DefinedErrorCodeDirectories {
		if isDirectory(directory) {
			codePackages, parseErr := parsePackage(directory)
			if parseErr != nil {
				log.Fatalf("parse package error: %v", parseErr)
			}
			definedErrorCodePackages = append(definedErrorCodePackages, codePackages...)
		} else {
			log.Fatalf("not a directory: %s", directory)
		}
	}

	// Generate code file
	genCodeErr := GenerateCodeFile(arg, definedErrorCodePackages)
	if genCodeErr != nil {
		log.Fatalf("generate code file error: %v", genCodeErr)
	}

	genDocErr := GenerateDocs(arg, definedErrorCodePackages)
	if genDocErr != nil {
		log.Fatalf("generate doc file error: %v", genDocErr)
	}

	log.Println("Generate files success")
}

// isDirectory reports whether the named file is a directory.
func isDirectory(name string) bool {
	info, err := os.Stat(name)
	if err != nil {
		log.Fatal(err)
	}

	return info.IsDir()
}

// Generator holds the state of the analysis. Primarily used to buffer
// the output for format.Source.
type Generator struct {
	buf             bytes.Buffer        // Accumulated output.
	sets            []*ErrorCodePackage // ErrorCodePackage we are scanning.
	trimPrefix      string
	codePackagePath string
}

// Printf like fmt.Printf, but add the string to g.buf.
func (g *Generator) Printf(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(&g.buf, format, args...)
}
