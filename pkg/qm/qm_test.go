package qm_test

import (
	"slices"
	"testing"

	"github.com/just-hms/mobo/pkg/qm"
	"github.com/just-hms/mobo/pkg/qm/cube"
	"github.com/stretchr/testify/require"
)

func TestGetCubes(t *testing.T) {
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
			got := qm.GetCubes(tt.input)

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
