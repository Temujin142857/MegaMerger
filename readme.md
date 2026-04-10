# MegaMerger Implementation

## Things to ask:
How should initiators be chosen?
Does the one goroutine per node make sense for simulating a distrubed system

## install instructions
go install github.com/goccy/go-graphviz/cmd/dot@latest
go run . -procedureNum=0 -filePath="C:\Users\Tomio\Programming\MegaMerger\input\edgeList.csv" -nodeNum 9

## Citations:
Didn't end up using anything directly but I took some inspiration from the pseudo code when making my state diagram.

1. Gallager, R. G, et al. “A Distributed Algorithm for Minimum-Weight Spanning Trees.” ACM Transactions on Programming Languages and Systems [New York, NY, USA], vol. 5, no. 1, January 1983, pp. 66–77, https://doi.org/10.1145/357195.357200.
