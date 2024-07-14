package bin_test

import (
	"testing"

	"github.com/just-hms/mobo/pkg/bin"
	"github.com/stretchr/testify/require"
)

func TestNextPowerOf2(t *testing.T) {
	t.Parallel()
	req := require.New(t)

	tests := []struct {
		input    uint
		expected uint
	}{
		{0, 1},
		{1, 1},
		{2, 2},
		{3, 4},
		{4, 4},
		{5, 8},
		{6, 8},
		{7, 8},
		{8, 8},
		{9, 16},
		{15, 16},
		{16, 16},
		{17, 32},
		{31, 32},
		{32, 32},
		{33, 64},
		{63, 64},
		{64, 64},
		{65, 128},
		{127, 128},
		{128, 128},
		{129, 256},
	}

	for _, tt := range tests {
		result := bin.NextPowerOf2(tt.input)
		req.Equal(tt.expected, result)
	}
}

// TestMinBitsNeeded tests the MinBitsNeeded function.
func TestMinBitsNeeded(t *testing.T) {
	t.Parallel()
	req := require.New(t)

	tests := []struct {
		input uint
		exp   uint
	}{
		{0, 1},
		{1, 1},
		{2, 2},
		{3, 2},
		{4, 3},
		{10, 4},
		{255, 8},
		{256, 9},
		{1023, 10},
		{1024, 11},
	}

	for _, tt := range tests {
		got := bin.MinBitsNeeded(tt.input)
		req.Equal(tt.exp, got)
	}
}
