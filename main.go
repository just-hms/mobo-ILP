package main

import (
	"fmt"
	"log"
	"math"

	"github.com/just-hms/mobo/pkg/cplex"
	"github.com/just-hms/mobo/pkg/opt"
	"github.com/just-hms/mobo/pkg/qm/cube"
)

func main() {
	outs := []*opt.Output{
		{Ones: []uint{1, 2, 3, 5}},
		{Ones: []uint{1, 5, 6, 7}},
	}
	problem, cubes := opt.Formalize(outs)

	sol, err := cplex.Solve(problem)
	if err != nil {
		log.Fatal(err)
	}

	solution := []*cube.Cube{}
	for _, v := range sol.Variables {
		c := cubes[v.Name]
		if math.Abs(v.Value-1) < 1e-3 {
			solution = append(solution, c)
		}
	}

	fmt.Println(solution)
}
