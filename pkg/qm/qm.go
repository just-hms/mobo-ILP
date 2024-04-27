package qm

import (
	"math/rand"
	"slices"

	"github.com/just-hms/mobo/pkg/qm/cube"
	"golang.org/x/exp/maps"
)

func initGroups(size int) []map[string]*cube.Cube {
	groups := make([]map[string]*cube.Cube, size+1)
	for i := range groups {
		groups[i] = make(map[string]*cube.Cube)
	}
	return groups
}

func Cubes(ones []*cube.Cube) []*cube.Cube {
	size := slices.MaxFunc(ones, func(a, b *cube.Cube) int { return a.Len() - b.Len() }).Len()

	groups := initGroups(size)

	// first iteration
	for _, c := range ones {
		i := c.Ones()
		groups[i][c.String()] = c.Clone()
	}

	cubes := ones

	for anyMerge := true; anyMerge; {

		anyMerge = false
		nextGroups := initGroups(size)

		// compare the i-group with the (i+1)-group
		// if one implicant merge :
		//   insert it to the next QM iteration i-group

		for i := 0; i < len(groups)-1; i++ {
			g := groups[i]
			gToCompare := groups[i+1]

			if len(g) == 0 || len(gToCompare) == 0 {
				continue
			}

			for _, c := range g {
				for _, cToCompare := range gToCompare {

					m, err := cube.Merge(c, cToCompare)
					if err != nil {
						continue
					}

					if _, ok := nextGroups[i][m.String()]; !ok {
						nextGroups[i][m.String()] = m.Clone()
						cubes = append(cubes, m)
						anyMerge = true
					}
				}
			}
		}
		groups = nextGroups
	}
	return cubes
}

func RandomOnes(seed int) []uint {
	type fill struct{}

	onesCount := rand.Intn(100) + 1
	ones := make(map[uint]fill, onesCount)
	for range onesCount {
		ones[uint(rand.Intn(200))] = fill{}
	}
	onesList := maps.Keys(ones)
	slices.Sort(onesList)
	return onesList
}
