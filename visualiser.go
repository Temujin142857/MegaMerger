package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

var dotFolderPath = filepath.Join("output", "dotFolder")
var pngFolderPath = filepath.Join("output", "pngFolder")

func VisualizeGraph(nodes map[int]Node, name string) error {

	dotFile := filepath.Join(dotFolderPath, name+".dot")
	pngFile := filepath.Join(pngFolderPath, name+".png")

	err := writeDOT(nodes, dotFile)
	if err != nil {
		return err
	}

	writePNG(dotFile, pngFile)

	return nil
}

func setupOutputFolder() {
	err := os.MkdirAll("output/dotFolder", 0755)
	if err != nil {
		panic(err)
	}
	err = os.MkdirAll("output/pngFolder", 0755)
	if err != nil {
		panic(err)
	}
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

			if seen[v.id] {
				continue
			}
			seen[v.id] = true

			fmt.Fprintf(
				f,
				"  %d -- %d [label=\"v%d\"];\n",
				v.node1.name,
				v.node2.name,
				v.id,
			)
		}
	}

	fmt.Fprintln(f, "}")
	return nil
}

func writePNG(dotFilePath string, pngFilePath string) {
	cmd := exec.Command("dot", "-Tpng", "-Kdot", "-o", pngFilePath, dotFilePath)
	if err := cmd.Run(); err != nil {
		panic(err)
	}
}
