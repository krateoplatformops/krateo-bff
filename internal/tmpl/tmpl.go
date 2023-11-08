package tmpl

import (
	"regexp"
	"strings"

	"github.com/itchyny/gojq"
)

type JQTemplate interface {
	Execute(query string, data any) (string, error)
}

var _ JQTemplate = (*jqTemplate)(nil)

func New() (JQTemplate, error) {
	re, err := regexp.Compile(`^\{\{\s+(.*)\}\}$`)
	if err != nil {
		return nil, err
	}

	return &jqTemplate{re: re}, nil
}

type jqTemplate struct {
	re *regexp.Regexp
}

func (t *jqTemplate) Execute(q string, data any) (string, error) {
	enc := newEncoder(false, 0)

	query, err := gojq.Parse(t.fixQuery(q))
	if err != nil {
		return "", err
	}

	iter := query.Run(data) // or query.RunWithContext
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		if err, ok := v.(error); ok {
			return "", err
		}
		if err := enc.encode(v); err != nil {
			return "", err
		}
	}

	return enc.w.String(), nil
}

func (t *jqTemplate) fixQuery(q string) string {
	if !t.re.MatchString(q) {
		return q
	}

	res := t.re.FindAllStringSubmatch(q, -1)
	if len(res) == 0 || len(res[0]) == 0 {
		return q
	}

	return strings.TrimSpace(res[0][1])
}
