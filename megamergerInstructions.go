package main

import "fmt"

func instructions(node *Node, complexity *int) {
	if node.initiator {
		start(node)
	}

	for true {

		message := <-node.inbox

		switch node.state {
		case "Downtown":
			downTownInstructions(node, &message, complexity)
		case "Village":
			villageInstructions(node, &message, complexity)
		case "Asleep":
			downTownInstructions(node, &message, complexity)
		case "Done":
			downTownInstructions(node, &message, complexity)
		}

		if node.state == "done" {
			break
		}
	}
	fmt.Println("done")
}

func findSmallestExternalEdge(node *Node) int { return 0 }

// target is the index of the edge in node.edges
func sendMessage(node *Node, target int, message *Message, complexity *int) {
	*complexity++

}

func broadcast(node *Node, message Message, compexity *int) {

}

func queryNonChildren(node *Node, message Message, compexity *int) {
	for i := 0; i < len(node.neighbors); i++ {
		if _, ok := node.internalNeighbors[i]; !ok {
			outMessage := Message{catagory: "cityCheck", sender: node.internalNeighbors[i], level: node.level, city: node.city}
			sendMessage(node, node.internalNeighbors[i], &outMessage, compexity)
			return
		}
	}
}

func start(node *Node) {
	node.state = "Downtown"
	//path := findSmallestExternalEdge(node)

}

func villageInstructions(node *Node, message *Message, complexity *int) {
	switch message.catagory {
	case "findSmallestFringeEdge":
		callback := edgePath{edges: append(message.callbackPath.edges, message.sender)}
		outMessage := Message{catagory: "findSmallestFringeEdge", sender: message.sender, level: node.level, city: node.city, callbackPath: callback}
		broadcast(node, outMessage, complexity)
		outMessage2 := Message{catagory: "cityCheck", sender: message.sender, level: node.level, city: node.city, callbackPath: callback}
		queryNonChildren(node, outMessage2, complexity)

	case "smallestFringeEdgeFound":
		if message.sender == node.parent {
			panic("village recieved smallest fringe edge from parent")
		}
		if message.payload > node.mySmallestInternalEdge {
			message.payload = node.mySmallestInternalEdge
			message.callbackPath = edgePath{edges: []int{node.parent}}
		} else {
			message.callbackPath.edges = append(message.callbackPath.edges, node.parent)
		}
		sendMessage(node, node.parent, message, complexity)

	case "mergeRequest":
		outMessage := Message{catagory: "getAbsorbed", sender: message.sender, level: node.level, city: node.city, destinationPath: message.callbackPath}
		if node.level > message.level {
			sendMessage(node, outMessage.destinationPath.Pop(), &outMessage, complexity)
			node.children = append(node.children, message.sender)
			node.internalNeighbors[message.sender] = 2
		} else {
			node.substate = "waitingToReply"

		}
		//note consider if if b and c both make a merge request to me
		//use array to store pending request
		//just need to store level, and sender

	case "mergeAccepted":
		outMessage := Message{catagory: "mergeAccepted", sender: message.sender, level: message.level, city: message.city}
		broadcast(node, outMessage, complexity)
		if message.city != node.city {
			node.city = message.city
			node.level = message.level
			node.state = "Village"
			node.parent = message.callbackPath.Peek()
			node.removeChild(message.callbackPath.Peek())
		}

	case "getAbsorbed":
		if message.level <= node.level {
			panic("downtown recieved asborb message with smaller level")
		}
		outMessage := Message{catagory: "getAbsorbed", sender: message.sender, level: message.level, city: message.city}
		broadcast(node, outMessage, complexity)
		node.city = message.city
		node.level = message.level
		node.state = "Village"
		node.parent = message.callbackPath.Peek()
		node.removeChild(message.callbackPath.Peek())

	case "cityCheck":
		if message.city == node.city {
			outMessage := Message{catagory: "cityCheckReply", sender: message.sender, level: node.level, city: node.city, answer: "internal"}
			sendMessage(node, outMessage.destinationPath.Pop(), &outMessage, complexity)
		} else if node.level >= message.level {
			outMessage := Message{catagory: "cityCheckReply", sender: message.sender, level: node.level, city: node.city, answer: "external"}
			sendMessage(node, outMessage.destinationPath.Pop(), &outMessage, complexity)
		} else {
			node.substate = "WatingToReply"
		}

	case "cityCheckReply":
		//need to write

	case "terminationBroadcast":
		panic("downtown recieved termination broadcast")

	}
}

func downTownInstructions(node *Node, message *Message, complexity *int) {
	switch message.catagory {
	case "findSmallestFringeEdge":
		panic("Downtown recieved find smallest fringe edge")

	case "smallestFringeEdgeFound":
		callback := edgePath{edges: []int{message.callbackPath.Peek()}}
		outMessage := Message{catagory: "mergeRequest", sender: message.sender, level: node.level, city: node.city, callbackPath: callback, destinationPath: message.callbackPath}
		sendMessage(node, outMessage.destinationPath.Pop(), &outMessage, complexity)

	case "mergeRequest":
		outMessage := Message{catagory: "getAbsorbed", sender: message.sender, level: node.level, city: node.city, destinationPath: message.callbackPath}
		if node.level > message.level {
			sendMessage(node, outMessage.destinationPath.Pop(), &outMessage, complexity)
			node.children = append(node.children, message.sender)
			node.internalNeighbors[message.sender] = 2
		} else {
			node.substate = "waitingToReply"
		}

	case "mergeAccepted":
		outMessage := Message{catagory: "mergeAccepted", sender: message.sender, level: message.level, city: message.city}
		broadcast(node, outMessage, complexity)
		if message.city != node.city {
			node.city = message.city
			node.level = message.level
			node.state = "Village"
			node.parent = message.callbackPath.Peek()
			node.removeChild(message.callbackPath.Peek())
		}

	case "getAbsorbed":
		if message.level <= node.level {
			panic("downtown recieved asborb message with smaller level")
		}
		outMessage := Message{catagory: "getAbsorbed", sender: message.sender, level: message.level, city: message.city}
		broadcast(node, outMessage, complexity)
		node.city = message.city
		node.level = message.level
		node.state = "Village"
		node.parent = message.callbackPath.Peek()
		node.removeChild(message.callbackPath.Peek())

	case "cityCheck":
		if message.city == node.city {
			outMessage := Message{catagory: "cityCheckReply", sender: message.sender, level: node.level, city: node.city, answer: "internal"}
			sendMessage(node, outMessage.destinationPath.Pop(), &outMessage, complexity)
		} else if node.level >= message.level {
			outMessage := Message{catagory: "cityCheckReply", sender: message.sender, level: node.level, city: node.city, answer: "external"}
			sendMessage(node, outMessage.destinationPath.Pop(), &outMessage, complexity)
		} else {
			node.substate = "WatingToReply"
		}

	case "cityCheckReply":
		//need to write

	case "terminationBroadcast":
		panic("downtown recieved termination broadcast")

	}
}
