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
		//fmt.Println("Node", node.name, "waiting...")
		message := <-node.inbox
		//fmt.Println("Node", node.name, "received", message.catagory)

		switch node.state {
		case "Downtown":
			downTownInstructions(node, &message, complexity)
		case "Village":
			villageInstructions(node, &message, complexity)
		case "Asleep":
			fmt.Println("no sleep", node.name)
			//asleepInstructions(node, &message, complexity)
		}
		if node.state == "Done" {
			break
		}
	}
	*leader = node.city
	fmt.Println("done, leader=", node.city)
}

// target is the index of the edge in node.edges
func sendMessage(node *Node, target int, message *Message, complexity *int) {
	*complexity++
	m := cloneMessage(*message)
	m.sender = target
	//fmt.Println(target, node.edges)
	if node.edges[target].node1.name == node.name {
		node.edges[target].node2.inbox <- m
		fmt.Println(node.name, node.edges[target].node2.name, m)
	} else {
		node.edges[target].node1.inbox <- m
		fmt.Println(node.name, node.edges[target].node1.name, m)
	}
}

func broadcast(node *Node, message *Message, complexity *int) {
	for k, v := range node.neighbors {
		if v == 2 {
			sendMessage(node, k, message, complexity)
		}
	}
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
	mink := math.MaxInt
	for k, _ := range node.neighbors {
		if k < mink {
			mink = k
		}
	}
	node.substate = "WaitingForReply"
	node.nodesIveRequested[mink] = 1
	outMessage := Message{catagory: "mergeRequest", level: node.level, city: node.city, payload2: node.name}
	sendMessage(node, mink, &outMessage, complexity)
}

