package main

import (
	"flag"
	"fmt"
	"math"
	"sync"
	//"gg"
)

//var nodes map[int]Node = make(map[int]Node)

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
			fmt.Println("Error:\n", r)
		}
	}()

	flag.Parse()
	fmt.Println("procedureNum:", procedureNum, ", filePath:", filePath, ", withWeight:", withWeight)

	switch procedureNum {
	case 0:
		runFromFile(filePath, withWeight, nodeNum)
		break
	case 1:
		procedure1(nodeNum)
		break
	case 2:
		procedure2()
	}
}

func runFromFile(filepath string, withWeight bool, nodeNum int) {
	var nodes map[int]Node = make(map[int]Node)
	fileSetup(filePath, withWeight, nodeNum, nodes)
	complexity := 0
	runAlgorithm(nodes, &complexity)
}

func procedure1(n int) [4]int {
	var averageComplexities [4]int
	m := n
	for i := 0; i < 4; i++ {
		switch i {
		case 0:
			m = n
			break
		case 1:
			m = n * int(math.Round(math.Log2(float64(n))))
			break
		case 2:
			m = n * int(math.Sqrt(float64(n)))
			break
		case 3:
			m = n ^ 2
			break
		}
		totalComplexity := 0
		for j := 0; j < 1000; j++ {
			var nodes map[int]Node = make(map[int]Node)
			generateRandomGraph(n, m, n, nodes)
			complexity := 0
			runAlgorithm(nodes, &complexity)
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
			generateRandomGraph(n, m, n, nodes)
			complexity := 0
			runAlgorithm(nodes, &complexity)
			totalComplexity += complexity
		}
		averageComplexities[i] = totalComplexity / 1000
	}
	return averageComplexities
}

func runAlgorithm(nodes map[int]Node, complexity *int) {
	var wg sync.WaitGroup
	for _, node := range nodes {
		wg.Go(func() { instructions(&node, complexity) })
	}

	wg.Wait()

	VisualizeGraph(nodes, "network")
}
