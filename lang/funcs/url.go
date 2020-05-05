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

		var username, password cty.Value
		if url.User == nil {
			username = cty.NullVal(cty.String)
			password = cty.NullVal(cty.String)
		} else {
			username = cty.StringVal(url.User.Username())
			v, ok := url.User.Password()
			if ok {
				password = cty.StringVal(v)
			} else {
				password = cty.NullVal(cty.String)
			}
		}

		return cty.ObjectVal(map[string]cty.Value{
			"scheme":   cty.StringVal(url.Scheme),
			"username": username,
			"password": password,
			"host":     cty.StringVal(url.Host),
			"path":     cty.StringVal(url.Path),
			"query":    cty.StringVal(url.RawQuery),
		}), nil
	},
})

// FormatURLFunc constructs a function that formats a URL object to a string.
var FormatURLFunc = function.New(&function.Spec{
	Params: []function.Parameter{
		{
			Name: "url",
			Type: cty.Object(map[string]cty.Type{
				"scheme":   cty.String,
				"username": cty.String,
				"password": cty.String,
				"host":     cty.String,
				"path":     cty.String,
				"query":    cty.String,
			}),
		},
	},
	Type: function.StaticReturnType(cty.String),
	Impl: func(args []cty.Value, retType cty.Type) (ret cty.Value, err error) {
		obj := args[0].AsValueMap()

		var user *url.Userinfo
		if !obj["user"].IsNull() {
			if obj["password"].IsNull() {
				user = url.User(obj["username"].AsString())
			} else {
				user = url.UserPassword(obj["username"].AsString(), obj["password"].AsString())
			}
		}

		u := &url.URL{
			Scheme:   obj["scheme"].AsString(),
			User:     user,
			Host:     obj["host"].AsString(),
			Path:     obj["path"].AsString(),
			RawQuery: obj["query"].AsString(),
		}

		return cty.StringVal(u.String()), nil
	},
})

// ParseURL parses a URL, returning a URL object.
func ParseURL(url cty.Value) (cty.Value, error) {
	return ParseURLFunc.Call([]cty.Value{url})
}

// FormatURL formats a URL object, returning a string.
func FormatURL(url cty.Value) (cty.Value, error) {
	return FormatURLFunc.Call([]cty.Value{url})
}
