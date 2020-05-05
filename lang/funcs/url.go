package funcs

import (
	"net/url"
	"strconv"

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
		"hostname": cty.String,
		"port":     cty.Number,
		"path":     cty.String,
		"query":    cty.Map(cty.Set(cty.String)),
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

		var port cty.Value
		if p := url.Port(); p != "" {
			v, err := strconv.Atoi(p)
			if err != nil {
				return cty.NilVal, err
			}

			port = cty.NumberIntVal(int64(v))
		} else {
			port = cty.NullVal(cty.Number)
		}

		return cty.ObjectVal(map[string]cty.Value{
			"scheme":   cty.StringVal(url.Scheme),
			"username": username,
			"password": password,
			"hostname": cty.StringVal(url.Hostname()),
			"port":     port,
			"path":     cty.StringVal(url.Path),
			"query":    urlValuesToCty(url.Query()),
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
				"hostname": cty.String,
				"port":     cty.Number,
				"path":     cty.String,
				"query":    cty.Map(cty.Set(cty.String)),
			}),
		},
	},
	Type: function.StaticReturnType(cty.String),
	Impl: func(args []cty.Value, retType cty.Type) (ret cty.Value, err error) {
		obj := args[0].AsValueMap()

		var user *url.Userinfo
		if !obj["username"].IsNull() {
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

func urlValuesToCty(query url.Values) cty.Value {
	ret := make(map[string]cty.Value, len(query))
	for key, values := range query {
		cValues := make([]cty.Value, len(values))
		for i, v := range values {
			cValues[i] = cty.StringVal(v)
		}
		ret[key] = cty.SetVal(cValues)
	}

	return cty.MapVal(ret)
}
