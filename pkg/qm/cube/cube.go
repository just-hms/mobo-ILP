package cube

import (
	"errors"
	"fmt"
	"slices"
	"strconv"

	"github.com/bits-and-blooms/bitset"
)

type Cube struct {
	val   *bitset.BitSet
	minus *bitset.BitSet
}

func New(val uint) *Cube {
	return FromString(strconv.FormatInt(int64(val), 2))
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

	diff := int(size) - len(res)
	if diff < 0 {
		return "", fmt.Errorf("cannot represent %q in %d bits", res, size)
	}

	padding := make([]rune, diff)
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

func (c *Cube) Covers(one uint) bool {
	cOnes := c.val.Difference(c.minus)

	remainingOnes := New(one).val.Difference(c.minus)
	if cOnes.Count() == 0 && remainingOnes.Count() == 0 {
		return true
	}
	return cOnes.Equal(remainingOnes)
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
		val:   a.val.Difference(i),
		minus: a.minus.Union(i),
	}, nil
}
