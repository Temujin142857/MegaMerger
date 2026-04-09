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
		//fmt.Println("Node", node.name, "received", message.catagory, "from v:", message.sender)

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
		//fmt.Println(node.name, node.edges[target].node2.name, m)
	} else {
		node.edges[target].node1.inbox <- m
		//fmt.Println(node.name, node.edges[target].node1.name, m)
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
	node.smallestExternalEdgeFound = Message{catagory: "smallestFringeEdgeFound", payload: math.MaxInt, callbackPath: EdgePath{}}
}

func queryNonChildren(node *Node, compexity *int) {
	for k, v := range node.neighbors {
		if v == 0 {
			outMessage := Message{catagory: "cityCheck", level: node.level, city: node.city}
			sendMessage(node, k, &outMessage, compexity)
			return
		}
	}
	//println("node has no children but called queryNonChildren", node.name)
	node.foundMySmallestExternalEdge = true
}

func start(node *Node, complexity *int) {
	mink := math.MaxInt
	for k, _ := range node.neighbors {
		if k < mink {
			mink = k
		}
	}
	node.waitingForReply = true
	node.nodesIveRequested[mink] = 1
	outMessage := Message{catagory: "mergeRequest", level: node.level, city: node.city, payload2: node.name}
	sendMessage(node, mink, &outMessage, complexity)
}

func friendlyMergeVillage(node *Node, sender int, payload2 int, complexity *int) {
	//friendly merge, since this is a request from the node I've already requests
	oldParent := node.parent
	node.neighbors[oldParent] = 2
	node.chidlrenCount++
	if node.name < payload2 {
		node.state = "Downtown"
		node.neighbors[sender] = 2
		node.chidlrenCount++
	} else {
		node.parent = sender
		node.neighbors[sender] = 3
	}
	node.city = sender
	node.level++
	outMessage := Message{catagory: "mergeAccepted", level: node.level, city: node.city}
	sendMessage(node, oldParent, &outMessage, complexity)
	node.waitingForReply = false
	nodeHasChangedLevel(node, complexity)
}

func friendlyMergeDowntown(node *Node, sender int, payload2 int, complexity *int) {
	//here we broadcast first so only the old children recieve it, not the possible new children from the merge
	node.city = sender
	node.level++
	outMessage := Message{catagory: "mergeAccepted", level: node.level, city: node.city}
	broadcast(node, &outMessage, complexity)
	if node.name < payload2 {
		node.state = "Downtown"
		node.neighbors[sender] = 2
		node.chidlrenCount++
	} else {
		node.parent = sender
		node.neighbors[sender] = 3
		node.state = "Village"
	}
	nodeHasChangedLevel(node, complexity)
	if node.state == "Downtown" {
		node.searchingForFringEdge = true
		outMessage2 := Message{catagory: "findSmallestFringeEdge", city: node.city, level: node.level, payload: math.MaxInt}
		broadcast(node, &outMessage2, complexity)
		queryNonChildren(node, complexity)
	}
}

func cityCheckLogic(node *Node, sender int, level int, city int, complexity *int, setWait bool) bool {
	if city == node.city {
		outMessage := Message{catagory: "cityCheckReply", level: node.level, city: node.city, answer: "internal"}
		sendMessage(node, sender, &outMessage, complexity)
		return true
	} else if node.level >= level {
		outMessage := Message{catagory: "cityCheckReply", level: node.level, city: node.city, answer: "external"}
		sendMessage(node, sender, &outMessage, complexity)
		return true
	} else if setWait {
		node.waitingToReplyToCityCheck = true
		node.pendingCityChecks = append(node.pendingCityChecks, PendingCityCheck{sender: sender, level: level, city: city})
	}
	return false
}

func sendGetAbsorbedToExternal(node *Node, target int, complexity *int) {
	outMessage := Message{catagory: "getAbsorbed", level: node.level, city: node.city}
	sendMessage(node, target, &outMessage, complexity)
	if node.searchingForFringEdge == true {
		outMessage2 := Message{catagory: "findSmallestFringeEdge", level: node.level, city: node.city, payload: math.MaxInt, callbackPath: EdgePath{edges: []int{}}}
		sendMessage(node, target, &outMessage2, complexity)
	}

	node.neighbors[target] = 2
	node.chidlrenCount++
}

