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

func getCubes(outs []*Output) [][]*cube.Cube {
	results := make([][]*cube.Cube, len(outs))
	var wg sync.WaitGroup // WaitGroup to wait for all goroutines to finish

	for i, o := range outs {
		wg.Add(1)   // Increment the WaitGroup counter
		go func() { // Launch a goroutine
			defer wg.Done()
			input := make([]*cube.Cube, 0, len(o.Ones)+len(o.DontCares))
			for _, one := range o.Ones {
				input = append(input, cube.New(one))
			}
			for _, dc := range o.DontCares {
				input = append(input, cube.New(dc))
			}

			res := qm.GetCubes(input)

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

func Formalize(outs []*Output) (string, map[string]*cube.Cube) {
	cubes := []*cube.Cube{}

	// heavy computation
	results := getCubes(outs)

	// problem generation

	constraints := []string{}
	mapping := map[string]*cube.Cube{}
	reverseMapping := map[*cube.Cube]string{}
	iCount := 1

	for i, o := range outs {
		res := results[i]

		for _, one := range o.Ones {
			covers := []string{}
			for i, c := range res {
				if c.Covers(one) {
					key := fmt.Sprintf("x%d", i+iCount)
					covers = append(covers, key)
					mapping[key] = c
					reverseMapping[c] = key
				}
			}

			constraints = append(constraints, fmt.Sprintf("c%d: ", len(constraints))+strings.Join(covers, "+")+" >= 1 ")
		}
		iCount += len(res)

		cubes = append(cubes, res...)
	}

	uniqueCubes := uniqueCubes(cubes)

	cost := []string{}
	for i, c := range uniqueCubes {
		key := fmt.Sprintf("z%d", i+1)
		mapping[key] = c.Cube

		refs := []string{}
		for _, c := range c.refs {
			refs = append(refs, fmt.Sprintf("%.5f ", 1/(float64(len(outs))+1))+reverseMapping[c])
		}

		constraints = append(constraints, fmt.Sprintf("c%d: ", len(constraints))+fmt.Sprintf("%s - %s >= 0", key, strings.Join(refs, "-")))
		cost = append(cost, key)
	}

	template := strings.ReplaceAll(template, "{obj}", " obj: "+strings.Join(cost, "+"))
	template = strings.ReplaceAll(template, "{constraints}", " "+strings.Join(constraints, "\n "))
	template = strings.ReplaceAll(template, "{bounds}", " "+strings.Join(maps.Keys(mapping), "\n "))
	return template, mapping
}
