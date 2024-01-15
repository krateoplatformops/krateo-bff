package tmpl

import (
	"encoding/json"
	"fmt"
	"regexp"
	"testing"
)

func TestRegexPatternBuild(t *testing.T) {
	leftDelim, rightDelim := "${", "}"
	pattern := fmt.Sprintf("^%s\\s+(.*)%s",
		regexp.QuoteMeta(leftDelim),
		regexp.QuoteMeta(rightDelim))

	fmt.Println(pattern)
}

func TestJQTemplate(t *testing.T) {
	test := []struct {
		input string
		want  string
	}{
		{`${ .age }`, "41"},
		{` .age }}`, ` .age }}`},
		{`${ .location.city }`, "San Fracisco"},
		{"hello world", "hello world"},
		{`${ .hobbies | join(",") }`, "chess,netflix"},
		{`${ .id }`, "1"},
		{`${ "/todos/" + (.id|tostring) +  "/comments" }`, "/todos/1/comments"},
	}

	ds, err := dataSource()
	if err != nil {
		t.Fatal(err)
	}

	tpl, err := New("${", "}")
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range test {
		got, err := tpl.Execute(tc.input, ds)
		if err != nil {
			t.Fatal(err)
		}

		if got != tc.want {
			t.Fatalf("got: %s, want: %s\n", got, tc.want)
		}
	}
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

func TestAcceptQuery(t *testing.T) {
	test := []struct {
		input string
		want  string
		ok    bool
	}{
		{`${ .age }`, `.age`, true},
		{` .age }}`, ` .age }}`, false},
		{`${ .location.city }`, `.location.city`, true},
		{`hello world`, `hello world`, false},
		{`${ .hobbies | join(",") }`, `.hobbies | join(",")`, true},
	}

	tpl, err := New("${", "}")
	if err != nil {
		t.Fatal(err)
	}
	r := tpl.(*jqTemplate)

	for _, tc := range test {
		got, ok := r.acceptQuery(tc.input)
		if got != tc.want {
			t.Fatalf("got: %s, want: %s\n", got, tc.want)
		}
		if ok != tc.ok {
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
		],
		"id": 1
	  }`

	res := map[string]any{}
	err := json.Unmarshal([]byte(sample), &res)
	return res, err
}
