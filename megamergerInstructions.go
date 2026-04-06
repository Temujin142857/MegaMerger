package main

import (
	"fmt"
	"math"
)

func instructions(node *Node, complexity *int, leader *int) {
	if node.initiator {
		start(node, complexity)
	}

	for true {
		message := <-node.inbox

		switch node.state {
		case "Downtown":
			downTownInstructions(node, &message, complexity)
		case "Village":
			villageInstructions(node, &message, complexity)
		case "Asleep":
			asleepInstructions(node, &message, complexity)
		}
		if node.state == "done" {
			break
		}
	}
	fmt.Println("done, leader=", node.city)
}

// target is the index of the edge in node.edges
func sendMessage(node *Node, target int, message *Message, complexity *int) {
	*complexity++
	message.sender = target
	node.edges[target].node2.inbox <- *message
}

func broadcast(node *Node, message *Message, complexity *int) {

}

func sendUpSmallestFringeEdgeFound(node *Node, complexity *int) {
	sendMessage(node, node.parent, &node.smallestExternalEdgeFound, complexity)
	node.fringeEdgeFoundResponceCount = 0
	node.foundMySmallestExternalEdge = false
	node.smallestExternalEdgeFound = Message{catagory: "smallestFringeEdgeFound", payload: math.MaxInt}
}

func queryNonChildren(node *Node, message Message, compexity *int) {
	for k, v := range node.neighbors {
		if v == 0 {
			outMessage := Message{catagory: "cityCheck", level: node.level, city: node.city}
			sendMessage(node, k, &outMessage, compexity)
			return
		}
	}
}

func start(node *Node, complexity *int) {
	node.state = "Downtown"
	node.city = node.name
	node.level = 1
	mink := math.MaxInt
	for k, _ := range node.neighbors {
		if k < mink {
			mink = k
		}
	}
	outMessage := Message{catagory: "mergeRequest", level: node.level, city: node.city}
	sendMessage(node, mink, &outMessage, complexity)
}

func asleepInstructions(node *Node, message *Message, complexity *int) {
	switch message.catagory {
	case "findSmallestFringeEdge":
		panic("Asleep node recieve findSmallestFringeEdge")

	case "smallestFringeEdgeFound":
		panic("Asleep node recieve smallestFringeEdgeFound")

	case "mergeRequest":
		outMessage := Message{catagory: "getAbsorbed", level: node.level, city: node.city, destinationPath: message.callbackPath}
		if node.level > message.level {
			sendMessage(node, outMessage.destinationPath.Pop(), &outMessage, complexity)
			node.neighbors[message.sender] = 2
			node.chidlrenCount++
		} else {
			node.substate = "WaitingToReply"
			node.pendingMergeRequests = append(node.pendingMergeRequests, PendingMergeRequest{sender: message.sender, level: message.level})
		}

	case "mergeAccepted":
		panic("Asleep node recieved mergeAccespted")

	case "getAbsorbed":
		panic("Asleep node recieve getAbsorbed")

	case "cityCheck":
		node.state = "Downtown"
		node.city = node.name
		node.level = 1
		outMessage := Message{catagory: "cityCheckReply", level: node.level, city: node.city, answer: "external"}
		sendMessage(node, message.sender, &outMessage, complexity)

	case "cityCheckReply":
		panic("Asleep node recieved cityCheckReply")

	case "terminationBroadcast":
		broadcast(node, message, complexity)
		node.state = "Done"
	}
}

