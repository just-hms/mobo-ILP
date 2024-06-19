package main

import (
	"fmt"
	"math/rand"

	"github.com/just-hms/mobo/pkg/mobo"
	"github.com/just-hms/mobo/pkg/opt"
)

func main() {
	for nOuts := range 2 {
		for nIn := range 2 {
			for k := range 10 {
				onesRatio := float64(k+1) / 10

				outs := mobo.TestOutputs(nOuts, nIn, onesRatio, rand.Int())

				ports, _, globalCost := mobo.Solve(outs)

				err := mobo.Assert(outs, ports)
				if err != nil {
					panic(err)
				}

				totalCost := 0.0
				for _, out := range outs {
					ports, _, cost := mobo.Solve([]*opt.Output{out})
					err := mobo.Assert(outs, ports)
					if err != nil {
						panic(err)
					}
					totalCost += cost
				}

				fmt.Printf("%d,%d,%.2f,%.2f\n", nOuts, nIn, globalCost, totalCost)
			}
		}
	}
}