func nodeHasChangedLevel(node *Node, complexity *int) {
	if node.waitingToReply == true {
		//loop through beinding absorbtions
		for i := 0; i < len(node.pendingMergeRequests); i++ {
			if node.pendingMergeRequests[i].friendly == false && node.pendingMergeRequests[i].level < node.level {
				sendGetAbsorbedToExternal(node, node.pendingMergeRequests[i].sender, complexity)
				node.pendingMergeRequests = removeMR(node.pendingMergeRequests, i)
				if len(node.pendingMergeRequests) == 0 {
					node.waitingToReply = false
					return
				}
				//since we're removing an element, we back up to get the new element in this spot
				i--
			} else if node.pendingMergeRequests[i].friendly == true && node.pendingMergeRequests[i].level == node.level {
				if node.state == "Downtown" {
					friendlyMergeDowntown(node, node.pendingMergeRequests[i].sender, node.pendingMergeRequests[i].payload2, complexity)
				} else if node.state == "Village" {
					friendlyMergeVillage(node, node.pendingMergeRequests[i].sender, node.pendingMergeRequests[i].payload2, complexity)
				}
				node.pendingMergeRequests = removeMR(node.pendingMergeRequests, i)
				if len(node.pendingMergeRequests) == 0 {
					node.waitingToFriendlyMerge = false
					if len(node.pendingMergeRequests) == 0 {
						node.waitingToReply = false
					}
					return
				}
				//since we're removing an element, we back up to get the new element in this spot
				i--
			}
		}
	}
	if node.waitingToReplyToCityCheck == true {
		for i := 0; i < len(node.pendingCityChecks); i++ {
			if node.pendingCityChecks[i].level == node.level {
				didSend := cityCheckLogic(node, node.pendingCityChecks[i].sender, node.pendingCityChecks[i].level, node.pendingCityChecks[i].city, complexity, false)
				if didSend {
					node.pendingCityChecks = removeCC(node.pendingCityChecks, i)
					if len(node.pendingCityChecks) == 0 {
						node.waitingToReplyToCityCheck = false
						return
					}
					i--
				}
			} else if node.pendingCityChecks[i].level > node.level {
				panic(fmt.Sprint("city check, I was waiting for equality, but somehow I'm bigger than it now? me, their edge", node.name, node.pendingCityChecks[i].sender))
			}
		}
	}
}

