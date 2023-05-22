package datafactory_test

import (
	"testing"

	keepertest "datafactory/testutil/keeper"
	"datafactory/testutil/nullify"
	"datafactory/x/datafactory"
	"datafactory/x/datafactory/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.DatafactoryKeeper(t)
	datafactory.InitGenesis(ctx, *k, genesisState)
	got := datafactory.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}
