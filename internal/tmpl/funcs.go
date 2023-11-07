package tmpl

import (
	"strings"
	"text/template"
	"time"

	"github.com/itchyny/gojq"
)

var genericMap = map[string]interface{}{
	"now": func() string {
		return time.Now().UTC().
			Format("2006-01-02T15:04:05Z")
	},

	"empty": func(s string) bool {
		return len(strings.TrimSpace(s)) == 0
	},

	"jq": jq,

	// "printf": func(format string, args ...any) string {
	// 	return fmt.Sprintf(format, args...)
	// },

	// "noescape": func(str string) template.HTML {
	// 	return template.HTML(str)
	// },
}

// FuncMap returns a copy of the basic function map as a map[string]interface{}.
func FuncMap() map[string]any {
	gfm := make(map[string]any, len(genericMap))
	for k, v := range genericMap {
		gfm[k] = v
	}
	return gfm
}

// funcMap returns a 'text/template'.FuncMap
func funcMap() template.FuncMap {
	return template.FuncMap(FuncMap())
}

func jq(q string, input map[string]any) (string, error) {
	enc := newEncoder(false, 0)

	query, err := gojq.Parse(q)
	if err != nil {
		return "", err
	}

	iter := query.Run(input) // or query.RunWithContext
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
