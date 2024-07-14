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
	"github.com/just-hms/mobo/pkg/optimizer"
	"github.com/just-hms/mobo/pkg/qm"
	"github.com/just-hms/mobo/pkg/qm/cube"
	"golang.org/x/sync/errgroup"
)

// Assert verifies that the generated ports correctly synthetize the provided circuit
func Assert(outs []*optimizer.Output, circuits []Circuit) error {
	var wg errgroup.Group

	if len(outs) != len(circuits) {
		return fmt.Errorf("provided %d output with %d circuits", len(outs), len(circuits))
	}

	for i, o := range outs {
		wg.Go(func() error {

			cubes := circuits[i]

			if len(o.Ones) == 0 {
				if len(cubes) != 0 {
					return fmt.Errorf("circuit should be empty but is not")
				}
				return nil
			}

			_max := slices.MaxFunc(slices.Concat(o.Ones, o.DontCares), func(a, b uint) int {
				return int(a) - int(b)
			})

			for j := range bin.NextPowerOf2(_max) {
				isDontCare := slices.Contains(o.DontCares, j)
				if isDontCare {
					continue
				}

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
					return fmt.Errorf("out: %d, %s shouldn't be covered but is covered by %v", i+1, binary, coverers)
				}
			}
			return nil
		})
	}
	return wg.Wait()
}

func inputSize(outs []*optimizer.Output) uint {
	_max := uint(0)
	for _, o := range outs {
		_max = max(_max, slices.MaxFunc(slices.Concat(o.Ones, o.DontCares), func(a, b uint) int {
			return int(a) - int(b)
		}))
	}
	return bin.MinBitsNeeded(_max)
}

// Solve given a thruth table returns the cube to use in each sub-circuit, the unique gates used and the cost of them using CPLEX
func Solve(outs []*optimizer.Output, cost optimizer.CostType) ([]Circuit, []*cube.Cube, float64) {
	nOnes := 0
	for _, o := range outs {
		nOnes += len(o.Ones)
	}

	if nOnes == 0 {
		return make([]Circuit, len(outs)), make([]*cube.Cube, 0), 0
	}

	size := inputSize(outs)

	problem, cubes := optimizer.Formalize(outs, cost, size)

	sol, err := cplex.Solve(problem)
	if err != nil {
		panic(err)
	}

	solutions := make([]Circuit, len(outs))
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

// RandomOutputs generates a random thruth table
func RandomOutputs(seed int) []*optimizer.Output {
	rnd := rand.New(rand.NewSource(int64(seed)))
	size := rnd.Intn(5) + 2
	outputs := make([]*optimizer.Output, size)
	for i := range outputs {
		rnd = rand.New(rand.NewSource(int64(seed * i)))
		inputSize := rnd.Intn(5) + 2
		onesRatio := rnd.Float64()
		outputs[i] = &optimizer.Output{Ones: qm.RandomOnes(inputSize, onesRatio, seed)}
	}
	return outputs
}

// TestOutputs given outputSize, outputSize and number of ones generates a random truth table
func TestOutputs(outputSize int, inputSize int, onesRatio float64, seed int) []*optimizer.Output {
	outputs := make([]*optimizer.Output, outputSize)
	for i := range outputs {
		outputs[i] = &optimizer.Output{Ones: qm.RandomOnes(inputSize, onesRatio, seed*i)}
	}
	return outputs
}
