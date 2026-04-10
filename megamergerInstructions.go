package main

import (
	"fmt"
	"math"
)

func instructions(node *Node, complexity *int, leader *int) {
	//fmt.Println("node starting", node.name)
	if node.initiator {
		start(node, complexity)
	}

	for true {

		//fmt.Println("Node", node.name, "level", node.level, "waiting...")
		message := <-node.inbox
		//fmt.Println("Node", node.name, "level", node.level, "received", message.catagory, "from v:", message.sender)

		switch node.state {
		case "Downtown":
			downTownInstructions(node, &message, complexity)
		case "Village":
			villageInstructions(node, &message, complexity)
		case "Asleep":
			fmt.Println("no sleep", node.name)
			//asleepInstructions(node, &message, complexity)
		}
		//not in switch, since this should check after the other function runs
		if node.state == "Done" {
			break
		}
	}
	if node.state == "Downtown" {
		//fmt.Println("I", node.name, "am leader")
	}
	*leader = node.leader
	//fmt.Println(node.name, "done, city=", node.city, "leader=", node.leader)
}

// target is the index of the edge in node.edges
func sendMessage(node *Node, target int, message *Message, complexity *int) {
	*complexity++
	m := cloneMessage(*message)
	m.sender = target
	//fmt.Println(target, node.edges)
	//fmt.Println(node.edges[target])
	if node.edges[target].node1.name == node.name {
		node.edges[target].node2.inbox <- m
		fmt.Println(node.name, node.edges[target].node2.name, m)
	} else {
		node.edges[target].node1.inbox <- m
		fmt.Println(node.name, node.edges[target].node1.name, m)
	}
}

func broadcast(node *Node, message *Message, ignore int, complexity *int) {
	for k, v := range node.neighbors {
		if v == 2 && k != ignore {
			sendMessage(node, k, message, complexity)
		}
	}
}

func sendUpSmallestFringeEdgeFound(node *Node, complexity *int) {
	sendMessage(node, node.parent, &node.smallestExternalEdgeFound, complexity)
	node.fringeEdgeFoundResponceCount = 0
	node.foundMySmallestExternalEdge = false
	node.smallestExternalEdgeFound = Message{catagory: "smallestFringeEdgeFound", smallestFringeEdgeFoundNum: math.MaxInt, callbackPath: EdgePath{edges: []int{}}, payload2: node.name}
}

func queryNonChildren(node *Node, compexity *int) {
	mink := math.MaxInt
	for k, v := range node.neighbors {
		if v == 0 {
			if k < mink {
				mink = k
			}
		}
	}
	if mink != math.MaxInt {
		outMessage := Message{catagory: "cityCheck", level: node.level, city: node.city}
		sendMessage(node, mink, &outMessage, compexity)
	} else {
		//fmt.println("node has no children but called queryNonChildren", node.name)
		node.foundMySmallestExternalEdge = true
	}
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
	//friendly merge, since this is a request from the node I've already requests/
	oldParent := node.parent
	node.neighbors[oldParent] = 2
	node.chidlrenCount++
	//fmt.Println("freindlyV", node.name, sender, payload2, oldParent)
	node.city = sender
	node.level++
	//fmt.Println("parent2", oldParent, node.parent)
	outMessage := Message{catagory: "mergeAccepted", level: node.level, city: node.city}
	broadcast(node, &outMessage, -1, complexity)
	if node.name < payload2 {
		node.state = "Downtown"
		node.neighbors[sender] = 2
		node.chidlrenCount++
		//fmt.Println("I'm in charge", node.name)
	} else {
		node.parent = sender
		node.neighbors[sender] = 3
	}

	if node.waitingForReply == true && node.nodesIveRequested[sender] == 1 {
		//node.waitingForReply = false
		node.nodesIveRequested[sender] = 0
	}
	if node.waitingToReply == true && node.nodesThatHaveRequestedMe[payload2] == 1 {
		node.nodesThatHaveRequestedMe[payload2] = 0
	}
	nodeHasChangedLevel(node, complexity)
}

