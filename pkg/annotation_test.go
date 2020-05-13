package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseAnnotation(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name    string
		args    args
		want    *Annotation
		wantErr bool
	}{
		{
			name: "NoParens",
			args: args{str: "@Component  \n"},
			want: &Annotation{
				Name: "Component",
				Args: nil,
			},
			wantErr: false,
		},
		{
			name: "EmptyParens",
			args: args{str: "@Component()  "},
			want: &Annotation{
				Name: "Component",
				Args: nil,
			},
			wantErr: false,
		},
		{
			name: "DefaultKey",
			args: args{str: "@Component(DefaultValue)"},
			want: &Annotation{
				Name: "Component",
				Args: map[string]string{
					"value": "DefaultValue",
				},
			},
		},
		{
			name: "NamedArgs",
			args: args{str: "@Component(key1 = value1, key2 = value2)"},
			want: &Annotation{
				Name: "Component",
				Args: map[string]string{
					"key1": "value1",
					"key2": "value2",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := ParseAnnotation(tt.args.str)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				if assert.NoError(t, err) {
					assert.Equal(t, tt.want, actual)
				}
			}
		})
	}
}
