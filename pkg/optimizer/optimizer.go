package optimizer

import (
	"fmt"
	"slices"
	"strings"
	"sync"

	"github.com/just-hms/mobo/pkg/qm"
	"github.com/just-hms/mobo/pkg/qm/cube"
	"golang.org/x/exp/maps"
)

type Output struct {
	Ones, DontCares []uint
}

type UniqueCube struct {
	*cube.Cube

	refs []*cube.Cube
}

const template = `Minimize
{obj}

Subject To
{constraints}

Binary
{bounds}

End
`

// uniqueCubes returns a list of unique cubes given a list of cubes
func uniqueCubes(cubes []*cube.Cube) []*UniqueCube {
	uniqueCubes := map[string]*UniqueCube{}
	for _, c := range cubes {
		key := c.String()
		if _, ok := uniqueCubes[key]; !ok {
			uniqueCubes[key] = &UniqueCube{Cube: c.Clone()}
		}
		uniqueCubes[key].refs = append(uniqueCubes[key].refs, c)
	}
	return maps.Values(uniqueCubes)
}

// getCubes retuns a list of all the cubes from a list of thruth tables using QM algorithm
func getCubes(outs []*Output) [][]*cube.Cube {
	results := make([][]*cube.Cube, len(outs))
	var wg sync.WaitGroup

	for i, o := range outs {
		wg.Add(1)
		go func() {
			defer wg.Done()
			input := make([]*cube.Cube, 0, len(o.Ones)+len(o.DontCares))
			for _, one := range o.Ones {
				input = append(input, cube.New(one))
			}
			for _, dc := range o.DontCares {
				input = append(input, cube.New(dc))
			}

			res := qm.Cubes(input)

			// remove cubes which cover only dontCares
			res = slices.DeleteFunc(res, func(c *cube.Cube) bool {
				for _, one := range o.Ones {
					if c.Covers(one) {
						return false
					}
				}
				return true
			})

			results[i] = res
		}()
	}

	wg.Wait()
	return results
}

// Formalize formalizes the problem of sintethizing a boolean function into a ILP problem in .lp format
func Formalize(outs []*Output, cost CostType, size uint) (string, map[string]*cube.Cube) {
	cubes := []*cube.Cube{}

	// variable generation
	results := getCubes(outs)

	// problem formulation
	constraints := []string{}

	// maps a variable to a cube
	mapping := map[string]*cube.Cube{}

	// maps a cube to a variable
	reverseMapping := map[*cube.Cube]string{}

	for i, o := range outs {
		res := results[i]

		for _, one := range o.Ones {
			covers := []string{}
			for j, c := range res {
				if c.Covers(one) {
					key := fmt.Sprintf("v_%d_%d", i+1, j+1)
					covers = append(covers, key)
					mapping[key] = c
					reverseMapping[c] = key
				}
			}

			// add coverage constraints
			constraints = append(constraints, strings.Join(covers, "+")+" >= 1 ")
		}

		cubes = append(cubes, res...)
		constraints = append(constraints, "")
	}

	uniqueCubes := uniqueCubes(cubes)

	for i, c := range uniqueCubes {
		key := fmt.Sprintf("z_%d", i+1)
		mapping[key] = c.Cube

		// get the ref variable of all the cubes the unique cube is linked to and an OR constraint to manage them
		refs := []string{}
		for _, c := range c.refs {
			refs = append(refs, fmt.Sprintf("%.5f ", 1/(float64(len(outs))+1))+reverseMapping[c])
		}

		// add choiche constraints
		constraints = append(constraints, fmt.Sprintf("%s - %s >= 0", key, strings.Join(refs, "-")))
	}

	costFunction := []string{}

	switch cost {
	case GATE_COST:
		for ref := range mapping {
			if strings.HasPrefix(ref, "z") {
				costFunction = append(costFunction, ref)
			}
		}
	case FAN_IN_COST:
		for ref, cube := range mapping {
			// consider every unique cube as +fan_in cost (inputs of AND gate)
			if strings.HasPrefix(ref, "z") {
				fanInCost := fmt.Sprintf("%d %s", cube.FanInCost(size), ref)
				costFunction = append(costFunction, fanInCost)
				continue
			}
			// consider every cube as +1 cost (inputs of the final OR gate)
			costFunction = append(costFunction, ref)

		}
	default:
		panic(fmt.Sprintf("unexpected optimizer.CostType: %#v", cost))
	}

	template := strings.ReplaceAll(template, "{obj}", " obj: "+strings.Join(costFunction, " + "))
	template = strings.ReplaceAll(template, "{constraints}", " "+strings.Join(constraints, "\n "))
	template = strings.ReplaceAll(template, "{bounds}", " "+strings.Join(maps.Keys(mapping), "\n "))

	return template, mapping
}