func friendlyMergeDowntown(node *Node, sender int, payload2 int, complexity *int) {
	//here we broadcast first so only the old children recieve it, not the possible new children from the merge
	node.city = sender
	node.level++
	outMessage := Message{catagory: "mergeAccepted", level: node.level, city: node.city}
	broadcast(node, &outMessage, -1, complexity)
	//fmt.Println("freindlyD", node.name, sender, payload2)
	if node.name <= payload2 {
		node.state = "Downtown"
		node.neighbors[sender] = 2
		node.chidlrenCount++
		//fmt.Println("I'm in charge", node.name)
	} else {
		node.parent = sender
		node.neighbors[sender] = 3
		node.state = "Village"
	}
	if node.waitingForReply == true && node.nodesIveRequested[sender] == 1 {
		//node.waitingForReply = false
		node.nodesIveRequested[sender] = 0
	}
	if node.waitingToReply == true && node.nodesThatHaveRequestedMe[payload2] == 1 {
		node.nodesThatHaveRequestedMe[payload2] = 0
	}
	nodeHasChangedLevel(node, complexity)
}

func cityCheckLogic(node *Node, sender int, level int, city int, complexity *int, setWait bool) bool {
	if city == node.city {
		outMessage := Message{catagory: "cityCheckReply", level: node.level, city: node.city, answer: "internal", payload2: node.name}
		sendMessage(node, sender, &outMessage, complexity)
		return true
	} else if node.level >= level {
		outMessage := Message{catagory: "cityCheckReply", level: node.level, city: node.city, answer: "external", payload2: node.name}
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
		// loop through pending merge requests (map)
		for key, mr := range node.pendingMergeRequests {

			if node.nodesThatHaveRequestedMe[mr.payload2] != 1 {
				delete(node.pendingMergeRequests, key)

				if len(node.pendingMergeRequests) == 0 {
					node.waitingToReply = false
					break
				}
				continue
			}

			if mr.level < node.level {
				sendGetAbsorbedToExternal(node, mr.sender, complexity)
				delete(node.pendingMergeRequests, key)

				if len(node.pendingMergeRequests) == 0 {
					node.waitingToReply = false
					break
				}
			}
		}
	}
	if node.waitingToReplyToCityCheck == true {
		//fmt.Println("chsking cty reply", node.name)
		for i := 0; i < len(node.pendingCityChecks); i++ {
			//fmt.Println("chsking cty reply22", node.name, node.level, node.pendingCityChecks[i].level)
			if node.pendingCityChecks[i].level == node.level {
				didSend := cityCheckLogic(node, node.pendingCityChecks[i].sender, node.pendingCityChecks[i].level, node.pendingCityChecks[i].city, complexity, false)
				if didSend {
					node.pendingCityChecks = removeCC(node.pendingCityChecks, i)
					if len(node.pendingCityChecks) == 0 {
						node.waitingToReplyToCityCheck = false
						return
					}
					i--
				} else {
					panic("should have sent the city check")
				}
			} else if node.pendingCityChecks[i].level < node.level {
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
			outMessage := Message{catagory: "smallestFringeEdgeFound", smallestFringeEdgeFoundNum: math.MaxInt, level: node.level, city: node.city, callbackPath: EdgePath{edges: []int{}}, payload2: node.name}
			sendMessage(node, node.parent, &outMessage, complexity)
		}
		if len(node.neighbors)-node.chidlrenCount == 1 {
			//fmt.Println("node has no external neighbors", node.name)
			node.foundMySmallestExternalEdge = true
		}
		outMessage := Message{catagory: "findSmallestFringeEdge", level: node.level, city: node.city, payload: math.MaxInt, callbackPath: EdgePath{edges: []int{}}}
		broadcast(node, &outMessage, -1, complexity)
		if !node.foundMySmallestExternalEdge {
			queryNonChildren(node, complexity)
		}

	case "smallestFringeEdgeFound":
		if message.sender == node.parent {
			panic("village recieved smallest fringe edge from parent")
		}
		node.fringeEdgeFoundResponceCount++
		//fmt.Println("smallestFringeEdgeVFound", "name:", node.name, "childCount:", node.chidlrenCount, "fringeResponse:", node.fringeEdgeFoundResponceCount, "foundMySEdge:", node.foundMySmallestExternalEdge, "neighbors:", node.neighbors)
		if message.smallestFringeEdgeFoundNum <= node.smallestExternalEdgeFound.smallestFringeEdgeFoundNum {
			node.smallestExternalEdgeFound = *message
			node.smallestExternalEdgeFound.callbackPath.edges = append(node.smallestExternalEdgeFound.callbackPath.edges, message.sender)
		}
		if node.fringeEdgeFoundResponceCount == node.chidlrenCount && node.foundMySmallestExternalEdge {
			//fmt.Println("sending up found edge", "name:", node.name, "payload:", node.smallestExternalEdgeFound.smallestFringeEdgeFoundNum)
			sendUpSmallestFringeEdgeFound(node, complexity)
			//no longer searching
			node.searchingForFringEdge = false
		}

	case "mergeRequest":
		//fmt.Println("mr recieved:", node.name, node.waitingToReply, node.nodesThatHaveRequestedMe[message.payload2], message.payload2, node.pendingMergeRequests[message.payload2].level, message.level)
		if message.payload > 0 {
			//case where the message needs to be passed on until it reaches the external node
			target := message.destinationPath.Pop()
			if message.payload == 1 {
				if node.waitingToReply == true && node.nodesThatHaveRequestedMe[message.payload2] == 1 && node.pendingMergeRequests[message.payload2].level == message.level {
					friendlyMergeVillage(node, target, message.payload2, complexity)
				} else {
					node.waitingForReply = true
					node.nodesIveRequested[target] = 1
					message.payload2 = node.name
				}
			}
			message.payload2 = node.name
			message.payload--
			sendMessage(node, target, message, complexity)
			if node.state == "Downtown" {
				//fmt.Println("I'm in charge now")
				node.searchingForFringEdge = true
				outMessage2 := Message{catagory: "findSmallestFringeEdge", city: node.city, level: node.level, payload: math.MaxInt}
				broadcast(node, &outMessage2, -1, complexity)
				queryNonChildren(node, complexity)
			}
		} else if node.level > message.level {
			//case where this node absorbes the requester
			sendGetAbsorbedToExternal(node, message.sender, complexity)
		} else if node.waitingForReply == true && node.nodesIveRequested[message.sender] == 1 && node.level == message.level {
			friendlyMergeVillage(node, message.sender, message.payload2, complexity)
			if node.state == "Downtown" {
				//fmt.Println("I'm in charge now")
				node.searchingForFringEdge = true
				outMessage2 := Message{catagory: "findSmallestFringeEdge", city: node.city, level: node.level, payload: math.MaxInt}
				broadcast(node, &outMessage2, -1, complexity)
				queryNonChildren(node, complexity)
			}

		} else {
			node.waitingToReply = true
			node.nodesThatHaveRequestedMe[message.payload2] = 1
			//fmt.Println("lvl", message.level)
			node.pendingMergeRequests[message.sender] = PendingMergeRequest{sender: message.sender, level: message.level, payload2: message.payload2}
		}

	case "mergeAccepted":
		node.level = message.level
		node.city = message.city
		if message.sender == node.parent {
			broadcast(node, message, -1, complexity)
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
		if message.level < node.level {
			panic("village recieved asborb message with smaller level")
		}
		node.city = message.city
		node.level = message.level
		outMessage := Message{catagory: "getAbsorbed", level: message.level, city: message.city}
		if message.sender == node.parent {
			broadcast(node, &outMessage, -1, complexity)
		} else if node.neighbors[message.sender] == 1 {
			panic("village recieved merge accepted from internal non child/parent node")
		} else {
			//case where it's coming from external neighbor or child
			sendMessage(node, node.parent, &outMessage, complexity)
			broadcast(node, &outMessage, message.sender, complexity)
			node.neighbors[node.parent] = 2
			node.parent = message.sender
			if node.neighbors[message.sender] == 0 {
				node.chidlrenCount++
			}
			node.neighbors[message.sender] = 3
		}
		if node.waitingForReply == true && node.nodesIveRequested[message.sender] == 1 {
			node.waitingForReply = false
			node.nodesIveRequested[message.sender] = 0
		}
		nodeHasChangedLevel(node, complexity)

	case "cityCheck":
		//fmt.Println("cityCheck", node.level, node.pendingCityChecks, node.waitingToReplyToCityCheck)
		cityCheckLogic(node, message.sender, message.level, message.city, complexity, true)
		//fmt.Println("4!", didSend, node.level, node.pendingCityChecks, node.waitingToReplyToCityCheck)

	case "cityCheckReply":
		if message.answer == "internal" {
			if node.neighbors[message.sender] == 0 {
				node.neighbors[message.sender] = 1
			}
			allVisited := true
			for _, v := range node.neighbors {
				if v == 0 {
					allVisited = false
					queryNonChildren(node, complexity)
					break
				}
			}
			if allVisited {
				node.foundMySmallestExternalEdge = true
			}
		} else {
			node.foundMySmallestExternalEdge = true
			if message.sender < node.smallestExternalEdgeFound.smallestFringeEdgeFoundNum {
				node.smallestExternalEdgeFound = Message{catagory: "smallestFringeEdgeFound", sender: message.sender, smallestFringeEdgeFoundNum: message.sender, level: node.level, city: node.city, callbackPath: EdgePath{edges: []int{}}, payload2: message.payload2}
				node.smallestExternalEdgeFound.callbackPath.edges = append(node.smallestExternalEdgeFound.callbackPath.edges, message.sender)
				node.smallestExternalEdgeFound.smallestFringeEdgeFoundNum = message.sender
			}
		}
		if node.fringeEdgeFoundResponceCount == node.chidlrenCount && node.foundMySmallestExternalEdge {
			sendUpSmallestFringeEdgeFound(node, complexity)
			//no longer searching
			node.searchingForFringEdge = false
		}

	case "terminationBroadcast":
		broadcast(node, message, -1, complexity)
		node.state = "Done"
		node.leader = message.payload2
	}
}

func downTownInstructions(node *Node, message *Message, complexity *int) {
	switch message.catagory {
	case "findSmallestFringeEdge":
		panic(fmt.Sprintf("Downtown recieved find smallest fringe edge %d", node.name))

	case "smallestFringeEdgeFound":
		node.fringeEdgeFoundResponceCount++
		if message.smallestFringeEdgeFoundNum < node.smallestExternalEdgeFound.smallestFringeEdgeFoundNum {
			//fmt.Println(*message)
			node.smallestExternalEdgeFound = *message
		}
		//fmt.Println("downtown recieved smallestFound", node.name, node.fringeEdgeFoundResponceCount, node.chidlrenCount, node.foundMySmallestExternalEdge, node.smallestExternalEdgeFound)
		if node.fringeEdgeFoundResponceCount == node.chidlrenCount && node.foundMySmallestExternalEdge {
			if node.smallestExternalEdgeFound.smallestFringeEdgeFoundNum == math.MaxInt {
				//fmt.Println("termination, these are my children", node.chidlrenCount, node.neighbors)
				outMessage := Message{catagory: "terminationBroadcast", payload2: node.name}
				broadcast(node, &outMessage, -1, complexity)
				node.state = "Done"
				node.leader = node.name
				break
			}
			//no longer searching
			node.searchingForFringEdge = false
			if node.name == 6 || node.name == 2 {
				//fmt.Println("wtf", node.smallestExternalEdgeFound)
			}
			outMessage := Message{catagory: "mergeRequest", level: node.level, city: node.city, destinationPath: node.smallestExternalEdgeFound.callbackPath, payload: len(node.smallestExternalEdgeFound.callbackPath.edges), payload2: node.smallestExternalEdgeFound.payload2}
			if outMessage.payload == 0 {
				if node.waitingToReply == true && node.nodesThatHaveRequestedMe[message.payload2] == 1 && node.level == message.level {
					friendlyMergeDowntown(node, node.smallestExternalEdgeFound.sender, message.payload2, complexity)

					sendMessage(node, node.smallestExternalEdgeFound.sender, &outMessage, complexity)
					if node.state == "Downtown" {
						//fmt.Println("I'm in charge 438")
						node.searchingForFringEdge = true
						outMessage2 := Message{catagory: "findSmallestFringeEdge", city: node.city, level: node.level, payload: math.MaxInt}
						node.smallestExternalEdgeFound.smallestFringeEdgeFoundNum = math.MaxInt
						broadcast(node, &outMessage2, -1, complexity)
						queryNonChildren(node, complexity)
					}
				} else {
					node.waitingForReply = true
					node.nodesIveRequested[node.smallestExternalEdgeFound.sender] = 1
					//println("setting waitingToReply from frnge found", node.smallestExternalEdgeFound.sender, message.payload2)
					sendMessage(node, node.smallestExternalEdgeFound.sender, &outMessage, complexity)
				}

			} else {
				//fmt.Println("name", node.name)
				sendMessage(node, node.smallestExternalEdgeFound.sender, &outMessage, complexity)
			}
		}

	case "mergeRequest":
		//fmt.Println("mr recieved:", node.name, node.waitingToReply, node.nodesThatHaveRequestedMe[message.payload2], node.waitingForReply, node.nodesIveRequested, node.level, message.level)
		if node.level > message.level {
			//absorbtion case
			sendGetAbsorbedToExternal(node, message.sender, complexity)
		} else if node.waitingForReply == true && node.nodesIveRequested[message.sender] == 1 && node.level == message.level {
			//friendly merge, since this is a request from the node I've already requests
			friendlyMergeDowntown(node, message.sender, message.payload2, complexity)
			if node.state == "Downtown" {
				//fmt.Println("I'm in charge 466")
				node.searchingForFringEdge = true
				outMessage2 := Message{catagory: "findSmallestFringeEdge", city: node.city, level: node.level, payload: math.MaxInt}
				node.smallestExternalEdgeFound.smallestFringeEdgeFoundNum = math.MaxInt
				broadcast(node, &outMessage2, -1, complexity)
				queryNonChildren(node, complexity)
			}

		} else {
			node.waitingToReply = true
			node.nodesThatHaveRequestedMe[message.payload2] = 1
			//println("setting waitingToReply", message.sender, message.payload2, message.level)
			node.pendingMergeRequests[message.sender] = PendingMergeRequest{sender: message.sender, level: message.level, payload2: message.payload2}
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
		broadcast(node, &outMessage, -1, complexity)
		nodeHasChangedLevel(node, complexity)

	case "getAbsorbed":
		if message.level <= node.level {
			panic("downtown recieved asborb message with smaller level")
		} else if node.neighbors[message.sender] == 1 {
			panic("downtown recieved asborb message from non child neighbor")
		}
		node.city = message.city
		node.level = message.level
		node.state = "Village"
		if node.waitingForReply == true && node.nodesIveRequested[message.sender] == 1 {
			node.waitingForReply = false
			node.nodesIveRequested[message.sender] = 0
		}
		if node.neighbors[message.sender] == 2 {
			node.chidlrenCount--
		}
		outMessage := Message{catagory: "getAbsorbed", level: message.level, city: message.city}
		broadcast(node, &outMessage, message.sender, complexity)
		node.parent = message.sender
		node.neighbors[message.sender] = 3
		nodeHasChangedLevel(node, complexity)

	case "cityCheck":
		//fmt.Println("cityCheck", node.level, node.pendingCityChecks, node.waitingToReplyToCityCheck)
		cityCheckLogic(node, message.sender, message.level, message.city, complexity, true)
		//fmt.Println("4!", didSend, node.level, node.pendingCityChecks, node.waitingToReplyToCityCheck)

	case "cityCheckReply":
		if message.answer == "internal" {
			if node.neighbors[message.sender] == 0 {
				node.neighbors[message.sender] = 1
			}
			allVisited := true
			for _, v := range node.neighbors {
				if v == 0 {
					allVisited = false
					queryNonChildren(node, complexity)
					break
				}
			}
			if allVisited {
				node.foundMySmallestExternalEdge = true
			}
		} else {
			node.foundMySmallestExternalEdge = true
			if message.sender < node.smallestExternalEdgeFound.smallestFringeEdgeFoundNum {
				node.smallestExternalEdgeFound = Message{catagory: "smallestFringeEdgeFound", sender: message.sender, smallestFringeEdgeFoundNum: message.sender, level: node.level, city: node.city, callbackPath: EdgePath{edges: []int{}}, payload2: node.name}
				node.smallestExternalEdgeFound.smallestFringeEdgeFoundNum = message.sender
			}
		}
		if node.fringeEdgeFoundResponceCount == node.chidlrenCount && node.foundMySmallestExternalEdge {
			outMessage := Message{catagory: "mergeRequest", level: node.level, city: node.city, destinationPath: node.smallestExternalEdgeFound.callbackPath, payload: len(node.smallestExternalEdgeFound.callbackPath.edges), payload2: node.smallestExternalEdgeFound.payload2}
			if outMessage.payload == 0 {
				if node.waitingToReply == true && node.nodesThatHaveRequestedMe[message.payload2] == 1 && node.level == message.level {
					friendlyMergeDowntown(node, node.smallestExternalEdgeFound.sender, message.payload2, complexity)
					if node.state == "Downtown" {
						//fmt.Println("I'm in charge 552")
						node.searchingForFringEdge = true
						outMessage2 := Message{catagory: "findSmallestFringeEdge", city: node.city, level: node.level, payload: math.MaxInt}
						node.smallestExternalEdgeFound.smallestFringeEdgeFoundNum = math.MaxInt
						broadcast(node, &outMessage2, -1, complexity)
						queryNonChildren(node, complexity)
					}
				} else {
					node.waitingForReply = true
					node.nodesIveRequested[node.smallestExternalEdgeFound.sender] = 1
				}
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
