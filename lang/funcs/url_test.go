package funcs

import (
	"fmt"
	"testing"

	"github.com/zclconf/go-cty/cty"
)

func TestParseURL(t *testing.T) {
	tests := []struct {
		URL  cty.Value
		Want cty.Value
		Err  bool
	}{
		{
			cty.StringVal("https://app.terraform.io/app?q=1"),
			cty.ObjectVal(map[string]cty.Value{
				"scheme":   cty.StringVal("https"),
				"username": cty.NullVal(cty.String),
				"password": cty.NullVal(cty.String),
				"host":     cty.StringVal("app.terraform.io"),
				"path":     cty.StringVal("/app"),
				"query":    cty.StringVal("q=1"),
			}),
			false,
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("parseurl(%#v)", test.URL), func(t *testing.T) {
			got, err := ParseURL(test.URL)

			if test.Err {
				if err == nil {
					t.Fatal("succeeded; want error")
				}
				return
			} else if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			if !got.RawEquals(test.Want) {
				t.Errorf("wrong result\ngot:  %#v\nwant: %#v", got, test.Want)
			}
		})
	}
}

func TestFormatURL(t *testing.T) {
	tests := []struct {
		URL  cty.Value
		Want cty.Value
		Err  bool
	}{
		{
			cty.ObjectVal(map[string]cty.Value{
				"scheme":   cty.StringVal("https"),
				"username": cty.NullVal(cty.String),
				"password": cty.NullVal(cty.String),
				"host":     cty.StringVal("app.terraform.io"),
				"path":     cty.StringVal("/app"),
				"query":    cty.StringVal("q=1"),
			}),
			cty.StringVal("https://app.terraform.io/app?q=1"),
			false,
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("parseurl(%#v)", test.URL), func(t *testing.T) {
			got, err := FormatURL(test.URL)

			if test.Err {
				if err == nil {
					t.Fatal("succeeded; want error")
				}
				return
			} else if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			if !got.RawEquals(test.Want) {
				t.Errorf("wrong result\ngot:  %#v\nwant: %#v", got, test.Want)
			}
		})
	}
}
