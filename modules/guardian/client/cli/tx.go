package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/bianjieai/irita/modules/guardian/internal/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "guardian transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	txCmd.AddCommand(client.PostCommands(
		GetCmdCreateProfiler(cdc),
		GetCmdDeleteProfiler(cdc),
		GetCmdCreateTrustee(cdc),
		GetCmdDeleteTrustee(cdc),
	)...)
	return txCmd
}

func GetCmdCreateProfiler(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-profiler",
		Short: "Add a new profiler",
		Example: "iritacli tx guardian add-profiler --chain-id=<chain-id> --from=<key-name> --fees=0.3iris " +
			"--address=<added address> --description=<name>",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			fromAddr := cliCtx.GetFromAddress()

			paStr := viper.GetString(FlagAddress)
			if len(paStr) == 0 {
				return fmt.Errorf("must use --address flag")
			}
			pAddr, err := sdk.AccAddressFromBech32(paStr)
			if err != nil {
				return err
			}
			description := viper.GetString(FlagDescription)
			if len(description) == 0 {
				return fmt.Errorf("must use --description flag")
			}
			msg := types.NewMsgAddProfiler(description, pAddr, fromAddr)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	cmd.Flags().AddFlagSet(FsAddGuardian)
	cmd.MarkFlagRequired(FlagAddress)
	cmd.MarkFlagRequired(FlagDescription)
	return cmd
}

func GetCmdDeleteProfiler(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-profiler",
		Short: "Delete a profiler",
		Example: "iritacli tx guardian delete-profiler --chain-id=<chain-id> --from=<key-name> --fees=0.3iris " +
			"--address=<deleted address>",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			fromAddr := cliCtx.GetFromAddress()
			paStr := viper.GetString(FlagAddress)
			if len(paStr) == 0 {
				return fmt.Errorf("must use --address flag")
			}
			pAddr, err := sdk.AccAddressFromBech32(paStr)
			if err != nil {
				return err
			}
			msg := types.NewMsgDeleteProfiler(pAddr, fromAddr)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	cmd.Flags().AddFlagSet(FsDeleteGuardian)
	cmd.MarkFlagRequired(FlagAddress)
	return cmd
}

func GetCmdCreateTrustee(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-trustee",
		Short: "Add a new trustee",
		Example: "iritacli tx guardian add-trustee --chain-id=<chain-id> --from=<key-name> --fees=0.3iris " +
			"--address=<added address> --description=<name>",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			fromAddr := cliCtx.GetFromAddress()
			taStr := viper.GetString(FlagAddress)
			if len(taStr) == 0 {
				return fmt.Errorf("must use --address flag")
			}
			tAddr, err := sdk.AccAddressFromBech32(taStr)
			if err != nil {
				return err
			}
			description := viper.GetString(FlagDescription)
			if len(description) == 0 {
				return fmt.Errorf("must use --description flag")
			}
			msg := types.NewMsgAddTrustee(description, tAddr, fromAddr)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	cmd.Flags().AddFlagSet(FsAddGuardian)
	return cmd
}

func GetCmdDeleteTrustee(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-trustee",
		Short: "Delete a trustee",
		Example: "iritacli tx guardian delete-trustee --chain-id=<chain-id> --from=<key-name> --fees=0.3iris " +
			"--address=<deleted address>",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			fromAddr := cliCtx.GetFromAddress()
			taStr := viper.GetString(FlagAddress)
			if len(taStr) == 0 {
				return fmt.Errorf("must use --address flag")
			}
			tAddr, err := sdk.AccAddressFromBech32(taStr)
			if err != nil {
				return err
			}
			msg := types.NewMsgDeleteTrustee(tAddr, fromAddr)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	cmd.Flags().AddFlagSet(FsDeleteGuardian)
	return cmd
}
