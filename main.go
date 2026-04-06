package main

import (
	"flag"
	"fmt"
	"math"
	"sync"

	//"gg"
	"runtime/debug"
)

var (
	procedureNum int
	filePath     string
	withWeight   bool
	nodeNum      int
)

func init() {
	flag.IntVar(&procedureNum, "procedureNum", 0, "0 for file, 1 for procedure1, 2 for procedure2")
	flag.StringVar(&filePath, "filePath", "default", "name of file with graph")
	flag.BoolVar(&withWeight, "withWeight", false, "Wether the file has weights")
	flag.IntVar(&nodeNum, "nodeNum", 2, "How many nodes")
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Error:", r)
			debug.PrintStack()
		}
	}()

	flag.Parse()
	fmt.Println("procedureNum:", procedureNum, ", filePath:", filePath, ", withWeight:", withWeight)
	setupOutputFolder()

	switch procedureNum {
	case 0:
		runFromFile()
	case 1:
		procedure1(nodeNum)
	case 2:
		procedure2()
	}
}

func runFromFile() {
	var nodes map[int]Node = make(map[int]Node)
	//note here nodenum is the amount of initiators
	fileSetup(filePath, withWeight, nodeNum, nodes)
	VisualizeGraph(nodes, "network")
	complexity := 0
	leader := -1
	//runAlgorithm(nodes, &complexity, &leader)
	fmt.Printf("algorithm terminated\n")
	fmt.Printf("leader: %d\n", leader)
	fmt.Printf("complexity: %d\n", complexity)
}

func procedure1(n int) [4]int {
	var averageComplexities [4]int
	m := n
	for i := 0; i < 4; i++ {
		switch i {
		case 0:
			m = n
		case 1:
			m = n * int(math.Round(math.Log2(float64(n))))
		case 2:
			m = n * int(math.Sqrt(float64(n)))
		case 3:
			m = n ^ 2
		}
		totalComplexity := 0
		for j := 0; j < 1000; j++ {
			var nodes map[int]Node = make(map[int]Node)
			//note we currently make every node an initiator
			generateRandomGraph(n, m, n, nodes)
			complexity := 0
			leader := -1
			runAlgorithm(nodes, &complexity, &leader)
			totalComplexity += complexity
		}
		averageComplexities[i] = totalComplexity / 1000
	}
	return averageComplexities
}

func procedure2() [4]int {
	var averageComplexities [4]int
	edgeGrowth := func(n int) int {
		return 3 * n
	}
	var nValues = [...]int{20, 30, 40, 60, 80, 100}

	for i := 0; i < 4; i++ {
		n := nValues[i]
		m := edgeGrowth(n)
		totalComplexity := 0
		for j := 0; j < 1000; j++ {
			var nodes map[int]Node = make(map[int]Node)
			//note we currently make every node an initiator
			generateRandomGraph(n, m, n, nodes)
			complexity := 0
			leader := -1
			runAlgorithm(nodes, &complexity, &leader)
			totalComplexity += complexity
		}
		averageComplexities[i] = totalComplexity / 1000
	}
	return averageComplexities
}

func runAlgorithm(nodes map[int]Node, complexity *int, leader *int) {
	//VisualizeGraph(nodes, "network")
	var wg sync.WaitGroup
	for _, node := range nodes {
		wg.Go(func() { instructions(&node, complexity, leader) })
	}

	wg.Wait()
}
