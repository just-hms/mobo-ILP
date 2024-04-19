package qm_test

import (
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
		},
	}

	for _, tt := range testcases {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := qm.GetCubes(tt.input)

			req.Equal(len(tt.exp), len(got), tt.name)
			for i := range tt.exp {
				req.Equal(tt.exp[i].String(), got[i].String(), tt.name)
			}
		})
	}
}
