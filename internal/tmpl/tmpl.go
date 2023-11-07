package tmpl

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"text/template"
)

type JQTemplate interface {
	Execute(query string, data any) (string, error)
}

var _ JQTemplate = (*jqTemplate)(nil)

func New() (JQTemplate, error) {
	re, err := regexp.Compile(`^jq\s+"([^"]*)"`)
	if err != nil {
		return nil, err
	}

	return &jqTemplate{
		re:  re,
		tpl: template.New("tmp").Funcs(FuncMap()),
	}, nil
}

type jqTemplate struct {
	re  *regexp.Regexp
	tpl *template.Template
}

func (t *jqTemplate) Execute(query string, data any) (string, error) {
	tpl, err := t.tpl.Parse(t.fixQuery(query))
	if err != nil {
		return "", err
	}

	buf := bytes.Buffer{}
	err = tpl.Execute(&buf, data)
	return buf.String(), err
}

func (t *jqTemplate) fixQuery(q string) string {
	if !strings.HasPrefix(q, "{{") {
		return q
	}

	if !strings.HasSuffix(q, "}}") {
		return q
	}

	s := strings.TrimPrefix(q, "{{")
	s = strings.TrimSuffix(s, "}}")
	s = strings.TrimSpace(s)

	if !t.re.MatchString(s) {
		return q
	}

	if strings.HasSuffix(s, ".") {
		return q
	}

	return fmt.Sprintf("{{ %s . }}", s)
}
