package main

import (
	"fmt"
	"os"
)

func VisualizeGraph(nodes map[int]Node, name string) error {

	dotFile := name + ".dot"

	err := writeDOT(nodes, dotFile)
	if err != nil {
		return err
	}

	return nil
}

func writeDOT(nodes map[int]Node, filename string) error {

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	fmt.Fprintln(f, "graph G {")
	fmt.Fprintln(f, "  node [shape=circle];")

	// Nodes
	for _, n := range nodes {

		shape := "circle"
		if n.initiator {
			shape = "doublecircle"
		}

		label := fmt.Sprintf("%d", n.name)

		fmt.Fprintf(
			f,
			"  %d [shape=%s label=\"%s\"];\n",
			n.name,
			shape,
			label,
		)
	}

	// Edges
	seen := map[int]bool{}

	for _, n := range nodes {
		for _, v := range n.edges {

			if seen[v.name] {
				continue
			}
			seen[v.name] = true

			fmt.Fprintf(
				f,
				"  %d -- %d [label=\"v%d\"];\n",
				v.node1,
				v.node2,
				v.name,
			)
		}
	}

	fmt.Fprintln(f, "}")
	return nil
}
