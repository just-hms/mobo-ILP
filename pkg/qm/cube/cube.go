package cube

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/bits-and-blooms/bitset"
)

type Cube struct {
	val   *bitset.BitSet
	minus *bitset.BitSet
}

func New(val uint) *Cube {
	return FromString(strconv.FormatInt(int64(val), 2))
}

// FromString generates a cube from a string
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

// Clone returns a clone of the current cube
func (c *Cube) Clone() *Cube {
	return &Cube{
		val:   c.val.Clone(),
		minus: c.minus.Clone(),
	}
}

// Ones retuns the number of ones in the cube
func (c *Cube) Ones() uint {
	return c.val.DifferenceCardinality(c.minus)
}

// Repr returns the representation of the cube in a specified number of bits
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

// Len returns the minumun number of bits to represent the cube
func (c *Cube) Len() int {
	return int(max(c.val.Len(), c.minus.Len()))
}

// String retuns the representation of the cube in the minimum amount of bits
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

// Equal returns true if a cube is equal to the provided one
func (c *Cube) Equal(b *Cube) bool {
	minusDiff := c.minus.SymmetricDifferenceCardinality(b.minus)
	if minusDiff != 0 {
		return false
	}
	allC := c.minus.Union(c.val)
	allB := b.minus.Union(b.val)
	return allC.SymmetricDifferenceCardinality(allB) == 0
}

// Covers returns true if the cube cover a minterm in the circuit
func (c *Cube) Covers(one uint) bool {

	input := New(one).val

	maskC := c.val.Union(c.minus)
	maskInput := input.Union(c.minus)

	return maskC.Equal(maskInput)
}

// Display prints the cube using x_1*!x_2 format
func (c *Cube) Display(size uint) string {

	s, err := c.Repr(size)
	if err != nil {
		return "Error: " + err.Error()
	}

	rev := []rune(s)
	slices.Reverse(rev)

	builder := strings.Builder{}

	sep := ""
	for i := 0; i < len(rev); i++ {
		if rev[i] == '-' {
			continue
		}

		builder.WriteString(sep)

		if rev[i] == '0' {
			builder.WriteString("!")
		}

		builder.WriteString(fmt.Sprintf("x_%d", i+1))
		sep = "*"
	}

	return builder.String()
}

// Merge merges two cubes using the QM rules, returns an error if they are equal or they are too distant
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
