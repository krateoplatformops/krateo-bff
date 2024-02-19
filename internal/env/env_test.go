package env

import (
	"os"
	"testing"
)

func TestServicePort(t *testing.T) {
	table := []struct {
		key  string
		val  string
		def  int
		want int
	}{
		{
			key:  "KRATEO_BFF_PORT",
			val:  "tcp://10.96.234.180:8081",
			def:  8888,
			want: 8081,
		},
		{
			key:  "KRATEO_BFF_PORT",
			val:  "tcp://8080",
			def:  8888,
			want: 8888,
		},
		{
			key:  "KRATEO_BFF_PORT",
			val:  "8080",
			def:  8888,
			want: 8080,
		},
	}

	for i, tc := range table {
		os.Setenv(tc.key, tc.val)
		got := ServicePort(tc.key, tc.def)
		if got != tc.want {
			t.Fatalf("[tc: %d] got: %d - expected: %d", i, got, tc.want)
		}
	}
}
