package main

import "math"

var STATES = [...]string{"DOWNTOWN", "ASLEEP", "VILLAGE", "DONE"}
var MESSAGE_CATEGORIES = [...]string{"FIND_SMALLEST_FRINGE_EDGE", "SMALLEST_FRINGE_EDGE_FOUND", "MERE_REQUEST", "MERGE_REQUESTED", "GET_ABSORBED", "WE_ABSORBED_THEM", "CITY_CHECK"}
var SUB_STATES = [...]string{"WAITING_TO_REPLY"}

type EdgePath struct {
	edges []int
}

type Message struct {
	catagory string
	//note this is the edge it came from, not a nodes id
	sender          int
	level           int
	city            int
	callbackPath    EdgePath
	destinationPath EdgePath
	//the next three are optional information, only used for some messages, in retrospect multiple types of message structs would be cleaner, might change if I have time
	//answer will be internal or external
	answer string
	//payload is for smallestFringe edge found, it's the smallest edge found so far
	//payload is the used as a countdown on the way back as part of a fringe
	payload int
	//payload 2 is the id of the sender
	//or in termination broadcast is the leader
	payload2 int
}

type PendingMergeRequest struct {
	sender   int
	level    int
	payload2 int
}

type PendingCityCheck struct {
	sender int
	level  int
	city   int
}

// parent and children are indexes for the edges array
type Node struct {
	name  int
	level int
	city  int
	//parent and children should be edge indixes
	parent int
	edges  map[int]Vertex
	//has different meaning in setup, but after gets converted to map edge id->0 if unkonwn, 1 if internal, 2 if a child, 3 if parent
	neighbors map[int]int
	//maps local vertex to 1 or 0
	nodesIveRequested map[int]int
	//maps the id of the node that sent me the request, to a 0 or 1
	nodesThatHaveRequestedMe    map[int]int
	chidlrenCount               int
	foundMySmallestExternalEdge bool
	smallestExternalEdgeFound   Message
	state                       string
	waitingToReply              bool
	waitingForReply             bool
	searchingForFringEdge       bool
	waitingToReplyToCityCheck   bool
	pendingCityChecks           []PendingCityCheck
	//requests that will be resolved via absorbtion or freindly merge later on
	pendingMergeRequests         []PendingMergeRequest
	fringeEdgeFoundResponceCount int
	initiator                    bool
	inbox                        chan Message
	leader                       int
}

func NewNode(id int, initiatior bool, nodesNum int) Node {
	state := "Asleep"
	if initiatior {
		state = "Downtown"
	}
	n := Node{name: id, level: 1, city: id, edges: make(map[int]Vertex), neighbors: make(map[int]int), chidlrenCount: 0, nodesIveRequested: make(map[int]int),
		nodesThatHaveRequestedMe: make(map[int]int), foundMySmallestExternalEdge: false, smallestExternalEdgeFound: Message{catagory: "smallestFringeEdgeFound", payload: math.MaxInt,
			callbackPath: EdgePath{edges: []int{}}}, state: state, fringeEdgeFoundResponceCount: 0, initiator: initiatior, inbox: make(chan Message, nodesNum)}
	return n
}

type Vertex struct {
	id    int
	node1 *Node
	node2 *Node
}

func (s *EdgePath) Push(item int) {
	s.edges = append(s.edges, item)
}

func (s *EdgePath) Pop() int {
	if len(s.edges) == 0 {
		return -1
	}
	item := s.edges[len(s.edges)-1]
	s.edges = s.edges[:len(s.edges)-1]
	return item
}

func (s *EdgePath) Peek() int {
	if len(s.edges) == 0 {
		return -1
	}
	item := s.edges[len(s.edges)-1]
	return item
}

func (s *EdgePath) IsEmpty() bool {
	return len(s.edges) == 0
}

func removeMR(s []PendingMergeRequest, i int) []PendingMergeRequest {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func removeCC(s []PendingCityCheck, i int) []PendingCityCheck {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func cloneMessage(m Message) Message {
	clone := m

	clone.callbackPath.edges = append([]int(nil), m.callbackPath.edges...)
	clone.destinationPath.edges = append([]int(nil), m.destinationPath.edges...)

	return clone
}
