package main

import "fmt"

func instructions(node *Node, complexity *int) {
	if node.initiator {
		start(node)
	}

	for true {

		message := <-node.inbox
		//senderIndex := message.sender

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

// target is the index of the edge in node.edges
func sendMessage(node *Node, target int, complexity *int) {
	*complexity++

}

func start(node *Node) {
	node.state = "Downtown"
	//path := findSmallestExternalEdge(node)

}
