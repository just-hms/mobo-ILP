package mobo_test

import (
	"testing"

	"github.com/just-hms/mobo/pkg/mobo"
	"github.com/stretchr/testify/require"
)

func FuzzMobo(f *testing.F) {
	f.Setenv("CPLEX_PATH", "/opt/ibm/ILOG/CPLEX_Studio2211/cplex/bin/x86-64_linux/cplex")

	f.Add(0)
	f.Fuzz(func(t *testing.T, seed int) {
		outs := mobo.RandomOutputs(seed)
		ports, _, _ := mobo.Solve(outs)

		err := mobo.Assert(outs, ports)
		require.NoError(t, err)
	})
}
