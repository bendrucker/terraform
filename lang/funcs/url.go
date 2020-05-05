package funcs

import (
	"net/url"

	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
)

// ParseURLFunc constructs a function that parses a given string as a URL.
var ParseURLFunc = function.New(&function.Spec{
	Params: []function.Parameter{
		{
			Name: "url",
			Type: cty.String,
		},
	},
	Type: function.StaticReturnType(cty.Object(map[string]cty.Type{
		"scheme":   cty.String,
		"username": cty.String,
		"password": cty.String,
		"host":     cty.String,
		"path":     cty.String,
		"query":    cty.String,
	})),
	Impl: func(args []cty.Value, retType cty.Type) (ret cty.Value, err error) {
		str := args[0].AsString()

		url, err := url.Parse(str)
		if err != nil {
			return cty.NilVal, err
		}

		var password cty.Value
		v, ok := url.User.Password()
		if ok {
			password = cty.StringVal(v)
		} else {
			password = cty.NullVal(cty.String)
		}

		return cty.ObjectVal(map[string]cty.Value{
			"scheme":   cty.StringVal(url.Scheme),
			"username": cty.StringVal(url.User.Username()),
			"password": password,
			"host":     cty.StringVal(url.Host),
			"path":     cty.StringVal(url.Path),
			"query":    cty.StringVal(url.RawQuery),
		}), nil
	},
})

// ParseURL parses a URL, returning a URL object.
func ParseURL(url cty.Value) (cty.Value, error) {
	return ParseURLFunc.Call([]cty.Value{url})
}