func villageInstructions(node *Node, message *Message, complexity *int) {
	switch message.catagory {
	case "findSmallestFringeEdge":
		if len(node.neighbors) == 1 {
			outMessage := Message{catagory: "smallestFringeEdgeFound", level: node.level, city: node.city, payload: math.MaxInt}
			sendMessage(node, node.parent, &outMessage, complexity)
		}
		if len(node.neighbors)-node.chidlrenCount == 1 {
			node.foundMySmallestExternalEdge = true
		}
		callback := EdgePath{edges: append(message.callbackPath.edges, message.sender)}
		outMessage := Message{catagory: "findSmallestFringeEdge", level: node.level, city: node.city, callbackPath: callback}
		broadcast(node, &outMessage, complexity)
		outMessage2 := Message{catagory: "cityCheck", level: node.level, city: node.city, callbackPath: callback}
		queryNonChildren(node, outMessage2, complexity)

	case "smallestFringeEdgeFound":
		if message.sender == node.parent {
			panic("village recieved smallest fringe edge from parent")
		}
		node.fringeEdgeFoundResponceCount++
		if message.payload < node.smallestExternalEdgeFound.payload {
			node.smallestExternalEdgeFound = *message
			node.smallestExternalEdgeFound.callbackPath.edges = append(node.smallestExternalEdgeFound.callbackPath.edges, message.sender)
		}
		if node.fringeEdgeFoundResponceCount == node.chidlrenCount && node.foundMySmallestExternalEdge {
			sendUpSmallestFringeEdgeFound(node, complexity)
		}

	case "mergeRequest":
		outMessage := Message{catagory: "getAbsorbed", level: node.level, city: node.city, destinationPath: message.callbackPath}
		if node.level > message.level {
			sendMessage(node, outMessage.destinationPath.Pop(), &outMessage, complexity)
			node.neighbors[message.sender] = 2
			node.chidlrenCount++
		} else {
			node.substate = "WaitingToReply"
			node.pendingMergeRequests = append(node.pendingMergeRequests, PendingMergeRequest{sender: message.sender, level: message.level})
		}

	case "mergeAccepted":
		node.level = message.level
		node.city = message.city
		if message.sender == node.parent {
			broadcast(node, message, complexity)
		} else if node.neighbors[message.sender] == 2 {
			outMessage := Message{catagory: "mergeAccepted", level: node.level, city: node.city}
			sendMessage(node, node.parent, &outMessage, complexity)
			node.neighbors[node.parent] = 2
			node.parent = message.sender
			node.neighbors[message.sender] = 3
		} else if node.neighbors[message.sender] == 1 {
			panic("village recieved merge accepted from internal non child/parent node")
		} else {
			//case where it's establishing a friendly merger with the neighbor
		}

	case "getAbsorbed":
		if message.level <= node.level {
			panic("village recieved asborb message with smaller level")
		}
		node.city = message.city
		node.level = message.level
		outMessage := Message{catagory: "getAbsorbed", level: message.level, city: message.city}
		if message.sender == node.parent {
			broadcast(node, &outMessage, complexity)
		} else if node.neighbors[message.sender] == 1 {
			panic("village recieved merge accepted from internal non child/parent node")
		} else {
			//case where it's coming from external neighbor or child
			sendMessage(node, node.parent, &outMessage, complexity)
			node.neighbors[node.parent] = 2
			node.parent = message.sender
			node.neighbors[message.sender] = 3
		}

	case "cityCheck":
		if message.city == node.city {
			outMessage := Message{catagory: "cityCheckReply", level: node.level, city: node.city, answer: "internal"}
			sendMessage(node, outMessage.destinationPath.Pop(), &outMessage, complexity)
		} else if node.level >= message.level {
			outMessage := Message{catagory: "cityCheckReply", level: node.level, city: node.city, answer: "external"}
			sendMessage(node, outMessage.destinationPath.Pop(), &outMessage, complexity)
		} else {
			node.substate = "WatingToReply"
			node.pendingMergeRequests = append(node.pendingMergeRequests, PendingMergeRequest{sender: message.sender, level: message.level})
		}

	case "cityCheckReply":
		if message.answer == "internal" {
			node.neighbors[message.sender] = 1
		} else {
			node.foundMySmallestExternalEdge = true
			if message.sender > node.smallestExternalEdgeFound.payload {
				node.smallestExternalEdgeFound = Message{catagory: "smallestFringeEdgeFound", city: node.city, level: node.level, callbackPath: EdgePath{edges: []int{message.sender}}}
			}
			if node.fringeEdgeFoundResponceCount == node.chidlrenCount {
				sendUpSmallestFringeEdgeFound(node, complexity)
			}
		}

	case "terminationBroadcast":
		broadcast(node, message, complexity)
		node.state = "Done"
	}
}

func downTownInstructions(node *Node, message *Message, complexity *int) {
	switch message.catagory {
	case "findSmallestFringeEdge":
		panic("Downtown recieved find smallest fringe edge")

	case "smallestFringeEdgeFound":
		callback := EdgePath{edges: []int{message.callbackPath.Peek()}}
		outMessage := Message{catagory: "mergeRequest", level: node.level, city: node.city, callbackPath: callback, destinationPath: message.callbackPath}
		sendMessage(node, outMessage.destinationPath.Pop(), &outMessage, complexity)

	case "mergeRequest":
		outMessage := Message{catagory: "getAbsorbed", level: node.level, city: node.city, destinationPath: message.callbackPath}
		if node.level > message.level {
			sendMessage(node, outMessage.destinationPath.Pop(), &outMessage, complexity)
			node.neighbors[message.sender] = 2
			node.chidlrenCount++
		} else {
			node.substate = "WaitingToReply"
			node.pendingMergeRequests = append(node.pendingMergeRequests, PendingMergeRequest{sender: message.sender, level: message.level})
		}

	case "mergeAccepted":
		//need to andle the frienfly merger case
		outMessage := Message{catagory: "mergeAccepted", level: message.level, city: message.city}
		if message.city != node.city {
			node.city = message.city
			node.level = message.level
			node.state = "Village"
			if node.neighbors[message.sender] == 2 {
				node.chidlrenCount--
			}
			node.parent = message.sender
			node.neighbors[message.sender] = 3
		}
		broadcast(node, &outMessage, complexity)

	case "getAbsorbed":
		if message.level <= node.level {
			panic("downtown recieved asborb message with smaller level")
		} else if node.neighbors[message.sender] == 1 {
			panic("downtown recieved asborb message from non child neighbor")
		}
		outMessage := Message{catagory: "getAbsorbed", level: message.level, city: message.city}
		broadcast(node, &outMessage, complexity)
		node.city = message.city
		node.level = message.level
		node.state = "Village"
		if node.neighbors[message.sender] == 2 {
			node.chidlrenCount--
		}
		node.parent = message.sender
		node.neighbors[message.sender] = 3

	case "cityCheck":
		if message.city == node.city {
			outMessage := Message{catagory: "cityCheckReply", level: node.level, city: node.city, answer: "internal"}
			sendMessage(node, outMessage.destinationPath.Pop(), &outMessage, complexity)
		} else if node.level >= message.level {
			outMessage := Message{catagory: "cityCheckReply", level: node.level, city: node.city, answer: "external"}
			sendMessage(node, outMessage.destinationPath.Pop(), &outMessage, complexity)
		} else {
			node.substate = "WatingToReply"
			node.pendingMergeRequests = append(node.pendingMergeRequests, PendingMergeRequest{sender: message.sender, level: message.level})
		}

	case "cityCheckReply":
		//need to write

	case "terminationBroadcast":
		panic("downtown recieved termination broadcast")

	}
}
