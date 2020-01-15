package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bianjieai/irita/modules/service/internal/types"
)

func GetCmdQuerySvcDef(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "definition",
		Short:   "Query service definition",
		Example: "iritacli query service definition <def-chain-id> <service name>",
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			params := types.QueryDefinitionParams{
				DefChainID:  args[0],
				ServiceName: args[1],
			}

			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryDefinition)
			res, _, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			var svcDefOutput types.DefinitionOutput
			err = cdc.UnmarshalJSON(res, &svcDefOutput)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(svcDefOutput)
		},
	}

	return cmd
}

func GetCmdQuerySvcBind(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "binding",
		Short:   "Query service binding",
		Example: "iritacli query service binding <def-chain-id> <service name> <bind-chain-id> <provider>",
		Args:    cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			provider, err := sdk.AccAddressFromBech32(args[3])
			if err != nil {
				return err
			}

			params := types.QueryBindingParams{
				DefChainID:  args[0],
				ServiceName: args[1],
				BindChainId: args[2],
				Provider:    provider,
			}

			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryBinding)
			res, _, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			var binding types.SvcBinding
			err = cdc.UnmarshalJSON(res, &binding)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(binding)
		},
	}

	return cmd
}

func GetCmdQuerySvcBinds(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "bindings",
		Short:   "Query service bindings",
		Example: "iritacli query service bindings <def-chain-id> <service name>",
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			params := types.QueryBindingsParams{
				DefChainID:  args[0],
				ServiceName: args[1],
			}

			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryBindings)
			res, _, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			var bindings []types.SvcBinding
			err = cdc.UnmarshalJSON(res, &bindings)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(bindings)
		},
	}

	return cmd
}

func GetCmdQuerySvcRequests(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "requests",
		Short:   "Query service requests",
		Example: "iritacli query service requests <def-chain-id> <service-name> <bind-chain-id> <provider>",
		Args:    cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			provider, err := sdk.AccAddressFromBech32(args[3])
			if err != nil {
				return err
			}

			params := types.QueryBindingParams{
				DefChainID:  args[0],
				ServiceName: args[1],
				BindChainId: args[2],
				Provider:    provider,
			}

			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryRequests)
			res, _, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			var requests []types.SvcRequest
			err = cdc.UnmarshalJSON(res, &requests)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(requests)
		},
	}

	return cmd
}

func GetCmdQuerySvcResponse(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "response",
		Short:   "Query a service response",
		Example: "iritacli query service response <req-chain-id> <request-id>",
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			params := types.QueryResponseParams{
				ReqChainId: args[0],
				RequestId:  args[1],
			}

			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryResponse)
			res, _, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			var response types.SvcResponse
			err = cdc.UnmarshalJSON(res, &response)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(response)
		},
	}

	return cmd
}

func GetCmdQuerySvcFees(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "fees",
		Short:   "Query returned and incoming fee of a particular address",
		Example: "iritacli query service fees <account address>",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			addr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			params := types.QueryFeesParams{
				Address: addr,
			}

			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryFees)
			res, _, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			var fees types.FeesOutput
			err = cdc.UnmarshalJSON(res, &fees)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(fees)
		},
	}

	return cmd
}
