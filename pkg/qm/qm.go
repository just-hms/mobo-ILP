package qm

import "github.com/just-hms/mobo/pkg/qm/cube"

func GetCubes(size uint, ones []*cube.Cube) []*cube.Cube {
	groups := make([]map[string]*cube.Cube, size+1)
	// first iteration
	for _, c := range ones {
		i := c.Ones()
		groups[i][c.String()] = c.Clone()
	}

	cubes := ones

	// todo check if any group is not empty
	nextGroupEmpty := false
	for !nextGroupEmpty {
		nextGroupEmpty = true

		nextGroups := make([]map[string]*cube.Cube, size+1)

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

					cubes = append(cubes, m)
					nextGroups[i][m.String()] = m.Clone()
					nextGroupEmpty = false
				}
			}
		}
		groups = nextGroups
	}
	return cubes
}
