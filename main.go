package main

import (
	"fmt"

	"github.com/just-hms/mobo/pkg/mobo"
)

func main() {
	outs := mobo.TestOutputs(200, 10, 0.1, 20)

	ports, uniquePorts, solution := mobo.Solve(outs)

	fmt.Println("ports used", uniquePorts)
	fmt.Println("circuits", ports)
	fmt.Println("cost", solution)

	err := mobo.Assert(outs, ports)
	if err != nil {
		panic(err)
	}
}
