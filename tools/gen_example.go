package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"reflect"
	"strings"

	"github.com/j-b-e/ossie/internal/config"
)

// Generates ossie.toml.example from the Config struct
func generateToml() (string, error) {
	var sb strings.Builder

	configType := reflect.TypeOf(config.Config{})
	configValue := reflect.ValueOf(config.Global)

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "./internal/config/config.go", nil, parser.ParseComments)
	if err != nil {
		return "", err
	}

	// Find the comments and default values for each field
	for _, decl := range f.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok || typeSpec.Name.Name != "Config" {
				continue
			}
			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				continue
			}
			fields := structType.Fields.List
			for i, field := range fields {
				fieldName := configType.Field(i).Name
				comment := field.Comment.Text()
				if comment == "" {
					// ignore uncommented fields, they are not intended for configuration
					continue
				}
				defaultValue := configValue.Field(i)
				if i != 0 {
					// no newline as first line
					sb.WriteString("\n")
				}
				sb.WriteString(fmt.Sprintf("# %s\n", strings.TrimSpace(comment)))

				// Format the default value based on its type
				var tomlValue string
				switch defaultValue.Kind() {
				case reflect.String:
					tomlValue = fmt.Sprintf("\"%s\"", defaultValue.String())
				case reflect.Bool:
					tomlValue = fmt.Sprintf("%t", defaultValue.Bool())
				default:
					fmt.Println("uh oh add default value for this type")
					os.Exit(1)
				}
				sb.WriteString(fmt.Sprintf("%s = %s\n", fieldName, tomlValue))

			}
		}
	}

	return sb.String(), nil
}

func main() {
	tomlContent, err := generateToml()
	if err != nil {
		fmt.Println(err)
		return
	}

	file, err := os.Create("ossie.toml.example")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(tomlContent)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("ossie.toml.example generated")
}
