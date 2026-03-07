package main

import (
	"math/rand/v2"
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

var nodes []Node

func main() {
	var wg sync.WaitGroup

	setup(5, 7, 2)
	for _, node := range nodes {
		wg.Go(func() { instructions(node) })
	}

	wg.Wait()
}

func setup(upperBoundOfNodes int, connectionsNum int, initiatorNum int) {
	for i := 1; i < upperBoundOfNodes+1; i++ {
		node := Node{name: i, level: 0, city: -1, parent: -1, state: "asleep", initiator: false}
		nodes = append(nodes, node)
		if i-1 < initiatorNum {
			nodes[i-1].initiator = true
		}
	}
	for i := 0; i < connectionsNum; i++ {
		n1 := rand.IntN(upperBoundOfNodes)
		for len(nodes[n1].edges)>=upperBoundOfNodes{
			n1 := rand.IntN(upperBoundOfNodes)
		}
		n2 := rand.IntN(upperBoundOfNodes)
		for n2 == n1 || nodes[n1].neighbors[n2]==1{
			n2 = rand.IntN(upperBoundOfNodes)
		}
		connect(i, &nodes[n1], &nodes[n2])
	}
	VisualizeGraph(nodes, "network")
}

func connect(i int, node1 *Node, node2 *Node) {
	v := Vertex{name: i, node1: node1.name, node2: node2.name, channel: make(chan Message)}
	node1.edges = append(node1.edges, v)
	node1.neighbors[]
	node2.edges = append(node2.edges, v)
}

func instructions(node Node) {

}

func transmit() {

}

func drawGraph() {
	//DrawCircle(x, y, r)

}
