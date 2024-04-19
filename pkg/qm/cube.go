package qm

import (
	"errors"
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

func (c *Cube) Equal(b *Cube) bool {
	allC := c.minus.Union(c.val)
	allB := b.minus.Union(b.val)
	return allC.SymmetricDifferenceCardinality(allB) == 0
}

func MergeCubes(a, b *Cube) (*Cube, error) {
	if a.Equal(b) {
		return nil, errors.New("cannot merge they are equal")
	}

	minusDiff := a.minus.SymmetricDifferenceCardinality(b.minus)
	if minusDiff != 0 {
		return nil, errors.New("cannot merge different number of -")
	}

	allA := a.minus.Union(a.val)
	allB := b.minus.Union(b.val)

	i := allA.SymmetricDifference(allB)
	if i.Count() != 1 {
		return nil, errors.New("cannot merge too many ones and zero differences")
	}

	// TODO: decide what to do with val
	return &Cube{
		val:   a.val.Union(i),
		minus: a.minus.Union(i),
	}, nil
}
