package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/just-hms/mobo/pkg/mobo"
	"github.com/just-hms/mobo/pkg/opt"
)

func main() {
	tests := 10
	for range tests {
		for nOuts := 1; nOuts <= 8; nOuts++ {
			for nIn := 1; nIn <= 8; nIn++ {
				for onesRatio := 0.1; onesRatio < 0.50+1e-3; onesRatio += 0.1 {

					seed := rand.Int()
					outs := mobo.TestOutputs(nOuts, nIn, onesRatio, seed)

					ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
					defer cancel()

					done := make(chan struct{})

					go func() {

						start := time.Now()
						_, _, globalCost := mobo.Solve(outs)
						duration := time.Since(start)

						totalCost := 0.0
						for _, out := range outs {
							singleOut := []*opt.Output{out}
							_, _, cost := mobo.Solve(singleOut)
							totalCost += cost
						}

						fmt.Printf("%d,%d,%d,%.2f,%.2f,%.2f,%d\n", seed, nOuts, nIn, onesRatio, globalCost, totalCost, duration.Milliseconds())
						done <- struct{}{}
					}()

					select {
					case <-done:
					case <-ctx.Done():
						fmt.Printf("%d,%d,%d,%.2f,%d,%d,timed out\n", seed, nOuts, nIn, onesRatio, -1, -1)
					}

				}
			}
		}
	}
}
