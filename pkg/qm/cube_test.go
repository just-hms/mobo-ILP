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
		a, b   *qm.Cube
		exp    string
		expErr bool
	}{
		{
			name:   "Identical Cubes",
			a:      qm.CubeFromString("010-"),
			b:      qm.CubeFromString("010-"),
			expErr: true,
		},
		{
			name: "Adjacent Cubes",
			a:    qm.CubeFromValue(0),
			b:    qm.CubeFromValue(1),
			exp:  "000-",
		},
		{
			name: "Distant",
			a:    qm.CubeFromValue(0),
			b:    qm.CubeFromValue(2),
			exp:  "00-0",
		},
		{
			name:   "Wrong",
			a:      qm.CubeFromValue(1),
			b:      qm.CubeFromValue(4),
			expErr: true,
		},
		{
			name: "Ok with minus",
			a:    qm.CubeFromString("00-0"),
			b:    qm.CubeFromString("10-0"),
			exp:  "-0-0",
		},
		{
			name:   "Wrong with minus",
			a:      qm.CubeFromString("00-1"),
			b:      qm.CubeFromString("10-0"),
			expErr: true,
		},
		{
			name: "Multiple Minus Signs",
			a:    qm.CubeFromString("-0-1"),
			b:    qm.CubeFromString("-0-0"),
			exp:  "-0--",
		},
	}

	for _, tt := range testcases {

		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			res, err := qm.MergeCubes(tt.a, tt.b)
			if tt.expErr {
				req.Error(err, tt.name)
				return
			}

			req.NoError(err, tt.name)
			size := uint(len(tt.exp))
			req.Equal(tt.exp, res.Repr(size), tt.name)
		})
	}
}
