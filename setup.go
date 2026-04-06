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
		nodes[i] = NewNode(i, initiator)
	}
	//first give each node an edge, so they are not isolated
	i := 0
	for i < nodesNum {
		n1 := i
		n2 := rand.IntN(nodesNum)
		for n2 == n1 || nodes[n1].neighbors[n2] == 1 || nodes[n2].neighbors[n1] == 1 {
			n2 = rand.IntN(nodesNum)
		}
		fmt.Println(n1, ",", n2)
		fmt.Println(nodes[n1].neighbors[n2])
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

	remakeNodeNeighbors(nodes)

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
				nodes[n1] = NewNode(n1, initiator)
			}
			n2, err := strconv.Atoi(line[1])
			check(err)
			if _, ok := nodes[n2]; !ok {
				//fmt.Println("not ok", n2)
				initiator := false
				if len(nodes) < initiatorNum {
					initiator = true
				}
				nodes[n2] = NewNode(n2, initiator)
			}
			connect(i, n1, n2, nodes)
		}
		i++
	}
	remakeNodeNeighbors(nodes)
}

func connect(i int, n1 int, n2 int, nodes map[int]Node) {
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
