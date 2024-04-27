package qm_test

import (
	"slices"
	"strconv"
	"testing"

	"github.com/just-hms/mobo/pkg/bin"
	"github.com/just-hms/mobo/pkg/qm"
	"github.com/just-hms/mobo/pkg/qm/cube"
	"github.com/stretchr/testify/require"
)

func TestCubes(t *testing.T) {
	t.Parallel()
	req := require.New(t)

	testcases := []struct {
		name  string
		input []*cube.Cube
		exp   []*cube.Cube
		size  int
	}{
		{
			name: "Simple",
			input: []*cube.Cube{
				cube.FromString("0000"),
				cube.FromString("0001"),
				cube.FromString("0101"),
				cube.FromString("1101"),
			},
			exp: []*cube.Cube{
				cube.FromString("0000"),
				cube.FromString("0001"),
				cube.FromString("0101"),
				cube.FromString("1101"),

				cube.FromString("000-"),
				cube.FromString("0-01"),
				cube.FromString("-101"),
			},
			size: 4,
		},
		{
			name: "Example 4.1-y1",
			input: []*cube.Cube{
				cube.New(1), cube.New(5),
				cube.New(2), cube.New(3),
			},
			exp: []*cube.Cube{
				cube.FromString("001"),
				cube.FromString("101"),
				cube.FromString("010"),
				cube.FromString("011"),

				cube.FromString("-01"),
				cube.FromString("0-1"),
				cube.FromString("01-"),
			},
			size: 3,
		},
		{
			name: "Example 4.1-y2",
			input: []*cube.Cube{
				cube.New(1), cube.New(5),
				cube.New(7), cube.New(6),
			},
			exp: []*cube.Cube{
				cube.FromString("001"),
				cube.FromString("101"),
				cube.FromString("111"),
				cube.FromString("110"),

				cube.FromString("-01"),
				cube.FromString("1-1"),
				cube.FromString("11-"),
			},
			size: 3,
		},
	}

	for _, tt := range testcases {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			// t.Parallel()
			got := qm.Cubes(tt.input)

			gotDump := []string{}
			for _, g := range got {
				repr, err := g.Repr(uint(tt.size))
				req.NoError(err)
				gotDump = append(gotDump, repr)
			}
			slices.Sort(gotDump)

			expDump := []string{}
			for _, g := range tt.exp {
				repr, err := g.Repr(uint(tt.size))
				req.NoError(err)
				expDump = append(expDump, repr)
			}
			slices.Sort(expDump)

			req.Equal(expDump, gotDump)

		})
	}
}

func FuzzCubes(f *testing.F) {
	req := require.New(f)

	f.Add(0)
	f.Fuzz(func(t *testing.T, seed int) {
		ones := qm.RandomOnes(seed)

		input := make([]*cube.Cube, 0, len(ones))
		for _, o := range ones {
			input = append(input, cube.New(o))
		}

		cubes := qm.Cubes(input)
		_max := slices.MaxFunc(ones, func(a, b uint) int {
			return int(a) - int(b)
		})

		for j := range bin.NextPowerOf2(_max) {

			mustBeCovered := slices.Contains(ones, j)
			coverers := make([]*cube.Cube, 0)
			for _, c := range cubes {
				if c.Covers(j) {
					coverers = append(coverers, c)
					break
				}
			}

			binary := strconv.FormatInt(int64(j), 2)
			if mustBeCovered && len(coverers) == 0 {
				req.Failf("fail", "out: %d, %s should be covered but is not", j+1, binary)
			}

			if !mustBeCovered && len(coverers) > 0 {
				req.Failf("fail", "out: %d, %s shouldn't be covered but it is %v", j+1, binary, coverers)
			}
		}

	})

}
