package main

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
	answer          string
	payload         int
}

// parent and children are indexes for the edges array
type Node struct {
	name  int
	level int
	city  int
	//parent and children should be edge indixes
	parent   int
	children []int
	edges    []Vertex
	//only used in setup
	neighbors map[int]int
	//maps vertex to a 1 if internal, 2 if a child
	internalNeighbors      map[int]int
	mySmallestInternalEdge int
	state                  string
	substate               string
	initiator              bool
	inbox                  chan Message
}

type Vertex struct {
	name  int
	node1 *Node
	node2 *Node
}

func (s *edgePath) Push(item int) {
	s.edges = append(s.edges, item)
}

func (s *edgePath) Pop() int {
	if len(s.edges) == 0 {
		return -1
	}
	item := s.edges[len(s.edges)-1]
	s.edges = s.edges[:len(s.edges)-1]
	return item
}

func (s *edgePath) Peek() int {
	if len(s.edges) == 0 {
		return -1
	}
	item := s.edges[len(s.edges)-1]
	return item
}

func (s *edgePath) IsEmpty() bool {
	return len(s.edges) == 0
}

func (n *Node) removeChild(target int) {
	for i, v := range n.children {
		if v == target {
			n.children = append(n.children[:i], n.children[i+1:]...)
			return
		}
	}
}
