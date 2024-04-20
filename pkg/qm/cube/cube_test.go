package cube_test

import (
	"fmt"
	"testing"

	"github.com/just-hms/mobo/pkg/qm/cube"
	"github.com/stretchr/testify/require"
)

func TestMergeCubes(t *testing.T) {
	t.Parallel()
	req := require.New(t)

	testcases := []struct {
		name   string
		a, b   *cube.Cube
		exp    *cube.Cube
		expErr bool
	}{
		{
			name:   "Identical Cubes",
			a:      cube.FromString("010-"),
			b:      cube.FromString("010-"),
			expErr: true,
		},
		{
			name: "Adjacent Cubes",
			a:    cube.New(0),
			b:    cube.New(1),
			exp:  cube.FromString("000-"),
		},
		{
			name: "Distant",
			a:    cube.New(0),
			b:    cube.New(2),
			exp:  cube.FromString("00-0"),
		},
		{
			name:   "Wrong",
			a:      cube.New(1),
			b:      cube.New(4),
			expErr: true,
		},
		{
			name: "Ok with minus",
			a:    cube.FromString("00-0"),
			b:    cube.FromString("10-0"),
			exp:  cube.FromString("-0-0"),
		},
		{
			name:   "Wrong with minus",
			a:      cube.FromString("00-1"),
			b:      cube.FromString("10-0"),
			expErr: true,
		},
		{
			name: "Multiple Minus Signs",
			a:    cube.FromString("-0-1"),
			b:    cube.FromString("-0-0"),
			exp:  cube.FromString("-0--"),
		},
	}

	for _, tt := range testcases {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := cube.Merge(tt.a, tt.b)
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
		a    *cube.Cube
		exp  string
	}{
		{
			name: "cutprefix",
			a:    cube.FromString("0-00"),
			exp:  "-00",
		},
		{
			name: "minus",
			a:    cube.FromString("1-00"),
			exp:  "1-00",
		},
		{
			name: "simple",
			a:    cube.FromString("0101"),
			exp:  "101",
		},
		{
			name: "long",
			a:    cube.FromString("10001000000-000000101"),
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

func TestCovers(t *testing.T) {
	t.Parallel()
	req := require.New(t)

	testcases := []struct {
		name string
		a    *cube.Cube
		one  uint
		exp  bool
	}{
		{
			name: "Simple",
			a:    cube.FromString("0001"),
			one:  0b001,
			exp:  true,
		},
		{
			name: "Strange",
			a:    cube.FromString("0-0-"),
			one:  0b001,
			exp:  true,
		},
		{
			name: "Strange with one",
			a:    cube.FromString("0-01"),
			one:  0b001,
			exp:  true,
		},
		{
			name: "Strange with one, seems broken",
			a:    cube.FromString("-1"),
			one:  0b001,
			exp:  true,
		},
		{
			name: "Something else broken",
			a:    cube.FromString("-01"),
			one:  0b101,
			exp:  true,
		},
	}

	for _, tt := range testcases {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.a.Covers(tt.one)
			req.Equal(tt.exp, got, tt.name)
		})
	}

}
