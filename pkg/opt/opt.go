package opt

import (
	"fmt"
	"slices"
	"strings"

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

func Formalize(outs []*Output) (string, map[string]*cube.Cube) {
	cubes := []*cube.Cube{}

	// problem generation ----
	constraints := []string{}
	mapping := map[string]*cube.Cube{}
	reverseMapping := map[*cube.Cube]string{}
	iCount := 0
	// problem generation ----

	for _, o := range outs {

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

		// problem generation ----
		covers := []string{}
		for _, one := range o.Ones {
			for i, c := range res {
				if c.Covers(one) {
					key := fmt.Sprintf("x%d", i+iCount)
					covers = append(covers, key)
					mapping[key] = c
					reverseMapping[c] = key
				}
			}

			constraints = append(constraints, strings.Join(covers, " + ")+" > 1")
		}
		iCount += len(res)
		constraints = append(constraints, "----\n")
		// problem generation ----

		cubes = append(cubes, res...)
	}

	uniqueCubes := uniqueCubes(cubes)

	// problem generation ----
	for i, c := range uniqueCubes {
		key := fmt.Sprintf("x%d", i+iCount)
		mapping[key] = c.Cube
		reverseMapping[c.Cube] = key

		s := ""
		for _, c := range c.refs {
			s += reverseMapping[c] + fmt.Sprintf(" * %.5f", 1/float64(len(outs))+1)
		}

		constraints = append(constraints, fmt.Sprintf("%s > %s", key, s))
	}
	// problem generation ----

	return strings.Join(constraints, "\n"), mapping
}
