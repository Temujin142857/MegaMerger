package main

import (
	"fmt"
	"sync"
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
	edges     []chan Message
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
	fmt.Printf("hello world")

	var wg sync.WaitGroup

	setup(5, 7)
	for _, node := range nodes {
		wg.Go(func() { instructions(node) })
	}

	wg.Wait()
}

func setup(upperBoundOfNodes int, averageConnections int) {
	//populate nodes
	//connect nodes
	//set some to initiators
}

func connect(node1 Node, node2 Node) {

}

func instructions(node Node) {

}

func transmit() {

}
