package main

import (
	"errors"
	"go/ast"
	"go/token"
	"go/types"
	"golang.org/x/tools/go/packages"
	"regexp"
	"sort"
	"strings"
)

func parsePackage(direPath string) ([]ErrorCodePackage, error) {
	parsedPackages := make([]ErrorCodePackage, 0)

	cfg := &packages.Config{
		Mode:  packages.NeedName | packages.LoadAllSyntax,
		Tests: false,
	}
	pkgs, err := packages.Load(cfg, direPath+"/...")
	if err != nil {
		return parsedPackages, err
	}

	for _, pkg := range pkgs {
		parsedPackage := ErrorCodePackage{
			packageName: pkg.Name,
			packagePath: pkg.PkgPath,
		}
		values, parsedValueErr := getConstFilesByPackage(pkg)
		if parsedValueErr != nil {
			return nil, parsedValueErr
		}

		if len(values) > 0 {
			sort.Slice(values, func(i, j int) bool {
				return values[i].value < values[j].value
			})
		}

		parsedPackage.codes = values
		parsedPackages = append(parsedPackages, parsedPackage)
	}
	return parsedPackages, nil
}

func collectComments(node *ast.File) (map[string]string, error) {
	dict := make(map[string]string)
	ast.Inspect(node, func(n ast.Node) bool {
		genDecl, ok := n.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.CONST {
			return true
		}

		// 遍历 const 声明的规范
		for _, spec := range genDecl.Specs {
			valueSpec, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}

			var comment string
			if valueSpec.Doc != nil && valueSpec.Doc.Text() != "" {
				comment = valueSpec.Doc.Text()
			} else if c := valueSpec.Comment; c != nil && len(c.List) == 1 {
				comment = c.Text()
			}
			comment = strings.TrimSpace(comment)
			dict[valueSpec.Names[0].Name] = comment
		}
		return true
	})
	return dict, nil
}

func collectConstFields(pkg *packages.Package) (map[string]string, error) {
	fields := make(map[string]string)
	if pkg == nil || pkg.TypesInfo == nil {
		return fields, errors.New("pkg.TypesInfo should not be nil")
	}

	if len(pkg.TypesInfo.Defs) == 0 {
		return fields, errors.New("pkg.TypesInfo.Defs should not be empty")
	}

	for key, value := range pkg.TypesInfo.Defs {
		if key.Obj == nil || key.Obj.Kind != ast.Con {
			continue
		}
		constValue := value.(*types.Const).Val().String()
		fields[key.Name] = constValue
	}
	return fields, nil
}

// ErrorCodePackage defines options for package.
type ErrorCodePackage struct {
	packageName string
	packagePath string
	codes       []Value
}

// Value represents a declared constant.
type Value struct {
	comment string
	name    string // The name with trimmed prefix.
	value   string // Will be converted to int64 when needed.
}

func getConstFilesByPackage(pkg *packages.Package) ([]Value, error) {
	values := make([]Value, 0)
	if pkg.Syntax == nil || len(pkg.Syntax) == 0 {
		return values, nil
	}

	comments := make(map[string]string)
	constFields := make(map[string]string)
	for _, node := range pkg.Syntax {
		// 收集注释
		nodeComments, collectCommentsError := collectComments(node)
		if collectCommentsError != nil {
			return nil, collectCommentsError
		}
		for key, value := range nodeComments {
			comments[key] = value
		}

		nodeConstFields, collectConstFieldsError := collectConstFields(pkg)
		if collectConstFieldsError != nil {
			return nil, collectConstFieldsError
		}
		for key, value := range nodeConstFields {
			constFields[key] = value
		}
	}
	for fieldName, value := range constFields {
		fileValue := Value{
			name:  fieldName,
			value: value,
		}
		if comment, ok := comments[fieldName]; ok {
			fileValue.comment = comment
		}
		values = append(values, fileValue)
	}
	return values, nil
}

// ParseComment parse comment to http code and error code description.
func (v *Value) ParseComment() (string, string) {
	reg := regexp.MustCompile(`\w\s*-\s*(\d{3})\s*:\s*([A-Z].*)\s*\.\n*`)
	if !reg.MatchString(v.comment) {
		//log.Printf("constant '%s' have wrong comment codeFileformat, register with 500 as default", v.originalName)
		return "500", "Internal server error"
	}

	groups := reg.FindStringSubmatch(v.comment)
	if len(groups) != 3 {
		return "500", "Internal server error"
	}

	return groups[1], groups[2]
}
