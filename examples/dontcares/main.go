package main

import (
	"fmt"

	"github.com/just-hms/mobo/pkg/mobo"
	"github.com/just-hms/mobo/pkg/optimizer"
)

func main() {
	truthTable := []*optimizer.Output{
		{
			Ones:      []uint{2, 3, 7, 12, 15},
			DontCares: []uint{4, 5, 13},
		},
		{
			Ones:      []uint{4, 7, 9, 11, 15},
			DontCares: []uint{6, 12, 14},
		},
	}

	// generate the solution
	circuits, gates, cost := mobo.Solve(truthTable, optimizer.FAN_IN)

	// assert that the solution correctly resembles the truthtable
	err := mobo.Assert(truthTable, circuits)
	if err != nil {
		panic(err)
	}

	fmt.Println("Circuits")
	fmt.Println("------------------------")
	for i, c := range circuits {
		fmt.Printf("%d: ", i+1)
		fmt.Println(c.Display(4))
	}
	fmt.Println()
	fmt.Println("Gates")
	fmt.Println("------------------------")
	for _, gate := range gates {
		fmt.Println(gate.Display(4))
	}
	fmt.Println()
	fmt.Println("Cost:", cost)

}
