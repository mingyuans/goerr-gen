project.root=$(shell pwd)


## gen: Generate all necessary files, such as error code files.
.PHONY: gen
gen:
	go run $(project.root)/cmd -output=$(project.root)/example/model/code_gen -docOutput=$(project.root)/example/docs/error_code_generated.md $(project.root)/example/model/code