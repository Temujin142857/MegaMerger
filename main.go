package main

import (
	"flag"
	"fmt"
	"reflect"
	"sync"
	//"gg"
)

type edgePath struct {
	edges []int
}

type Message struct {
	catagory        string
	sender          int
	level           int
	city            int
	callbackPath    edgePath
	destinationPath edgePath
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
		start(node)
	}

	cases := make([]reflect.SelectCase, len(node.edges))
	for i, ch := range node.edges {
		cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ch)}
	}

	for true {

		chosen, value, ok := reflect.Select(cases)
		if ok {
			panic("unexpected channel closed")
		}
		senderIndex := node.edges[chosen]
		message := value.String()

		switch {
		case message == "" && node.state == "":
			fmt.Println("")
		case message == "a" && node.state == "":

		case message == "complete":

		}
		if node.state == "done" {
			break
		}
	}
	fmt.Println("done")

}

func start(node *Node) {
	node.state = "Downtown"
	path := findSmallestExternalEdge(node)

}

func findSmallestExternalEdge(node *Node) int { return 0 }

func (s *edgePath) Push(item int) {
	s.edges = append(s.edges, item)
}

func (s *edgePath) Pop() (int, bool) {
	if len(s.edges) == 0 {
		return 0, false
	}
	item := s.edges[len(s.edges)-1]
	s.edges = s.edges[:len(s.edges)-1]
	return item, true
}

func (s *edgePath) IsEmpty() bool {
	return len(s.edges) == 0
}
