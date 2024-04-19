package qm_test

import (
	"fmt"
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
		exp    *qm.Cube
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
			exp:  qm.CubeFromString("000-"),
		},
		{
			name: "Distant",
			a:    qm.CubeFromValue(0),
			b:    qm.CubeFromValue(2),
			exp:  qm.CubeFromString("00-0"),
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
			exp:  qm.CubeFromString("-0-0"),
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
			exp:  qm.CubeFromString("-0--"),
		},
	}

	for _, tt := range testcases {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := qm.MergeCubes(tt.a, tt.b)
			if tt.expErr {
				req.Error(err, fmt.Sprintf("test: %q a: %v b: %v", tt.name, tt.a, tt.b))
				return
			}

			req.NoError(err, fmt.Sprintf("test: %q a: %v b: %v", tt.name, tt.a, tt.b))
			req.True(tt.exp.Equal(got), fmt.Sprintf("test: %q exp: %v got: %v", tt.name, tt.exp, got))
		})
	}
}

func TestString(t *testing.T) {
	t.Parallel()
	req := require.New(t)

	testcases := []struct {
		name string
		a    *qm.Cube
		exp  string
	}{
		{
			name: "cutprefix",
			a:    qm.CubeFromString("0-00"),
			exp:  "-00",
		},
		{
			name: "minus",
			a:    qm.CubeFromString("1-00"),
			exp:  "1-00",
		},
		{
			name: "simple",
			a:    qm.CubeFromString("0101"),
			exp:  "101",
		},
		{
			name: "long",
			a:    qm.CubeFromString("10001000000-000000101"),
			exp:  "10001000000-000000101",
		},
	}

	for _, tt := range testcases {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := tt.a.String()

			req.Equal(tt.exp, got, tt.name)
		})
	}
}
