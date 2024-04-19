package qm

import (
	"errors"
	"fmt"
	"slices"

	"github.com/bits-and-blooms/bitset"
)

type Cube struct {
	val   *bitset.BitSet
	minus *bitset.BitSet
}

func CubeFromValue(val uint) *Cube {
	c := &Cube{
		val:   &bitset.BitSet{},
		minus: &bitset.BitSet{},
	}
	if val == 0 {
		return c
	}
	c.val.Set(val - 1)
	return c
}

func CubeFromString(val string) *Cube {
	c := &Cube{
		val:   &bitset.BitSet{},
		minus: &bitset.BitSet{},
	}

	rev := []rune(val)
	slices.Reverse(rev)

	for i, r := range rev {
		switch r {
		case '1':
			c.val.Set(uint(i))
		case '-':
			c.minus.Set(uint(i))
		}
	}
	return c
}

func (c *Cube) Clone() *Cube {
	return &Cube{
		val:   c.val.Clone(),
		minus: c.minus.Clone(),
	}
}

func (c *Cube) Ones() uint {
	return c.val.DifferenceCardinality(c.minus)
}

func (c *Cube) Repr(size uint) string {
	res := make([]rune, 0, size)
	for i := range size {
		toApp := '0'
		switch {
		case c.minus.Test(i):
			toApp = '-'
		case c.val.Test(i):
			toApp = '1'
		}

		res = append(res, toApp)
	}
	slices.Reverse(res)
	return string(res)
}

func MergeCubes(a, b *Cube) (*Cube, error) {
	repra := a.Repr(10)
	reprb := b.Repr(10)
	fmt.Println(repra, reprb)

	minusDiff := a.minus.SymmetricDifferenceCardinality(b.minus)
	if minusDiff != 0 {
		return nil, errors.New("cannot merge different number of -")
	}

	i := a.val.SymmetricDifference(b.val)
	if i.Count() != 1 {
		return nil, errors.New("cannot merge too many ones and zero differences")
	}

	return &Cube{
		val:   a.val.Union(i),
		minus: a.minus.Union(i),
	}, nil
}
