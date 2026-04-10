package main

import (
	"bufio"
	"fmt"
	"math/rand/v2"
	"os"
	"strconv"
	s "strings"
)

func generateRandomGraph(nodesNum int, connectionsNum int, initiatorNum int, nodes map[int]Node) {
	for i := 0; i < nodesNum; i++ {
		initiator := false
		if i < initiatorNum {
			initiator = true
		}
		nodes[i] = NewNode(i, initiator, nodesNum)
	}
	//first give each node an edge, so they are not isolated
	i := 0
	for i := 1; i < nodesNum; i++ {
		n1 := i
		n2 := rand.IntN(i) // connect to any previous node

		connect(i, n1, n2, nodes)
		i++
	}

	//then assign the rest randomly
	for i < connectionsNum {
		n1 := rand.IntN(nodesNum)
		for len(nodes[n1].edges) >= nodesNum {
			n1 = rand.IntN(nodesNum)
		}
		n2 := rand.IntN(nodesNum)
		for n2 == n1 || nodes[n1].neighbors[n2] == 1 || nodes[n2].neighbors[n1] == 1 {
			n2 = rand.IntN(nodesNum)
		}
		//fmt.Println(n1, ",", n2)
		//fmt.Println(nodes[n1].neighbors[n2])
		connect(i, n1, n2, nodes)
		i++
	}

	/*
		components := getComponents(nodes, nodesNum)

		if len(components) > 1 {
			connectComponents(nodes, components, connectionsNum)
		}
	*/

	remakeNodeNeighbors(nodes)

}

func getComponents(nodes map[int]Node, nodesNum int) [][]int {
	visited := make(map[int]bool)
	var components [][]int

	for i := 0; i < nodesNum; i++ {
		if visited[i] {
			continue
		}

		stack := []int{i}
		var component []int
		visited[i] = true

		for len(stack) > 0 {
			curr := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			component = append(component, curr)

			for neighbor := range nodes[curr].neighbors {
				if nodes[curr].neighbors[neighbor] == 1 && !visited[neighbor] {
					visited[neighbor] = true
					stack = append(stack, neighbor)
				}
			}
		}

		components = append(components, component)
	}

	return components
}

func connectComponents(nodes map[int]Node, components [][]int, connectionsNum int) {
	for i := 0; i < len(components)-1; i++ {
		c1 := components[i]
		c2 := components[i+1]

		n1 := c1[rand.IntN(len(c1))]
		n2 := c2[rand.IntN(len(c2))]

		connect(i+connectionsNum, n1, n2, nodes)
	}
}

func fileSetup(filePath string, withWeight bool, initiatorNum int, nodes map[int]Node) {
	file, err := os.Open(filePath)
	check(err)
	defer file.Close()
	scanner := bufio.NewScanner(file)

	i := 0
	for scanner.Scan() {
		line := s.Split(scanner.Text(), ",")
		if (len(line) != 3 && withWeight) || (len(line) < 2 && !withWeight) {
			panic("invalid file format")
		}
		if i != 0 {
			//fmt.Println("line:", line)
			n1, err := strconv.Atoi(line[0])
			check(err)

			if _, ok := nodes[n1]; !ok {
				//fmt.Println("not ok", n1)
				initiator := false
				if len(nodes) < initiatorNum {
					initiator = true
				}
				nodes[n1] = NewNode(n1, initiator, 15)
			}
			n2, err := strconv.Atoi(line[1])
			check(err)
			if _, ok := nodes[n2]; !ok {
				//fmt.Println("not ok", n2)
				initiator := false
				if len(nodes) < initiatorNum {
					initiator = true
				}
				nodes[n2] = NewNode(n2, initiator, 15)
			}
			connect(i, n1, n2, nodes)
		}
		i++
	}
	remakeNodeNeighbors(nodes)
}

func connect(i int, n1 int, n2 int, nodes map[int]Node) {
	fmt.Println(i)
	node1 := nodes[n1]
	node2 := nodes[n2]
	v := Vertex{id: i, node1: &node1, node2: &node2}
	node1.edges[i] = v
	node1.neighbors[node2.name] = 1
	nodes[n1] = node1

	node2.edges[i] = v
	node2.neighbors[node1.name] = 1
	nodes[n2] = node2
	//fmt.Println(n1, n2, node1.edges)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func remakeNodeNeighbors(nodes map[int]Node) {
	for k, v := range nodes {
		v.neighbors = make(map[int]int)
		for ek, _ := range v.edges {
			v.neighbors[ek] = 0
		}
		nodes[k] = v
	}
}
