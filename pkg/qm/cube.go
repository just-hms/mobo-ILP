package qm

import (
	"errors"
	"strings"

	"github.com/bits-and-blooms/bitset"
)

type Cube struct {
	*bitset.BitSet
}

func NewCube(val uint) *Cube {
	c := &Cube{
		BitSet: &bitset.BitSet{},
	}
	if val == 0 {
		return c
	}
	c.BitSet.Set(val - 1)
	return c
}

func (c *Cube) Clone() *Cube {
	return &Cube{
		BitSet: c.BitSet.Clone(),
	}
}

func (c *Cube) Repr(size int) string {
	bits := c.DumpAsBits()
	if bits == "" {
		// TODO: make this better
		bits = "0000000000000000000000000000000000000000000000000000"
	}
	r, _ := strings.CutSuffix(bits, ".")
	return r[len(r)-size:]
}

func MergeCubes(a, b *Cube) (*Cube, error) {
	distance := a.DifferenceCardinality(b.BitSet)
	if distance > 1 {
		return nil, errors.New("cubes are too distant")
	}

	return &Cube{
		BitSet: a.Union(b.BitSet),
	}, nil
}
