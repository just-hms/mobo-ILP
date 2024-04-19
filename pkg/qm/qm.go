package qm

func GetCubes(size uint, ones []uint) []*Cube {
	bins := make([]*Cube, len(ones))
	for i, o := range ones {
		bins[i] = NewCube(o)
	}

	groups := make([]map[string]*Cube, size+1)

	// first iteration
	for _, bin := range bins {
		i := bin.Count()
		groups[i][bin.DumpAsBits()] = bin.Clone()
	}

	cubes := bins

	// todo check if any group is not empty
	nextGroupEmpty := false
	for !nextGroupEmpty {
		nextGroupEmpty = true

		nextGroups := make([]map[string]*Cube, size+1)

		// compare the i-group with the (i+1)-group
		// if one implicant merge :
		//   insert it to the next QM iteration i-group
		//   insert it into the implicant list and update A

		for i := 0; i < len(groups)-1; i++ {
			g := groups[i]
			gToCompare := groups[i+1]

			if len(g) == 0 || len(gToCompare) == 0 {
				continue
			}

			for _, cube := range g {
				for _, cubeToCompare := range gToCompare {

					m, err := MergeCubes(cube, cubeToCompare)
					if err != nil {
						continue
					}

					if m.Count() == 0 {
						continue
					}

					cubes = append(cubes, m)
					nextGroups[i][m.DumpAsBits()] = m.Clone()
					nextGroupEmpty = false
				}

			}

		}

		groups = nextGroups
	}
	return cubes
}
