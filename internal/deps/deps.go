package deps

import (
	"fmt"
)

type Graph interface {
	DependOn(child, parent string) error
	Leaves() []string
	TopoSorted() []string
}

func New() Graph {
	return &graph{
		dependencies: make(depmap),
		dependents:   make(depmap),
		nodes:        make(nodeset),
	}
}

var _ Graph = (*graph)(nil)

type nodeset map[string]struct{}

type depmap map[string]nodeset

type graph struct {
	nodes        nodeset
	dependencies depmap
	dependents   depmap
}

func (g *graph) DependOn(child, parent string) error {
	if child == parent {
		return fmt.Errorf("self-referential dependencies not allowed")
	}

	if g.dependsOn(parent, child) {
		return fmt.Errorf("circular dependencies not allowed")
	}

	// Add nodes.
	g.nodes[parent] = struct{}{}
	g.nodes[child] = struct{}{}

	// Add edges.
	addNodeToNodeset(g.dependents, parent, child)
	addNodeToNodeset(g.dependencies, child, parent)

	return nil
}

func (g *graph) Leaves() []string {
	leaves := make([]string, 0)

	for node := range g.nodes {
		if _, ok := g.dependencies[node]; !ok {
			leaves = append(leaves, node)
		}
	}

	return leaves
}
func (g *graph) TopoSorted() []string {
	nodeCount := 0
	layers := g.TopoSortedLayers()
	for _, layer := range layers {
		nodeCount += len(layer)
	}

	allNodes := make([]string, 0, nodeCount)
	for _, layer := range layers {
		for _, node := range layer {
			allNodes = append(allNodes, node)
		}
	}

	return allNodes
}

func (g *graph) TopoSortedLayers() [][]string {
	layers := [][]string{}

	shrinkingGraph := g.clone()
	for {
		leaves := shrinkingGraph.Leaves()
		if len(leaves) == 0 {
			break
		}

		layers = append(layers, leaves)
		for _, leafNode := range leaves {
			shrinkingGraph.remove(leafNode)
		}
	}

	return layers
}

func (g *graph) deps(child string) nodeset {
	return g.buildTransitive(child, g.immediateDependencies)
}

func (g *graph) dependsOn(child, parent string) bool {
	deps := g.deps(child)
	_, ok := deps[parent]
	return ok
}

func (g *graph) hasDependent(parent, child string) bool {
	deps := g.Dependents(parent)
	_, ok := deps[child]
	return ok
}

func (g *graph) remove(node string) {
	// Remove edges from things that depend on `node`.
	for dependent := range g.dependents[node] {
		removeFromDepmap(g.dependencies, dependent, node)
	}
	delete(g.dependents, node)

	// Remove all edges from node to the things it depends on.
	for dependency := range g.dependencies[node] {
		removeFromDepmap(g.dependents, dependency, node)
	}
	delete(g.dependencies, node)

	// Finally, remove the node itself.
	delete(g.nodes, node)
}

func (g *graph) immediateDependencies(node string) nodeset {
	return g.dependencies[node]
}

func (g *graph) Dependents(parent string) nodeset {
	return g.buildTransitive(parent, g.immediateDependents)
}

func (g *graph) immediateDependents(node string) nodeset {
	return g.dependents[node]
}

func (g *graph) clone() *graph {
	return &graph{
		dependencies: copyDepmap(g.dependencies),
		dependents:   copyDepmap(g.dependents),
		nodes:        copyNodeset(g.nodes),
	}
}

// buildTransitive starts at `root` and continues calling `nextFn` to keep discovering more nodes until
// the graph cannot produce any more. It returns the set of all discovered nodes.
func (g *graph) buildTransitive(root string, nextFn func(string) nodeset) nodeset {
	if _, ok := g.nodes[root]; !ok {
		return nil
	}

	out := make(nodeset)
	searchNext := []string{root}
	for len(searchNext) > 0 {
		// List of new nodes from this layer of the dependency graph. This is
		// assigned to `searchNext` at the end of the outer "discovery" loop.
		discovered := []string{}
		for _, node := range searchNext {
			// For each node to discover, find the next nodes.
			for nextNode := range nextFn(node) {
				// If we have not seen the node before, add it to the output as well
				// as the list of nodes to traverse in the next iteration.
				if _, ok := out[nextNode]; !ok {
					out[nextNode] = struct{}{}
					discovered = append(discovered, nextNode)
				}
			}
		}
		searchNext = discovered
	}

	return out
}

func removeFromDepmap(dm depmap, key, node string) {
	nodes := dm[key]
	if len(nodes) == 1 {
		// The only element in the nodeset must be `node`, so we
		// can delete the entry entirely.
		delete(dm, key)
	} else {
		// Otherwise, remove the single node from the nodeset.
		delete(nodes, node)
	}
}

func copyNodeset(s nodeset) nodeset {
	out := make(nodeset, len(s))
	for k, v := range s {
		out[k] = v
	}
	return out
}

func copyDepmap(m depmap) depmap {
	out := make(depmap, len(m))
	for k, v := range m {
		out[k] = copyNodeset(v)
	}
	return out
}

func addNodeToNodeset(dm depmap, key, node string) {
	nodes, ok := dm[key]
	if !ok {
		nodes = make(nodeset)
		dm[key] = nodes
	}
	nodes[node] = struct{}{}
}
