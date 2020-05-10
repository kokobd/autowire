package pkg

import (
	"fmt"
	"github.com/pkg/errors"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

// Parse golang source file to extract provider names
//
// Params:
//
// - source: valid golang source code
//
// Returns: slice of provider names
func ParseSource(source string) ([]string, error) {
	fileSet := token.NewFileSet()
	expr, err := parser.ParseFile(fileSet, "", source, parser.ParseComments)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	var providers []string
	for _, decl := range expr.Decls {
		switch typedDecl := decl.(type) {
		case *ast.FuncDecl:
			commentGroup := typedDecl.Doc
			if commentGroup != nil {
				for _, c := range commentGroup.List {
					begin := strings.Index(c.Text, "@")
					if begin == -1 {
						continue
					}
					annotationStr := c.Text[begin+1:]
					ann, err := parser.ParseExpr(annotationStr)
					if err != nil {
						return nil, errors.WithStack(err)
					}
					if annExpr, ok := ann.(*ast.CallExpr); ok {
						if ident, ok := annExpr.Fun.(*ast.Ident); ok {
							if ident.Name == "Component" {
								providers = append(providers, typedDecl.Name.Name)
							}
						}
					}
				}
			}
		case *ast.GenDecl:
			commentGroup := typedDecl.Doc
			if commentGroup != nil {
				for _, c := range commentGroup.List {
					begin := strings.Index(c.Text, "@")
					if begin == -1 {
						continue
					}
					annotationStr := c.Text[begin+1:]
					ann, err := parser.ParseExpr(annotationStr)
					if err != nil {
						return nil, errors.WithStack(err)
					}
					if annExpr, ok := ann.(*ast.CallExpr); ok {
						if ident, ok := annExpr.Fun.(*ast.Ident); ok {
							fmt.Printf("%+v\n", ident)
						}
					}
				}
			}
		}
	}
	return providers, nil
}

type Bind struct {

}

type ISomeService interface {}

//func parseAnnotation(src string) (Annotation, error) {
//	expr, err := parser.ParseExpr(src)
//	if err != nil {
//		return Annotation{}, errors.WithStack(err)
//	}
//
//	fmt.Printf("%+v\n", expr)
//	return Annotation{
//		Name: "",
//		Args: "",
//	}
//}
