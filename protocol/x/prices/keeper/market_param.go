package keeper

import (
	"sort"

	errorsmod "cosmossdk.io/errors"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/metrics"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	indexerevents "github.com/dydxprotocol/v4-chain/protocol/indexer/events"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
)

// newMarketParamStore creates a new prefix store for MarketParams.
func (k Keeper) newMarketParamStore(ctx sdk.Context) prefix.Store {
	return prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.MarketParamKeyPrefix))
}

// ModifyMarketParam modifies an existing market param in the store.
func (k Keeper) ModifyMarketParam(
	ctx sdk.Context,
	marketParam types.MarketParam,
) (types.MarketParam, error) {
	// Validate input.
	if err := marketParam.Validate(); err != nil {
		return types.MarketParam{}, err
	}

	// Get existing market param.
	existingParam, exists := k.GetMarketParam(ctx, marketParam.Id)
	if !exists {
		return types.MarketParam{}, errorsmod.Wrap(
			types.ErrMarketParamDoesNotExist,
			lib.UintToString(marketParam.Id),
		)
	}

	// Validate update is permitted.
	if marketParam.Exponent != existingParam.Exponent {
		return types.MarketParam{},
			errorsmod.Wrapf(types.ErrMarketExponentCannotBeUpdated, lib.UintToString(marketParam.Id))
	}

	// Store the modified market param.
	marketParamStore := k.newMarketParamStore(ctx)
	b := k.cdc.MustMarshal(&marketParam)
	marketParamStore.Set(lib.Uint32ToKey(marketParam.Id), b)

	k.GetIndexerEventManager().AddTxnEvent(
		ctx,
		indexerevents.SubtypeMarket,
		indexerevents.MarketEventVersion,
		indexer_manager.GetBytes(
			indexerevents.NewMarketModifyEvent(
				marketParam.Id,
				marketParam.Pair,
				marketParam.MinPriceChangePpm,
			),
		),
	)

	// Update the in-memory market pair map for labelling metrics.
	metrics.AddMarketPairForTelemetry(marketParam.Id, marketParam.Pair)

	return marketParam, nil
}

// GetMarketParam returns a market param from its id.
func (k Keeper) GetMarketParam(
	ctx sdk.Context,
	id uint32,
) (
	market types.MarketParam,
	exists bool,
) {
	marketParamStore := k.newMarketParamStore(ctx)
	b := marketParamStore.Get(lib.Uint32ToKey(id))
	if b == nil {
		return types.MarketParam{}, false
	}

	k.cdc.MustUnmarshal(b, &market)
	return market, true
}

// GetAllMarketParams returns all market params.
func (k Keeper) GetAllMarketParams(ctx sdk.Context) []types.MarketParam {
	marketParamStore := k.newMarketParamStore(ctx)

	marketParams := make([]types.MarketParam, 0)

	iterator := marketParamStore.Iterator(nil, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		marketParam := types.MarketParam{}
		k.cdc.MustUnmarshal(iterator.Value(), &marketParam)
		marketParams = append(marketParams, marketParam)
	}

	// Sort the market params to return them in ascending order based on Id.
	sort.Slice(marketParams, func(i, j int) bool {
		return marketParams[i].Id < marketParams[j].Id
	})

	return marketParams
}
