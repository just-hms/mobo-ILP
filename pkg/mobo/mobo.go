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
	"github.com/just-hms/mobo/pkg/qm"
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

			binary := strconv.FormatInt(int64(j), 2)
			if mustBeCovered && len(coverers) == 0 {
				return fmt.Errorf("out: %d, %s should be covered but is not", i+1, binary)
			}

			if !mustBeCovered && len(coverers) > 0 {
				return fmt.Errorf("out: %d, %s shouldn't be covered but it is %v", i+1, binary, coverers)
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
		if math.Abs(v.Value-1) > 1e-4 {
			continue
		}

		c := cubes[v.Name]

		if strings.HasPrefix(v.Name, "z") {
			uniquePorts = append(uniquePorts, c)
			continue
		}

		i, err := strconv.Atoi(strings.Split(v.Name, "_")[1])
		if err != nil {
			panic(err)
		}
		i -= 1
		solutions[i] = append(solutions[i], c)
	}

	return solutions, uniquePorts, math.Ceil(sol.Header.ObjectiveValue)
}

func RandomOutputs(seed int) []*opt.Output {
	rand := rand.New(rand.NewSource(int64(seed)))
	n := rand.Intn(16) + 1
	outputs := make([]*opt.Output, n)
	for i := range outputs {
		outputs[i] = &opt.Output{Ones: qm.RandomOnes(seed)}
	}
	return outputs
}
