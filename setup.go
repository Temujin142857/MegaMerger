package main

import (
	"bufio"
	"fmt"
	"math/rand/v2"
	"os"
	"strconv"
	s "strings"
)

func randomSetup(upperBoundOfNodes int, connectionsNum int, initiatorNum int) {
	for i := 0; i < upperBoundOfNodes; i++ {
		initiator := false
		if i < initiatorNum {
			initiator = true
		}
		nodes[i] = Node{name: i, level: 0, city: -1, parent: -1, state: "asleep", initiator: initiator, neighbors: make(map[int]int)}
	}
	//fmt.Println(nodes)
	for i := 0; i < connectionsNum; i++ {
		n1 := rand.IntN(upperBoundOfNodes)
		for len(nodes[n1].edges) >= upperBoundOfNodes {
			n1 = rand.IntN(upperBoundOfNodes)
		}
		n2 := rand.IntN(upperBoundOfNodes)
		for n2 == n1 || nodes[n1].neighbors[n2] == 1 || nodes[n2].neighbors[n1] == 1 {
			n2 = rand.IntN(upperBoundOfNodes)
		}
		fmt.Println(n1, ",", n2)
		fmt.Println(nodes[n1].neighbors[n2])
		connect(i, n1, n2)
	}
}

func fileSetup(filePath string, withWeight bool, initiatorNum int) {
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
			fmt.Println(line)
			n1, err := strconv.Atoi(line[0])
			check(err)

			if _, ok := nodes[n1]; !ok {
				fmt.Println("not ok", n1)
				initiator := false
				if len(nodes) < initiatorNum {
					initiator = true
				}
				nodes[n1] = Node{name: n1, level: 0, city: -1, parent: -1, state: "asleep", initiator: initiator, neighbors: make(map[int]int)}
			}
			n2, err := strconv.Atoi(line[1])
			check(err)
			if _, ok := nodes[n2]; !ok {
				fmt.Println("not ok", n2)
				initiator := false
				if len(nodes) < initiatorNum {
					initiator = true
				}
				nodes[n2] = Node{name: n2, level: 0, city: -1, parent: -1, state: "asleep", initiator: initiator, neighbors: make(map[int]int)}
			}
			connect(i, n1, n2)
		}
		i++
	}
}

func connect(i int, n1 int, n2 int) {
	node1 := nodes[n1]
	node2 := nodes[n2]
	v := Vertex{name: i, node1: node1.name, node2: node2.name, channel: make(chan Message)}
	node1.edges = append(node1.edges, v)
	node1.neighbors[node2.name] = 1
	nodes[n1] = node1

	node2.edges = append(node2.edges, v)
	node2.neighbors[node1.name] = 1
	nodes[n2] = node2
	fmt.Println(node1.edges)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
