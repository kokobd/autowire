package pkg

import (
	"github.com/pkg/errors"
	"go/ast"
	"go/parser"
	"go/token"
)

type Symbol struct {
	Kind        SymbolKind
	Name        string
	Annotations []Annotation
}

type SymbolKind uint8

const (
	SymbolKindStruct SymbolKind = iota
	SymbolKindFunction
	SymbolKindVariable
)

func ParseSymbols(fileContent string) ([]Symbol, error) {
	fileSet := token.NewFileSet()
	file, err := parser.ParseFile(fileSet, "", fileContent, parser.ParseComments)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	var symbols []Symbol
	for _, decl_ := range file.Decls {
		switch decl := decl_.(type) {
		case *ast.FuncDecl:
			symbols = append(symbols, Symbol{
				Name:        decl.Name.Name,
				Kind:        SymbolKindFunction,
				Annotations: nil, // TODO
			})
		case *ast.GenDecl:
			switch decl.Tok {
			case token.CONST:
				fallthrough
			case token.VAR:
				for i := range decl.Specs {
					names := decl.Specs[i].(*ast.ValueSpec).Names
					for j := range names {
						symbols = append(symbols, Symbol{
							Kind:        SymbolKindVariable,
							Name:        names[j].Name,
							Annotations: nil, // TODO
						})
					}
				}
			case token.TYPE:
				for i := range decl.Specs {
					spec := decl.Specs[i].(*ast.TypeSpec)
					_, ok := spec.Type.(*ast.StructType)
					if !ok {
						continue
					}
					symbols = append(symbols, Symbol{
						Kind:        SymbolKindStruct,
						Name:        spec.Name.Name,
						Annotations: nil, // TODO
					})
				}
			}
		}
	}

	panic("not implemented")
}
