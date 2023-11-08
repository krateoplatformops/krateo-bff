package tmpl

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestJQTemplate(t *testing.T) {
	const query = `{{ .hobbies | join(",") }}`

	ds, err := dataSource()
	if err != nil {
		t.Fatal(err)
	}

	tpl, err := New()
	if err != nil {
		t.Fatal(err)
	}

	res, err := tpl.Execute(query, ds)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(res)
}

func TestJQ(t *testing.T) {
	ds, err := dataSource()
	if err != nil {
		t.Fatal(err)
	}

	str, err := jq(`.hobbies | join(",")`, ds)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(str)
}

func TestFixQuery(t *testing.T) {
	test := []struct {
		input string
		want  string
	}{
		{`{{ .age }}`, `.age`},
		{` .age }}`, ` .age }}`},
		{`{{ .location.city }}`, `.location.city`},
		{`hello world`, `hello world`},
		{`{{ .hobbies | join(",") }}`, `.hobbies | join(",")`},
	}

	tpl, err := New()
	if err != nil {
		t.Fatal(err)
	}
	r := tpl.(*jqTemplate)

	for _, tc := range test {
		got := r.fixQuery(tc.input)
		if got != tc.want {
			t.Fatalf("got: %s, want: %s\n", got, tc.want)
		}
	}
}

func dataSource() (map[string]any, error) {
	const sample = `
	{
		"firstName": "Charles",
		"lastName": "Doe",
		"age": 41,
		"location": {
		  "city": "San Fracisco",
		  "postalCode": "94103"
		},
		"hobbies": [
		  "chess",
		  "netflix"
		]
	  }`

	res := map[string]any{}
	err := json.Unmarshal([]byte(sample), &res)
	return res, err
}
