package codegen

import (
	"go/format"
	goparser "go/parser"
	"go/token"
	"io"

	"github.com/beinan/gql-server/logging"
	"github.com/vektah/gqlparser/ast"
	"github.com/vektah/gqlparser/parser"
)

//GenConfig is the config for code generation
type GenConfig struct {
	SchemaPath       string
	GenPath          string
	UserModelPath    string
	UserResolverPath string
}

var logger = logging.StandardLogger(logging.DEBUG)

func generate(
	cfg GenConfig,
	w io.Writer,
	execTmpl func(*ast.SchemaDocument) []byte,
) error {
	schemaStr := loadSchema(cfg.SchemaPath)
	doc, err := parser.ParseSchema(&ast.Source{
		Name:  "Sche",
		Input: schemaStr})
	if err != nil {
		return err
	}
	rawCodeBytes := execTmpl(doc)
	//fmt.Println(string(rawCodeBytes))
	fset := token.NewFileSet()
	astNodes, goerr := goparser.ParseFile(fset, "ignore", rawCodeBytes, goparser.ParseComments)
	if goerr != nil {
		return goerr
	}
	return format.Node(w, fset, astNodes)
}
