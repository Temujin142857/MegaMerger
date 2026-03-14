package main

import (
	"flag"
	"fmt"
	"sync"
	//"gg"
)

type Message struct {
	catagory string
	sender   int
	level    int
	city     int
}

type Node struct {
	name      int
	level     int
	city      int
	parent    int
	children  []int
	edges     []Vertex
	neighbors map[int]int
	state     string
	initiator bool
}

type Vertex struct {
	name    int
	node1   int
	node2   int
	channel chan Message
}

var nodes map[int]Node = make(map[int]Node)

var (
	fileFlag     bool
	filePath     string
	withWeight   bool
	initiatorNum int
)

func init() {
	flag.BoolVar(&fileFlag, "fileFlag", false, "Wether making graph from file or rand")
	flag.StringVar(&filePath, "filePath", "default", "name of file with graph")
	flag.BoolVar(&withWeight, "withWeight", false, "Wether the file has weights")
	flag.IntVar(&initiatorNum, "initiatorNum", 1, "How many itiators")
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Error:\n", r)
		}
	}()

	flag.Parse()
	fmt.Println("fileFlag:", fileFlag, ", filePath:", filePath, ", withWeight:", withWeight)

	var wg sync.WaitGroup

	if fileFlag {
		fileSetup(filePath, withWeight, initiatorNum)
	} else {
		randomSetup(5, 7, 2)
	}
	fmt.Println(nodes)
	VisualizeGraph(nodes, "network")

	for _, node := range nodes {
		wg.Go(func() { instructions(&node) })
	}

	wg.Wait()
}

func instructions(node *Node) {
	if node.initiator {

	}

	for true {
		break
		fmt.Println("wtf")
	}
	fmt.Println("done")

}

func transmit() {

}
