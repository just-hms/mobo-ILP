package qm

import (
	"math/rand"
	"slices"

	"github.com/just-hms/mobo/pkg/qm/cube"
)

// initGroups initialize the QM groups
func initGroups(size int) []map[string]*cube.Cube {
	groups := make([]map[string]*cube.Cube, size+1)
	for i := range groups {
		groups[i] = make(map[string]*cube.Cube)
	}
	return groups
}

// Cubes is a variant of the QM algorithm which return all the implicant(cubes), not just the prime ones
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

// RandomOnes generate a slice of ones with the provided info
func RandomOnes(size int, onesRatio float64, seed int) []uint {
	if onesRatio == 0 {
		return []uint{}
	}

	rand := rand.New(rand.NewSource(int64(seed)))

	_max := 1 << uint(size)

	ones := make([]uint, _max)
	for i := range _max {
		ones[i] = uint(i)
	}
	rand.Shuffle(len(ones), func(i, j int) {
		ones[i], ones[j] = ones[j], ones[i]
	})

	// set a least one 1 if ones ratio is not 0
	onesNumber := max(1, int(float64(_max)*onesRatio))
	ones = ones[:onesNumber]

	slices.Sort(ones)

	return ones
}