func nodeHasChangedLevel(node *Node, complexity *int) {
	if node.substate == "WaitingToReply" {
		for i := 0; i < len(node.pendingMergeRequests); i++ {
			if node.pendingMergeRequests[i].level < node.level {
				outMessage := Message{catagory: "getAbsorbed", level: node.level, city: node.city}
				sendMessage(node, node.pendingMergeRequests[i].sender, &outMessage, complexity)
				node.pendingMergeRequests = remove(node.pendingMergeRequests, i)
				fmt.Println("levelc:", len(node.pendingMergeRequests), node.name)
				if len(node.pendingMergeRequests) == 0 {
					node.substate = ""
					return
				}
				//since we're removing an element, we back up to get the new element in this spot
				i--
			}
		}
		fmt.Println("exited")
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
		outMessage := Message{catagory: "findSmallestFringeEdge", level: node.level, city: node.city, payload: math.MaxInt}
		broadcast(node, &outMessage, complexity)
		outMessage2 := Message{catagory: "cityCheck", level: node.level, city: node.city}
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
		if message.payload > 0 {
			target := message.destinationPath.Pop()
			if message.payload == 1 {
				node.substate = "WaitingForReply"
				node.nodesIveRequested[target] = 1
			}
			message.payload--
			message.payload2 = node.name
			sendMessage(node, target, message, complexity)
		} else if node.level > message.level {
			outMessage := Message{catagory: "getAbsorbed", level: node.level, city: node.city}
			sendMessage(node, message.sender, &outMessage, complexity)
			node.neighbors[message.sender] = 2
			node.chidlrenCount++
		} else if node.substate == "WaitingForReply" && node.nodesIveRequested[message.sender] == 1 {
			oldParent := node.parent
			node.neighbors[oldParent] = 2
			node.chidlrenCount++
			if node.name < message.payload2 {
				node.state = "Downtown"
				node.neighbors[message.sender] = 2
				node.chidlrenCount++
			} else {
				node.parent = message.sender
				node.neighbors[message.sender] = 3
			}
			node.city = message.sender
			node.level++
			outMessage := Message{catagory: "mergeAccepted", level: node.level, city: node.city}
			sendMessage(node, oldParent, &outMessage, complexity)
			nodeHasChangedLevel(node, complexity)
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
		nodeHasChangedLevel(node, complexity)

	case "cityCheck":
		if message.city == node.city {
			outMessage := Message{catagory: "cityCheckReply", level: node.level, city: node.city, answer: "internal"}
			sendMessage(node, message.sender, &outMessage, complexity)
		} else if node.level >= message.level {
			outMessage := Message{catagory: "cityCheckReply", level: node.level, city: node.city, answer: "external"}
			sendMessage(node, message.sender, &outMessage, complexity)
		} else {
			node.substate = "WaitingToReply"
			node.pendingMergeRequests = append(node.pendingMergeRequests, PendingMergeRequest{sender: message.sender, level: message.level})
		}

	case "cityCheckReply":
		fmt.Println("here4", node.name, node.fringeEdgeFoundResponceCount, node.chidlrenCount, node.foundMySmallestExternalEdge)
		if message.answer == "internal" {
			node.neighbors[message.sender] = 1
		} else {
			node.foundMySmallestExternalEdge = true
			if message.sender < node.smallestExternalEdgeFound.payload {
				node.smallestExternalEdgeFound = *message
				node.smallestExternalEdgeFound.payload = message.sender
			}
		}
		fmt.Println("here4", node.name, node.fringeEdgeFoundResponceCount, node.chidlrenCount, node.foundMySmallestExternalEdge)
		if node.fringeEdgeFoundResponceCount == node.chidlrenCount && node.foundMySmallestExternalEdge {
			sendUpSmallestFringeEdgeFound(node, complexity)
		}

	case "terminationBroadcast":
		broadcast(node, message, complexity)
		node.state = "Done"
	}
}

func downTownInstructions(node *Node, message *Message, complexity *int) {
	switch message.catagory {
	case "findSmallestFringeEdge":
		panic(fmt.Sprintf("Downtown recieved find smallest fringe edge %d", node.name))

	case "smallestFringeEdgeFound":
		node.fringeEdgeFoundResponceCount++
		if message.payload < node.smallestExternalEdgeFound.payload {
			node.smallestExternalEdgeFound = *message
		}
		if node.fringeEdgeFoundResponceCount == node.chidlrenCount && node.foundMySmallestExternalEdge {
			if node.smallestExternalEdgeFound.payload == math.MaxInt {
				outMessage := Message{catagory: "terminationBroadcast"}
				broadcast(node, &outMessage, complexity)
				node.state = "Done"
				break
			}
			outMessage := Message{catagory: "mergeRequest", level: node.level, city: node.city, destinationPath: node.smallestExternalEdgeFound.callbackPath, payload: len(node.smallestExternalEdgeFound.callbackPath.edges), payload2: node.name}
			if outMessage.payload == 0 {
				node.substate = "WaitingForReply"
				node.nodesIveRequested[node.smallestExternalEdgeFound.sender] = 1
			}
			sendMessage(node, node.smallestExternalEdgeFound.sender, &outMessage, complexity)
		}

	case "mergeRequest":
		outMessage := Message{catagory: "getAbsorbed", level: node.level, city: node.city}
		if node.level > message.level {
			//fmt.Println("here me2", node.name, "target", message.sender)
			sendMessage(node, message.sender, &outMessage, complexity)
			node.neighbors[message.sender] = 2
			node.chidlrenCount++
		} else if node.substate == "WaitingForReply" && node.nodesIveRequested[message.sender] == 1 {
			fmt.Println("here me", node.name, "target", message.sender)
			if node.name < message.payload2 {
				node.state = "Downtown"
				node.neighbors[message.sender] = 2
				node.chidlrenCount++
			} else {
				node.parent = message.sender
				node.neighbors[message.sender] = 3
				node.state = "Village"
			}
			node.city = message.sender
			node.level++
			outMessage := Message{catagory: "mergeAccepted", level: node.level, city: node.city}
			broadcast(node, &outMessage, complexity)
			nodeHasChangedLevel(node, complexity)
			if node.state == "Downtown" {
				outMessage2 := Message{catagory: "findSmallestFringeEdge", city: node.city, level: node.level, payload: math.MaxInt}
				broadcast(node, &outMessage2, complexity)
			}

		} else {
			//fmt.Println("here me2", node.name, "target", message.sender)
			node.substate = "WaitingToReply"
			node.pendingMergeRequests = append(node.pendingMergeRequests, PendingMergeRequest{sender: message.sender, level: message.level})
		}

	case "mergeAccepted":
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
		nodeHasChangedLevel(node, complexity)

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
		nodeHasChangedLevel(node, complexity)

	case "cityCheck":
		if message.city == node.city {
			outMessage := Message{catagory: "cityCheckReply", level: node.level, city: node.city, answer: "internal"}
			sendMessage(node, message.sender, &outMessage, complexity)
		} else if node.level >= message.level {
			outMessage := Message{catagory: "cityCheckReply", level: node.level, city: node.city, answer: "external"}
			sendMessage(node, message.sender, &outMessage, complexity)
		} else {
			node.substate = "WaitingToReply"
			node.pendingMergeRequests = append(node.pendingMergeRequests, PendingMergeRequest{sender: message.sender, level: message.level})
		}

	case "cityCheckReply":
		if message.answer == "internal" {
			node.neighbors[message.sender] = 1
		} else {
			node.foundMySmallestExternalEdge = true
			if message.sender < node.smallestExternalEdgeFound.payload {
				node.smallestExternalEdgeFound = *message
				node.smallestExternalEdgeFound.payload = message.sender
			}
		}
		fmt.Println("here4", node.name, node.fringeEdgeFoundResponceCount, node.chidlrenCount, node.foundMySmallestExternalEdge)
		if node.fringeEdgeFoundResponceCount == node.chidlrenCount && node.foundMySmallestExternalEdge {
			outMessage := Message{catagory: "mergeRequest", level: node.level, city: node.city, destinationPath: node.smallestExternalEdgeFound.callbackPath, payload: len(node.smallestExternalEdgeFound.callbackPath.edges), payload2: node.name}
			if outMessage.payload == 0 {
				node.substate = "WaitingForReply"
				node.nodesIveRequested[node.smallestExternalEdgeFound.sender] = 1
			}
			fmt.Println("here3", node.name)
			sendMessage(node, node.smallestExternalEdgeFound.sender, &outMessage, complexity)
			node.foundMySmallestExternalEdge = false
			node.fringeEdgeFoundResponceCount = 0
		}

	case "terminationBroadcast":
		panic("downtown recieved termination broadcast")

	}
}
