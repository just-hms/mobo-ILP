package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/just-hms/mobo/pkg/mobo"
	"github.com/just-hms/mobo/pkg/opt"
)

func main() {
	test := 10
	for tt := range test {
		for nOuts := range 10 {
			for nIn := range 10 {
				for onesRatio := 0.1; onesRatio < 0.50+1e-3; onesRatio += 0.1 {

					seed := rand.Int()
					outs := mobo.TestOutputs(nOuts, nIn, onesRatio, seed)

					start := time.Now()
					ports, _, globalCost := mobo.Solve(outs)
					duration := time.Since(start)

					err := mobo.Assert(outs, ports)
					if err != nil {
						panic(fmt.Errorf("err: %s, seed: %d", err, seed))
					}

					totalCost := 0.0
					for outIdx, out := range outs {
						singleOut := []*opt.Output{out}
						ports, _, cost := mobo.Solve(singleOut)
						err := mobo.Assert(singleOut, ports)
						if err != nil {
							panic(fmt.Errorf("err: %s, seed: %d, out: %d", err, seed, outIdx))
						}
						totalCost += cost
					}

					fmt.Printf("%d,%d,%d,%d,%.2f,%.2f,%.2f,%d\n", tt, seed, nOuts, nIn, onesRatio, globalCost, totalCost, duration.Milliseconds())
				}
			}
		}
	}
}
