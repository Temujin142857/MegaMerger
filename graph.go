package main

import (
	"fmt"
)

var STATES = [...]string{"DOWNTOWN", "ASLEEP", "VILLAGE", "DONE"}
var MESSAGE_CATEGORIES = [...]string{"FIND_SMALLEST_FRINGE_EDGE", "SMALLEST_FRINGE_EDGE_FOUND", "MERE_REQUEST", "MERGE_REQUESTED", "GET_ABSORBED", "WE_ABSORBED_THEM", "CITY_CHECK"}
var SUB_STATES = [...]string{"WAITING_TO_REPLY"}

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
	inbox     chan Message
}

type Vertex struct {
	name  int
	node1 *Node
	node2 *Node
}

func instructions(node *Node, complexity *int) {
	if node.initiator {
		start(node)
	}

	for true {

		message := <-node.inbox
		senderIndex := message.sender

		switch {
		case message.catagory == "" && node.state == "":
			fmt.Println("")
		case message.catagory == "a" && node.state == "":

		case message.catagory == "complete":

		}
		if node.state == "done" {
			break
		}
	}
	fmt.Println("done")
}

func findSmallestExternalEdge(node *Node) int { return 0 }

// start here next time
func sendMessage(node *Node, target int, complexity *int) {
	*complexity++
}

func start(node *Node) {
	node.state = "Downtown"
	path := findSmallestExternalEdge(node)

}

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
