package codegen

import (
	"bytes"
	"html/template"
	"io"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/vektah/gqlparser/ast"
)

func GenerateModel(cfg GenConfig, w io.Writer) error {
	return generate(
		cfg, w, executeModelTmpl,
	)
}

const modelTmpl = `
//Generated by gql-server
//DO NOT EDIT
package gen

import "github.com/beinan/gql-server/graphql"

type ID = string
type StringOption = graphql.StringOption

type Context = graphql.Context

{{range .Definitions}}
type {{.Name}} struct {
  {{range .Fields}}
    {{.Name | titlePipe}} {{.| fieldTypePipe}}
  {{end}}
}
{{end}}
`

func executeModelTmpl(doc *ast.SchemaDocument) []byte {
	funcMap := template.FuncMap{
		"fieldTypePipe": fieldTypePipe,
		"titlePipe":     strings.Title,
	}
	tmpl, err := template.New("model").Funcs(funcMap).Parse(modelTmpl)
	if err != nil {
		panic(err)
	}
	var buf bytes.Buffer
	tmpl.Execute(&buf, doc)
	return buf.Bytes()
}

func fieldTypePipe(field *ast.FieldDefinition) string {
	if len(field.Arguments) > 0 {
		//generate a function type
		return "func(ctx Context," + argsPipe(field.Arguments) + ") (" +
			typeNamePipe(field.Type) + ", error)"
	} else {
		return typeNamePipe(field.Type)
	}
}

func argsPipe(args ast.ArgumentDefinitionList) string {
	result := ""
	for _, arg := range args {
		if arg.DefaultValue != nil {
			arg.Type.NonNull = true
		}
		result += arg.Name + " " + typeNamePipe(arg.Type) + ","
	}

	return result
}

func typeNamePipe(t *ast.Type) string {
	if t.NamedType != "" {
		//non-array type
		return goTypeNamePipe(t.NamedType, t.NonNull)
	} else {
		// array type
		return "[]" + typeNamePipe(t.Elem)
	}
}

func goTypeNamePipe(name string, nonNull bool) string {
	if nonNull {
		return goNonNullTypeNamePipe(name)
	} else {
		return goNullableTypeNamePipe(name)
	}
}

func goNullableTypeNamePipe(name string) string {
	switch name {
	case "Int":
		return "IntOption"
	case "Float":
		return "FloatOption"
	case "String":
		return "StringOption"
	case "Boolean":
		return "BoolOption"
	case "ID":
		return "IDOption"
	default:
		return "*" + name
	}
}

func goNonNullTypeNamePipe(name string) string {
	switch name {
	case "Int":
		return "int64"
	case "Float":
		return "float64"
	case "String":
		return "string"
	case "Boolean":
		return "bool"
	case "ID":
		return "ID"
	default:
		return "*" + name
	}
}

func loadSchema(path string) string {
	files, err := filepath.Glob(path + "/*.*")
	if err != nil {
		panic(nil)
	}
	schema := ""
	for _, filename := range files {
		bytes, err := ioutil.ReadFile(filename)
		if err != nil {
			panic(err)
		}
		schema += string(bytes)
	}
	return schema
}