package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	errorlib "github.com/dydxprotocol/v4-chain/protocol/lib/error"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

// UpdateBlockRateLimitConfiguration updates the block rate limit configuration returning an error
// if the configuration is invalid.
func (k msgServer) UpdateBlockRateLimitConfiguration(
	goCtx context.Context,
	msg *types.MsgUpdateBlockRateLimitConfiguration,
) (resp *types.MsgUpdateBlockRateLimitConfigurationResponse, err error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	defer func() {
		metrics.IncrSuccessOrErrorCounter(err, types.ModuleName, metrics.UpdateBlockRateLimitConfiguration, metrics.DeliverTx)
		if err != nil {
			errorlib.LogDeliverTxError(k.Keeper.Logger(ctx), err, ctx.BlockHeight(), "UpdateBlockRateLimitConfiguration", msg)
		}
	}()

	if !k.Keeper.HasAuthority(msg.Authority) {
		return nil, errorsmod.Wrapf(
			govtypes.ErrInvalidSigner,
			"invalid authority %s",
			msg.Authority,
		)
	}

	if err := k.Keeper.InitializeBlockRateLimit(ctx, msg.BlockRateLimitConfig); err != nil {
		return nil, err
	}
	return &types.MsgUpdateBlockRateLimitConfigurationResponse{}, nil
}
