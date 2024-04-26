package mobo

import (
	"fmt"
	"math"
	"math/rand"
	"slices"
	"strconv"
	"strings"

	"github.com/just-hms/mobo/pkg/bin"
	"github.com/just-hms/mobo/pkg/cplex"
	"github.com/just-hms/mobo/pkg/opt"
	"github.com/just-hms/mobo/pkg/qm/cube"
)

func Assert(outs []*opt.Output, ports [][]*cube.Cube) error {
	for i, o := range outs {
		cubes := ports[i]

		_max := slices.MaxFunc(slices.Concat(o.Ones, o.DontCares), func(a, b uint) int {
			return int(a) - int(b)
		})

		for j := range bin.NextPowerOf2(_max) {
			mustBeCovered := slices.Contains(o.Ones, j)
			coverers := make([]*cube.Cube, 0)
			for _, c := range cubes {
				if c.Covers(j) {
					coverers = append(coverers, c)
					break
				}
			}

			if mustBeCovered && len(coverers) == 0 {
				return fmt.Errorf("out: %d, %d should be covered but is not", i+1, j)
			}

			if !mustBeCovered && len(coverers) > 0 {
				return fmt.Errorf("out: %d, %d shouldn't be covered but it is %v", i+1, j, coverers)
			}
		}
	}
	return nil
}

func Solve(outs []*opt.Output) ([][]*cube.Cube, []*cube.Cube, float64) {
	problem, cubes := opt.Formalize(outs)

	sol, err := cplex.Solve(problem)
	if err != nil {
		panic(err)
	}

	solutions := make([][]*cube.Cube, len(outs))
	uniquePorts := make([]*cube.Cube, 0, len(outs))
	for _, v := range sol.Variables {
		if math.Abs(v.Value-1) > 1e-3 {
			continue
		}

		c := cubes[v.Name]

		if strings.HasPrefix(v.Name, "z") {
			uniquePorts = append(uniquePorts, c)
			continue
		}

		i, err := strconv.Atoi(string(v.Name[1]))
		if err != nil {
			panic(err)
		}
		i -= 1
		solutions[i] = append(solutions[i], c)
	}

	return solutions, uniquePorts, math.Ceil(sol.Header.ObjectiveValue)
}

func RandomOutputs(seed int64) []*opt.Output {
	rand := rand.New(rand.NewSource(seed))
	n := rand.Intn(10) + 1 // At least 1 output, up to 10.
	outputs := make([]*opt.Output, n)
	for i := range outputs {
		onesCount := rand.Intn(100) + 1 // Each output has between 1 and 10 Ones.
		ones := make([]uint, onesCount)
		for j := range ones {
			ones[j] = uint(rand.Intn(200)) // Random uint from 0 to 99.
		}
		outputs[i] = &opt.Output{Ones: ones}
	}
	return outputs
}
