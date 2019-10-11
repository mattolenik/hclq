package queryast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var tests = []struct {
	query       string
	description string
	expected    *Expr
}{
	{".ab.de-f.xz[ab] | fn1(.one, .two)", "simple query", &Expr{
		Node: &Path{
			Crumbs: []*Crumb{
				&Crumb{
					Key: &Key{
						Ident:    "ab",
						Selector: nil,
					},
				},
				&Crumb{
					Key: &Key{
						Ident:    "de-f",
						Selector: nil,
					},
				},
				&Crumb{
					Key: &Key{
						Ident: "xz",
						Selector: &MapSelector{
							Key: "ab",
						},
					},
				},
			},
		},
		Next: &Expr{
			Node: &FunctionCall{
				Name: "fn1",
				Params: []interface{}{
					&Expr{
						Node: &Path{
							Crumbs: []*Crumb{
								&Crumb{
									Key: &Key{
										Ident:    "one",
										Selector: nil,
									},
								},
							},
						},
						Next: nil,
					},
					&Expr{
						Node: &Path{
							Crumbs: []*Crumb{
								&Crumb{
									Key: &Key{
										Ident:    "two",
										Selector: nil,
									},
								},
							},
						},
						Next: nil,
					},
				},
			},
			Next: nil,
		},
	},
	},
}

func TestGrammar(t *testing.T) {
	for _, test := range tests {
		result, err := Parse("inline", []byte(test.query))
		assert.Nil(t, err)
		assert.Equal(t, result, test.expected)
	}
}
