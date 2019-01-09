package codegen

import (
	"io/ioutil"
	"path/filepath"

	"github.com/vektah/gqlparser/ast"
)

func argumentPipe(field *ast.FieldDefinition) string {
	if len(field.Arguments) > 0 {
		//exist arguments, so generate a function type
		return "ctx Context," + argsPipe(field.Arguments)
	}
	//no arguments, return empty string
	return ""
}

func modelFieldTypePipe(field *ast.FieldDefinition) string {
	return typeNamePipe(field.Type)
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
	}
	// array type
	return "[]" + typeNamePipe(t.Elem)
}

func goTypeNamePipe(name string, nonNull bool) string {
	if nonNull {
		return goNonNullTypeNamePipe(name)
	}
	return goNullableTypeNamePipe(name)
}

func goNullableTypeNamePipe(name string) string {
	switch name {
	case "Int":
		return "graphql.IntOption"
	case "Float":
		return "graphql.FloatOption"
	case "String":
		return "graphql.StringOption"
	case "Boolean":
		return "graphql.BoolOption"
	case "ID":
		return "graphql.IDOption"
	default:
		return name + "Option"
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
		return name
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

func isImmediate(t *ast.Type) bool {
	n := t.NamedType
	if n == "Int" || n == "Float" || n == "String" || n == "Boolean" || n == "ID" {
		return true
	}
	return false
}

//return type of resolver's func
func resolverFieldTypePipe(field *ast.FieldDefinition) string {
	t := field.Type
	if t.NamedType != "" {
		//non-array type
		return t.NamedType + "Resolver"
	}
	// array type
	return "[]" + t.Elem.NamedType + "Resolver"
}
func argTypePipe(arg *ast.ArgumentDefinition) string {
	if arg.DefaultValue != nil {
		arg.Type.NonNull = true
	}
	return typeNamePipe(arg.Type)
}
