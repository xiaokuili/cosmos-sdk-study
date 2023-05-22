package keeper_test

import (
	"context"
	"testing"

	keepertest "datafactory/testutil/keeper"
	"datafactory/x/datafactory/keeper"
	"datafactory/x/datafactory/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func setupMsgServer(t testing.TB) (types.MsgServer, context.Context) {
	k, ctx := keepertest.DatafactoryKeeper(t)
	return keeper.NewMsgServerImpl(*k), sdk.WrapSDKContext(ctx)
}
