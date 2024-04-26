package main

import (
	"fmt"

	"github.com/just-hms/mobo/pkg/mobo"
)

func main() {
	outs := mobo.RandomOutputs(22)
	ports, uniquePorts, solution := mobo.Solve(outs)

	fmt.Println("ports used", uniquePorts)
	fmt.Println("cost", solution)
	fmt.Println("circuits", ports)

	err := mobo.Assert(outs, ports)
	if err != nil {
		panic(err)
	}
}
