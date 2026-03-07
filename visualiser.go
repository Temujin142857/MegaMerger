package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

func VisualizeGraph(nodes []Node, name string) error {

	dotFile := name + ".dot"
	pngFile := name + ".png"

	err := writeDOT(nodes, dotFile)
	if err != nil {
		return err
	}

	// Run Graphviz
	cmd := exec.Command("neato", "-Tpng", dotFile, "-o", pngFile)
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("graphviz failed: %w", err)
	}

	// Open image automatically
	openImage(pngFile)

	return nil
}

func writeDOT(nodes []Node, filename string) error {

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

		label := fmt.Sprintf("%d\\nL%d\\n%s", n.name, n.level, n.state)

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

func openImage(file string) {

	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("xdg-open", file)
	case "darwin":
		cmd = exec.Command("open", file)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", file)
	default:
		return
	}

	cmd.Start()
}
