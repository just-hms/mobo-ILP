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
	"golang.org/x/sync/errgroup"
)

func Assert(outs []*opt.Output, ports [][]*cube.Cube) error {
	var wg errgroup.Group

	for i, o := range outs {

		wg.Go(func() error {

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
			return nil
		})
	}
	return wg.Wait()
}

func Solve(outs []*opt.Output) ([][]*cube.Cube, []*cube.Cube, float64) {
	problem, cubes := opt.Formalize(outs)

	fmt.Println(problem)

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
	rnd := rand.New(rand.NewSource(int64(seed)))
	size := rnd.Intn(200) + 1
	outputs := make([]*opt.Output, size)
	for i := range outputs {
		rnd = rand.New(rand.NewSource(int64(seed * i)))

		// TODO: change seed
		inputSize := rnd.Intn(200) + 1
		onesRatio := rnd.Float64()
		outputs[i] = &opt.Output{Ones: qm.RandomOnes(inputSize, onesRatio, seed)}
	}
	return outputs
}

func TestOutputs(size int, inputSize int, onesRatio float64, seed int) []*opt.Output {
	outputs := make([]*opt.Output, size)
	for i := range outputs {
		// TODO: change seed better
		outputs[i] = &opt.Output{Ones: qm.RandomOnes(inputSize, onesRatio, seed*i)}
	}
	return outputs
}
