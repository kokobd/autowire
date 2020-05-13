package pkg

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"

	"github.com/pkg/errors"
)

type Symbol struct {
	Kind        SymbolKind
	Name        string
	Annotations []Annotation
	// only for kind == struct
	Fields []SymbolField
}

type SymbolKind uint8

const (
	SymbolKindStruct SymbolKind = iota
	SymbolKindFunction
	SymbolKindVariable
)

type SymbolField struct {
	Name        string
	Annotations []Annotation
}

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
			annotations, err := extractAnnotations(decl.Doc)
			if err != nil {
				return nil, err
			}
			if annotations == nil {
				continue
			}
			symbols = append(symbols, Symbol{
				Name:        decl.Name.Name,
				Kind:        SymbolKindFunction,
				Annotations: annotations,
			})
		case *ast.GenDecl:
			switch decl.Tok {
			case token.CONST:
				fallthrough
			case token.VAR:
				for i := range decl.Specs {
					spec := decl.Specs[i].(*ast.ValueSpec)
					annotations, err := extractAnnotations(spec.Doc)
					if err != nil {
						return nil, err
					}
					if annotations == nil {
						continue
					}
					names := spec.Names
					for j := range names {
						symbols = append(symbols, Symbol{
							Kind:        SymbolKindVariable,
							Name:        names[j].Name,
							Annotations: annotations,
						})
					}
				}
			case token.TYPE:
				for i := range decl.Specs {
					spec := decl.Specs[i].(*ast.TypeSpec)
					annotations, err := extractAnnotations(decl.Doc)
					if err != nil {
						return nil, err
					}
					if annotations == nil {
						continue
					}
					_, ok := spec.Type.(*ast.StructType)
					if !ok {
						continue
					}
					symbols = append(symbols, Symbol{
						Kind:        SymbolKindStruct,
						Name:        spec.Name.Name,
						Annotations: annotations,
					})
				}
			}
		}
	}

	return symbols, nil
}

func extractAnnotations(doc *ast.CommentGroup) ([]Annotation, error) {
	if doc == nil {
		return nil, nil
	}
	text := doc.Text()
	lines := strings.Split(text, "\n")
	var annotations []Annotation
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "@") {
			annotation, err := ParseAnnotation(line)
			if err != nil {
				return nil, err
			}
			if annotation != nil {
				annotations = append(annotations, *annotation)
			}
		}
	}
	return annotations, nil
}
