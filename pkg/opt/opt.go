package opt

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
func Formalize(outs []*Output) (string, map[string]*cube.Cube) {
	cubes := []*cube.Cube{}

	// heavy computation
	results := getCubes(outs)

	// problem generation
	constraints := []string{}
	mapping := map[string]*cube.Cube{}
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

			constraints = append(constraints, strings.Join(covers, "+")+" >= 1 ")
		}

		cubes = append(cubes, res...)
		constraints = append(constraints, "")
	}

	uniqueCubes := uniqueCubes(cubes)

	cost := []string{}
	for i, c := range uniqueCubes {
		key := fmt.Sprintf("z_%d", i+1)
		mapping[key] = c.Cube

		refs := []string{}
		for _, c := range c.refs {
			refs = append(refs, fmt.Sprintf("%.5f ", 1/(float64(len(outs))+1))+reverseMapping[c])
		}

		constraints = append(constraints, fmt.Sprintf("%s - %s >= 0", key, strings.Join(refs, "-")))
		cost = append(cost, key)
	}

	template := strings.ReplaceAll(template, "{obj}", " obj: "+strings.Join(cost, "+"))
	template = strings.ReplaceAll(template, "{constraints}", " "+strings.Join(constraints, "\n "))
	template = strings.ReplaceAll(template, "{bounds}", " "+strings.Join(maps.Keys(mapping), "\n "))
	return template, mapping
}
