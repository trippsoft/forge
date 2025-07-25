package function

import (
	"github.com/trippsoft/forge/internal/log"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
)

var (
	SensitiveFunc = function.New(&function.Spec{
		Params: []function.Parameter{
			{
				Name: "value",
				Type: cty.String,
			},
		},
		Type: function.StaticReturnType(cty.String),
		RefineResult: func(rb *cty.RefinementBuilder) *cty.RefinementBuilder {
			return rb.NotNull()
		},
		Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
			value := args[0].AsString()

			log.LogSecretFilter.AddSecret(value)

			return args[0], nil
		},
	})
)

func Sensitive(value cty.Value) (cty.Value, error) {
	return SensitiveFunc.Call([]cty.Value{value})
}
