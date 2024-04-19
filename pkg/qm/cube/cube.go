package cube

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

func New(val uint) *Cube {
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

func FromString(val string) *Cube {
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
		case '0':
			continue
		default:
			panic(fmt.Sprintf("value %q not known", r))
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

func (c *Cube) Repr(size uint) (string, error) {
	res := c.String()

	diff := len(res) - int(size)
	if diff < 0 {
		return "", fmt.Errorf("cannot represent %v in %d bits", res, size)
	}

	padding := make([]rune, 0, diff)
	for i := range diff {
		padding[i] = '0'
	}

	return string(padding) + res, nil

}

func (c *Cube) Len() int {
	return int(max(c.val.Len(), c.minus.Len()))
}

func (c *Cube) String() string {
	len := c.Len()
	res := make([]rune, 0, len)
	for i := range uint(len) {
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
	minusDiff := c.minus.SymmetricDifferenceCardinality(b.minus)
	if minusDiff != 0 {
		return false
	}
	allC := c.minus.Union(c.val)
	allB := b.minus.Union(b.val)
	return allC.SymmetricDifferenceCardinality(allB) == 0
}

func Merge(a, b *Cube) (*Cube, error) {
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

	return &Cube{
		val:   a.val.Union(i),
		minus: a.minus.Union(i),
	}, nil
}