func villageInstructions(node *Node, message *Message, complexity *int) {
	switch message.catagory {
	case "findSmallestFringeEdge":
		//fmt.Println("find edge", node.name)
		node.searchingForFringEdge = true
		if len(node.neighbors) == 1 {
			outMessage := Message{catagory: "smallestFringeEdgeFound", level: node.level, city: node.city, callbackPath: EdgePath{edges: []int{}}, payload: math.MaxInt}
			sendMessage(node, node.parent, &outMessage, complexity)
		}
		if len(node.neighbors)-node.chidlrenCount == 1 {
			//fmt.Println("node has no external neighbors", node.name)
			node.foundMySmallestExternalEdge = true
		}
		outMessage := Message{catagory: "findSmallestFringeEdge", level: node.level, city: node.city, payload: math.MaxInt, callbackPath: EdgePath{edges: []int{}}}
		broadcast(node, &outMessage, complexity)
		if !node.foundMySmallestExternalEdge {
			queryNonChildren(node, complexity)
		}

	case "smallestFringeEdgeFound":
		if message.sender == node.parent {
			panic("village recieved smallest fringe edge from parent")
		}
		node.fringeEdgeFoundResponceCount++
		if message.payload <= node.smallestExternalEdgeFound.payload {
			node.smallestExternalEdgeFound = *message
			node.smallestExternalEdgeFound.callbackPath.edges = append(node.smallestExternalEdgeFound.callbackPath.edges, message.sender)
		}
		if node.fringeEdgeFoundResponceCount == node.chidlrenCount && node.foundMySmallestExternalEdge {
			sendUpSmallestFringeEdgeFound(node, complexity)
			//no longer searching
			node.searchingForFringEdge = false
		}

	case "mergeRequest":
		if message.payload > 0 {
			//case where the message needs to be passed on until it reaches the external node
			target := message.destinationPath.Pop()
			if message.payload == 1 {
				node.waitingForReply = true
				node.nodesIveRequested[target] = 1
			}
			message.payload--
			message.payload2 = node.name
			sendMessage(node, target, message, complexity)
		} else if node.level > message.level {
			//case where this node absorbes the requester
			sendGetAbsorbedToExternal(node, message.sender, complexity)
		} else if node.waitingForReply == true && node.nodesIveRequested[message.sender] == 1 && node.level == message.level {
			friendlyMergeVillage(node, message.sender, message.payload2, complexity)

		} else {
			node.waitingToReply = true
			node.pendingMergeRequests = append(node.pendingMergeRequests, PendingMergeRequest{sender: message.sender, level: message.level, payload2: message.payload2})
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
			//this would be for if they recieved it from an external node
			//but that would be a freindly merger, which I handle when they exchange merge requests, without needing an accept message
		}
		nodeHasChangedLevel(node, complexity)

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
			if node.waitingForReply == true && node.nodesIveRequested[message.sender] == 1 {
				node.waitingForReply = false
			}
			if node.neighbors[message.sender] == 0 {
				node.chidlrenCount++
			}
			node.neighbors[message.sender] = 3
		}
		nodeHasChangedLevel(node, complexity)

	case "cityCheck":
		cityCheckLogic(node, message.sender, message.level, message.city, complexity, true)

	case "cityCheckReply":
		if message.answer == "internal" {
			node.neighbors[message.sender] = 1
			allVisited := true
			for _, v := range node.neighbors {
				if v == 0 {
					allVisited = false
					queryNonChildren(node, complexity)
				}
			}
			if allVisited {
				node.foundMySmallestExternalEdge = true
			}
		} else {
			node.foundMySmallestExternalEdge = true
			if message.sender < node.smallestExternalEdgeFound.payload {
				node.smallestExternalEdgeFound = *message
				node.smallestExternalEdgeFound.payload = message.sender
			}
		}
		if node.fringeEdgeFoundResponceCount == node.chidlrenCount && node.foundMySmallestExternalEdge {
			sendUpSmallestFringeEdgeFound(node, complexity)
			//no longer searching
			node.searchingForFringEdge = false
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
			//no longer searching
			node.searchingForFringEdge = false
			outMessage := Message{catagory: "mergeRequest", level: node.level, city: node.city, destinationPath: node.smallestExternalEdgeFound.callbackPath, payload: len(node.smallestExternalEdgeFound.callbackPath.edges), payload2: node.name}
			if outMessage.payload == 0 {
				node.waitingForReply = true
				node.nodesIveRequested[node.smallestExternalEdgeFound.sender] = 1
			}
			sendMessage(node, node.smallestExternalEdgeFound.sender, &outMessage, complexity)
		}

	case "mergeRequest":
		if node.level > message.level {
			//absorbtion case
			sendGetAbsorbedToExternal(node, message.sender, complexity)
		} else if node.waitingForReply == true && node.nodesIveRequested[message.sender] == 1 && node.level == message.level {
			//friendly merge, since this is a request from the node I've already requests
			friendlyMergeDowntown(node, message.sender, message.payload2, complexity)

		} else {
			node.waitingToReply = true
			node.pendingMergeRequests = append(node.pendingMergeRequests, PendingMergeRequest{sender: message.sender, level: message.level, payload2: message.payload2})
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
		if node.waitingForReply == true && node.nodesIveRequested[message.sender] == 1 {
			node.waitingForReply = false
		}
		if node.neighbors[message.sender] == 2 {
			node.chidlrenCount--
		}
		node.parent = message.sender
		node.neighbors[message.sender] = 3
		nodeHasChangedLevel(node, complexity)

	case "cityCheck":
		cityCheckLogic(node, message.sender, message.level, message.city, complexity, true)

	case "cityCheckReply":
		if message.answer == "internal" {
			node.neighbors[message.sender] = 1
			allVisited := true
			for _, v := range node.neighbors {
				if v == 0 {
					allVisited = false
					queryNonChildren(node, complexity)
				}
			}
			if allVisited {
				node.foundMySmallestExternalEdge = true
			}
		} else {
			node.foundMySmallestExternalEdge = true
			if message.sender < node.smallestExternalEdgeFound.payload {
				node.smallestExternalEdgeFound = *message
				node.smallestExternalEdgeFound.payload = message.sender
			}
		}
		if node.fringeEdgeFoundResponceCount == node.chidlrenCount && node.foundMySmallestExternalEdge {
			outMessage := Message{catagory: "mergeRequest", level: node.level, city: node.city, destinationPath: node.smallestExternalEdgeFound.callbackPath, payload: len(node.smallestExternalEdgeFound.callbackPath.edges), payload2: node.name}
			if outMessage.payload == 0 {
				node.waitingForReply = true
				node.nodesIveRequested[node.smallestExternalEdgeFound.sender] = 1
			}
			//no longer searching
			node.searchingForFringEdge = false
			sendMessage(node, node.smallestExternalEdgeFound.sender, &outMessage, complexity)
			node.foundMySmallestExternalEdge = false
			node.fringeEdgeFoundResponceCount = 0
		}

	case "terminationBroadcast":
		panic("downtown recieved termination broadcast")

	}
}
