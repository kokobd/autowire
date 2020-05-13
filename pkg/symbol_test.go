package pkg

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseSymbols(t *testing.T) {
	type args struct {
		fileContent string
	}
	tests := []struct {
		name    string
		args    args
		want    []Symbol
		wantErr bool
	}{
		{
			name: "Function",
			args: args{
				fileContent: `
package pkg

// ...
// @Component
func NewConfig() Config {

}

// abc
func OtherFunc() {

}
`,
			},
			wantErr: false,
			want: []Symbol{
				{
					Kind:        SymbolKindFunction,
					Name:        "NewConfig",
					Annotations: []Annotation{{Name: "Component", Args: nil}},
				},
			},
		},
		{
			name: "Struct",
			args: args{fileContent: `
package pkg

// ...
// @Component(allFields = false)
type MyService struct {
	// @Autowire
	Dao IMyDao
}
`},
			wantErr: false,
			want: []Symbol{
				{
					Kind: SymbolKindStruct,
					Name: "MyService",
					Annotations: []Annotation{
						{Name: "Component", Args: map[string]string{"allFields": "false"}},
					},
					Fields: []SymbolField{
						{
							Name: "Dao", // TODO
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseSymbols(tt.args.fileContent)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				if assert.NoError(t, err) {
					assert.Equal(t, tt.want, got)
				}
			}
		})
	}
}
