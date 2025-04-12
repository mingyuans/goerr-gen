package cmd

import (
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

// usage is a replacement usage function for the flags package.
func usage() {
	_, _ = fmt.Fprintf(os.Stderr, "usage of goerr-gen:\n")
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

	flag.Usage = usage
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

func Run() {
	log.SetFlags(0)
	log.SetPrefix("goerr-gen: ")

	arg := initArg()

	//codePackagePath := getCodePackagePath()
	var definedErrorCodePackages []ErrorCodePackage
	for _, directory := range arg.DefinedErrorCodeDirectories {
		if isDirectory(directory) {
			codePackages, parseErr := ParsePackage(directory)
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
