package qm_test

import (
	"testing"

	"github.com/just-hms/mobo/pkg/qm"
	"github.com/stretchr/testify/require"
)

func TestMergeCubes(t *testing.T) {
	t.Parallel()
	req := require.New(t)

	testcases := []struct {
		name   string
		ones   [2]uint
		exp    string
		expErr bool
	}{
		{
			name:   "0000 x 0000",
			ones:   [2]uint{0, 0},
			exp:    "0000",
			expErr: false,
		},
		{
			name:   "0000 x 0001",
			ones:   [2]uint{0, 1},
			exp:    "0001",
			expErr: false,
		},
		{
			name:   "0001 x 0100",
			ones:   [2]uint{1, 4},
			exp:    "",
			expErr: true,
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			a := qm.NewCube(tt.ones[0])
			b := qm.NewCube(tt.ones[1])

			res, err := qm.MergeCubes(a, b)
			if tt.expErr {
				req.Error(err)
			} else {
				req.NoError(err)
			}

			if !tt.expErr {
				req.Equal(tt.exp, res.Repr(len(tt.exp)))
			}
		})
	}
}
