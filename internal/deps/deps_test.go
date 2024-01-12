package deps

import (
	"fmt"
	"testing"
)

func TestDeps(t *testing.T) {
	all := []item{
		{
			name:     "apiPipelineRunDetail",
			dependOn: "apiPipeline",
		},
		{
			name: "apiPipeline",
		},
		{
			name:     "api3",
			dependOn: "api4",
		},
		{
			name: "api4",
		},
	}

	g := New()

	for _, el := range all {
		if len(el.dependOn) > 0 {
			_ = g.DependOn(el.name, el.dependOn)
		}
	}

	fmt.Println(g.TopoSorted())
}

type item struct {
	name     string
	dependOn string
}
