package pkg

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go/parser"
	"go/token"
	"testing"
)

func TestParseSource(t *testing.T) {
	type args struct {
		source string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Basic",
			args: args{source: `
package service

import "fmt"

// create new config
// @Component{}
func NewConfig() Config {
	fmt.Println("Initializing config")
	return Config{Host: "127.0.0.1"}
}

type Config struct {
	Host string
}

func NotAProvider() Config {

}

type ISomeService interface {

}

// @Component{AllFields: true}
// @Bind(ISomeService)
type SomeService struct {

}

// @Component{}
func WrapSomeService(s *SomeService) ISomeService {
	return s
}

`},
			want: []string{"NewConfig", "wire.Struct(new(SomeService), \"*\")"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := ParseSource(tt.args.source)
			if assert.NoError(t, err) {
				assert.Equal(t, tt.want, actual)
			}
		})
	}
}

func TestDemoParser(t *testing.T) {
	source := `
package main

func (handler *Handler) AddEntity(id UUID /*@PathParam("id")*/, /*@QueryParam("name")*/ name string) {

}
`
	fileSet := token.NewFileSet()
	tree, err := parser.ParseFile(fileSet, "", source, parser.ParseComments)
	require.NoError(t, err)
	fmt.Printf("%+v\n", tree)
}
