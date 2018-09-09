package codegen

import (
	"bytes"
	"html/template"
	"strings"

	"github.com/beinan/gql-server/logging"
	"github.com/vektah/gqlparser/ast"
	"github.com/vektah/gqlparser/parser"
)

type GenConfig struct {
	schemaPath       string
	userModelPath    string
	userResolverPath string
}

var logger = logging.StandardLogger(logging.DEBUG)

func Generate(cfg GenConfig) error {
	schemaStr := loadSchema(cfg.schemaPath)
	doc, err := parser.ParseSchema(&ast.Source{
		Name:  "Sche",
		Input: schemaStr})
	if err != nil {
		return err
	}
	//logger.Info(doc.Schema[0], err)
	//logger.Info(doc.Definitions[0])
	// for _, def := range doc.Definitions {
	// 	logger.Info(def.Name, def.Kind)
	// 	for _, field := range def.Fields {
	// 		logger.Info("field", field.Name, field.Type)
	// 		for _, arg := range field.Arguments {
	// 			logger.Info("arg", arg.Name, arg.DefaultValue, arg.Type)
	// 		}
	// 	}
	// }
	logger.Info(generateModel(doc))
	return nil
}

const modelTmpl = `
package genmodel

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

func generateModel(doc *ast.SchemaDocument) string {
	funcMap := template.FuncMap{
		"fieldTypePipe": fieldTypePipe,
		"titlePipe":     strings.Title,
	}
	tmpl, err := template.New("test").Funcs(funcMap).Parse(modelTmpl)
	if err != nil {
		panic(err)
	}
	var buf bytes.Buffer
	tmpl.Execute(&buf, doc)
	return buf.String()
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
	result := "arg struct{ \n"
	for _, arg := range args {
		if arg.DefaultValue != nil {
			arg.Type.NonNull = true
		}
		result += arg.Name + " " + typeNamePipe(arg.Type) + "\n"
	}
	result += "}"
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
		return "int"
	case "Float":
		return "float64"
	case "String":
		return "string"
	case "Boolean":
		return "bool"
	case "ID":
		return "ID"
	default:
		return name
	}
}

func loadSchema(path string) string {
	return `
type User {
  id: ID!
  name: String
  friends(start: Int = 0, pageSize:Int = 20): [User!]
}
type Query {
  getUser(id: ID!): User  
}

schema {
  query: Query
}

`
}
