package keeper

import (
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/bianjieai/irita/modules/service/internal/types"
)

func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case types.QueryDefinition:
			return queryDefinition(ctx, req, k)
		case types.QueryBinding:
			return queryBinding(ctx, req, k)
		case types.QueryBindings:
			return queryBindings(ctx, req, k)
		case types.QueryRequests:
			return queryRequests(ctx, req, k)
		case types.QueryResponse:
			return queryResponse(ctx, req, k)
		case types.QueryFees:
			return queryFees(ctx, req, k)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown query path: %s", path[0])
		}
	}
}

func queryDefinition(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	var params types.QueryDefinitionParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	svcDef, found := keeper.GetServiceDefinition(ctx, params.DefChainID, params.ServiceName)
	if !found {
		return nil, types.ErrUnknownSvcDef
	}

	iterator := keeper.GetMethods(ctx, params.DefChainID, params.ServiceName)
	defer iterator.Close()

	var methods []types.MethodProperty
	for ; iterator.Valid(); iterator.Next() {
		var method types.MethodProperty
		keeper.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &method)
		methods = append(methods, method)
	}

	definitionOutput := types.DefinitionOutput{Definition: svcDef, Methods: methods}

	bz, err := codec.MarshalJSONIndent(keeper.cdc, definitionOutput)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryBinding(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	var params types.QueryBindingParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	svcBinding, found := keeper.GetServiceBinding(ctx, params.DefChainID, params.ServiceName, params.BindChainId, params.Provider)
	if !found {
		return nil, types.ErrUnknownSvcBinding
	}

	bz, err := codec.MarshalJSONIndent(keeper.cdc, svcBinding)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryBindings(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	var params types.QueryBindingsParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	iterator := keeper.ServiceBindingsIterator(ctx, params.DefChainID, params.ServiceName)
	defer iterator.Close()

	var bindings []types.SvcBinding
	for ; iterator.Valid(); iterator.Next() {
		var binding types.SvcBinding
		keeper.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &binding)
		bindings = append(bindings, binding)
	}

	bz, err := codec.MarshalJSONIndent(keeper.cdc, bindings)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryRequests(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	var params types.QueryRequestsParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	iterator := keeper.ActiveBindRequestsIterator(ctx, params.DefChainID, params.ServiceName, params.BindChainId, params.Provider)
	defer iterator.Close()

	var requests []types.SvcRequest
	for ; iterator.Valid(); iterator.Next() {
		var request types.SvcRequest
		keeper.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &request)
		requests = append(requests, request)
	}

	bz, err := codec.MarshalJSONIndent(keeper.cdc, requests)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryResponse(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	var params types.QueryResponseParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	eHeight, rHeight, counter, err := types.ConvertRequestID(params.RequestId)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid request id: %s", params.RequestId)
	}

	response, found := keeper.GetResponse(ctx, params.ReqChainId, eHeight, rHeight, counter)
	if !found {
		return nil, types.ErrUnknownResponse
	}

	bz, err := codec.MarshalJSONIndent(keeper.cdc, response)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryFees(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	var params types.QueryFeesParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	var feesOutput types.FeesOutput

	if returnFee, found := keeper.GetReturnFee(ctx, params.Address); found {
		feesOutput.ReturnedFee = returnFee.Coins
	}

	if incomingFee, found := keeper.GetIncomingFee(ctx, params.Address); found {
		feesOutput.IncomingFee = incomingFee.Coins
	}

	bz, err := codec.MarshalJSONIndent(keeper.cdc, feesOutput)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	return bz, nil
}
