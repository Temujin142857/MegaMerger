package main

import (
	"fmt"
	"reflect"
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
}

type Vertex struct {
	name    int
	node1   int
	node2   int
	channel chan Message
}

func instructions(node *Node, complexity *int) {
	if node.initiator {
		start(node)
	}

	cases := make([]reflect.SelectCase, len(node.edges))
	for i, ch := range node.edges {
		cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ch)}
	}

	for true {

		chosen, value, ok := reflect.Select(cases)
		if !ok {
			panic("unexpected channel closed")
		}
		senderIndex := node.edges[chosen]
		message := value.Interface().(Message)

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
