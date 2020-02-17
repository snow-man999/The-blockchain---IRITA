package cli

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"

	"github.com/bianjieai/irita/modules/service/internal/types"
)

func GetCmdSvcDef(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "define",
		Short: "Create a new service definition",
		Example: "iritacli tx service define --chain-id=<chain-id> --from=<key-name> --fee=0.3iris " +
			"--service-name=<service name> --service-description=<service description> --author-description=<author description> " +
			"--tags=tag1,tag2 --idl-content=<interface description content> --file=test.proto",
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			name := viper.GetString(FlagServiceName)
			description := viper.GetString(FlagServiceDescription)
			authorDescription := viper.GetString(FlagAuthorDescription)
			tags := viper.GetStringSlice(FlagTags)

			content := viper.GetString(FlagIdlContent)
			if len(content) > 0 {
				content = strings.Replace(content, `\n`, "\n", -1)
			}

			filePath := viper.GetString(FlagFile)
			if len(filePath) > 0 {
				contentBytes, err := ioutil.ReadFile(filePath)
				if err != nil {
					return err
				}
				content = string(contentBytes)
			}

			fmt.Printf("idl condent: \n%s\n", content)

			chainID := viper.GetString(flags.FlagChainID)
			fromAddr := cliCtx.GetFromAddress()

			msg := types.NewMsgSvcDef(name, chainID, description, tags, fromAddr, authorDescription, content)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(FsServiceDefinitionCreate)
	cmd.MarkFlagRequired(FlagServiceName)

	return cmd
}

func GetCmdSvcBind(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bind",
		Short: "Create a new service binding",
		Example: "iritacli tx service bind --chain-id=<chain-id> --from=<key-name> --fee=0.3iris " +
			"--service-name=<service name> --def-chain-id=<chain-id> --bind-type=Local " +
			"--deposit=1iris --prices=1iris,2iris --avg-rsp-time=10000 --usable-time=100",
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			fromAddr := cliCtx.GetFromAddress()

			chainID := viper.GetString(flags.FlagChainID)
			name := viper.GetString(FlagServiceName)
			defChainID := viper.GetString(FlagDefChainID)
			initialDeposit := viper.GetString(FlagDeposit)
			initialPrices := viper.GetStringSlice(FlagPrices)
			avgRspTime := viper.GetInt64(FlagAvgRspTime)
			usableTime := viper.GetInt64(FlagUsableTime)
			bindingTypeStr := viper.GetString(FlagBindType)

			bindingType, err := types.BindingTypeFromString(bindingTypeStr)
			if err != nil {
				return err
			}

			deposit, err := sdk.ParseCoins(initialDeposit)
			if err != nil {
				return err
			}

			var prices []sdk.Coin
			for _, ip := range initialPrices {
				price, err := sdk.ParseCoin(ip)
				if err != nil {
					return err
				}
				prices = append(prices, price)
			}

			level := types.Level{AvgRspTime: avgRspTime, UsableTime: usableTime}

			msg := types.NewMsgSvcBind(defChainID, name, chainID, fromAddr, bindingType, deposit, prices, level)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(FsServiceDefinition)
	cmd.Flags().AddFlagSet(FsServiceBindingCreate)
	cmd.MarkFlagRequired(FlagDefChainID)
	cmd.MarkFlagRequired(FlagServiceName)
	cmd.MarkFlagRequired(FlagBindType)
	cmd.MarkFlagRequired(FlagPrices)
	cmd.MarkFlagRequired(FlagAvgRspTime)
	cmd.MarkFlagRequired(FlagUsableTime)

	return cmd
}

func GetCmdSvcBindUpdate(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-binding",
		Short: "Update a service binding",
		Example: "iritacli tx service update-binding --chain-id=<chain-id> --from=<key-name> --fee=0.3iris " +
			"--service-name=<service name> --def-chain-id=<chain-id> --bind-type=Local " +
			"--deposit=1iris --prices=1iris,2iris --avg-rsp-time=10000 --usable-time=100",
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			fromAddr := cliCtx.GetFromAddress()

			chainID := viper.GetString(flags.FlagChainID)
			name := viper.GetString(FlagServiceName)
			defChainID := viper.GetString(FlagDefChainID)
			initialDeposit := viper.GetString(FlagDeposit)
			initialPrices := viper.GetStringSlice(FlagPrices)
			avgRspTime := viper.GetInt64(FlagAvgRspTime)
			usableTime := viper.GetInt64(FlagUsableTime)
			bindingTypeStr := viper.GetString(FlagBindType)

			var err error

			var bindingType types.BindingType
			if bindingTypeStr != "" {
				bindingType, err = types.BindingTypeFromString(bindingTypeStr)
				if err != nil {
					return err
				}
			}

			var deposit sdk.Coins
			if initialDeposit != "" {
				deposit, err = sdk.ParseCoins(initialDeposit)
				if err != nil {
					return err
				}
			}

			var prices []sdk.Coin
			for _, ip := range initialPrices {
				price, err := sdk.ParseCoin(ip)
				if err != nil {
					return err
				}
				prices = append(prices, price)
			}

			level := types.Level{AvgRspTime: avgRspTime, UsableTime: usableTime}

			msg := types.NewMsgSvcBindingUpdate(defChainID, name, chainID, fromAddr, bindingType, deposit, prices, level)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(FsServiceDefinition)
	cmd.Flags().AddFlagSet(FsServiceBindingUpdate)
	cmd.MarkFlagRequired(FlagDefChainID)
	cmd.MarkFlagRequired(FlagServiceName)

	return cmd
}

func GetCmdSvcDisable(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "disable",
		Short: "Disable an available service binding",
		Example: "iritacli tx service disable --chain-id=<chain-id> --from=<key-name> --fee=0.3iris " +
			"--service-name=<service name> --def-chain-id=<chain-id>",
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			fromAddr := cliCtx.GetFromAddress()

			chainID := viper.GetString(flags.FlagChainID)
			name := viper.GetString(FlagServiceName)
			defChainID := viper.GetString(FlagDefChainID)

			msg := types.NewMsgSvcDisable(defChainID, name, chainID, fromAddr)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(FsServiceDefinition)
	cmd.MarkFlagRequired(FlagDefChainID)
	cmd.MarkFlagRequired(FlagServiceName)

	return cmd
}

func GetCmdSvcEnable(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "enable",
		Short: "Enable an unavailable service binding",
		Example: "iritacli tx service enable --chain-id=<chain-id> --from=<key-name> --fee=0.3iris " +
			"--service-name=<service name> --def-chain-id=<chain-id> --deposit=1iris",
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			fromAddr := cliCtx.GetFromAddress()

			chainID := viper.GetString(flags.FlagChainID)
			name := viper.GetString(FlagServiceName)
			defChainID := viper.GetString(FlagDefChainID)

			initialDeposit := viper.GetString(FlagDeposit)
			deposit, err := sdk.ParseCoins(initialDeposit)
			if err != nil {
				return err
			}

			msg := types.NewMsgSvcEnable(defChainID, name, chainID, fromAddr, deposit)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(FsServiceDefinition)
	cmd.Flags().String(FlagDeposit, "", "additional deposit of binding")
	cmd.MarkFlagRequired(FlagDefChainID)
	cmd.MarkFlagRequired(FlagServiceName)

	return cmd
}

func GetCmdSvcRefundDeposit(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "refund-deposit",
		Short: "Refund all deposit from a service binding",
		Example: "iritacli tx service refund-deposit --chain-id=<chain-id> --from=<key-name> --fee=0.3iris " +
			"--service-name=<service name> --def-chain-id=<chain-id>",
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			fromAddr := cliCtx.GetFromAddress()

			chainID := viper.GetString(flags.FlagChainID)
			name := viper.GetString(FlagServiceName)
			defChainID := viper.GetString(FlagDefChainID)

			msg := types.NewMsgSvcRefundDeposit(defChainID, name, chainID, fromAddr)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(FsServiceDefinition)
	cmd.MarkFlagRequired(FlagDefChainID)
	cmd.MarkFlagRequired(FlagServiceName)

	return cmd
}

func GetCmdSvcCall(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "call",
		Short: "Call a service method",
		Example: "iritacli tx service call --chain-id=<chain-id> --from=<key-name> --fee=0.3iris --def-chain-id=<bind-chain-id> " +
			"--service-name=<service name> --method-id=<method-id> --bind-chain-id=<chain-id> --provider=<provider> --service-fee=1iris --request-data=<req>",
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			fromAddr := cliCtx.GetFromAddress()

			chainID := viper.GetString(flags.FlagChainID)
			defChainID := viper.GetString(FlagDefChainID)
			name := viper.GetString(FlagServiceName)
			bindChainID := viper.GetString(FlagBindChainID)
			methodID := int16(viper.GetInt(FlagMethodID))

			providerStr := viper.GetString(FlagProvider)
			provider, err := sdk.AccAddressFromBech32(providerStr)
			if err != nil {
				return err
			}

			serviceFeeStr := viper.GetString(FlagServiceFee)
			serviceFee, err := sdk.ParseCoins(serviceFeeStr)
			if err != nil {
				return err
			}

			inputString := viper.GetString(FlagReqData)
			input, err := hex.DecodeString(inputString)
			if err != nil {
				return err
			}

			profiling := viper.GetBool(FlagProfiling)

			msg := types.NewMsgSvcRequest(defChainID, name, bindChainID, chainID, fromAddr, provider, methodID, input, serviceFee, profiling)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(FsServiceDefinition)
	cmd.Flags().AddFlagSet(FsServiceBinding)
	cmd.Flags().AddFlagSet(FsServiceRequest)
	cmd.MarkFlagRequired(FlagDefChainID)
	cmd.MarkFlagRequired(FlagServiceName)
	cmd.MarkFlagRequired(FlagBindChainID)
	cmd.MarkFlagRequired(FlagProvider)
	cmd.MarkFlagRequired(FlagMethodID)

	return cmd
}

func GetCmdSvcRespond(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "respond",
		Short: "Respond a service method invocation",
		Example: "iritacli service respond --chain-id=<chain-id> --from=<key-name> --fee=0.3iris --request-chain-id=<call-chain-id> " +
			"--request-id=<request-id> --response-data=<resp>",
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			fromAddr := cliCtx.GetFromAddress()

			reqChainID := viper.GetString(FlagReqChainID)

			outputString := viper.GetString(FlagRespData)
			output, err := hex.DecodeString(outputString)
			if err != nil {
				return err
			}

			errMsgString := viper.GetString(FlagErrMsg)
			errMsg, err := hex.DecodeString(errMsgString)
			if err != nil {
				return err
			}

			reqID := viper.GetString(FlagReqID)

			msg := types.NewMsgSvcResponse(reqChainID, reqID, fromAddr, output, errMsg)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(FsServiceResponse)
	cmd.MarkFlagRequired(FlagReqChainID)
	cmd.MarkFlagRequired(FlagReqID)

	return cmd
}

func GetCmdSvcRefundFees(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "refund-fees",
		Short:   "Refund all fees from service call timeout",
		Example: "iritacli tx service refund-fees --chain-id=<chain-id> --from=<key-name> --fee=0.3iris --dest-address=<account address> --withdraw-amount 1iris",
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			fromAddr := cliCtx.GetFromAddress()

			msg := types.NewMsgSvcRefundFees(fromAddr)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return cmd
}

func GetCmdSvcWithdrawFees(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "withdraw-fees",
		Short:   "withdraw all fees from service call reward",
		Example: "iritacli tx service withdraw-fees --chain-id=<chain-id> --from=<key-name> --fee=0.3iris",
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			fromAddr := cliCtx.GetFromAddress()

			msg := types.NewMsgSvcWithdrawFees(fromAddr)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return cmd
}

func GetCmdSvcWithdrawTax(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "withdraw-tax",
		Short:   "withdraw service fee tax to an account",
		Example: "iritacli tx service withdraw-tax --chain-id=<chain-id> --from=<key-name> --fee=0.3iris --dest-address=<account address> --withdraw-amount=1iris",
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			fromAddr := cliCtx.GetFromAddress()

			destAddressStr := viper.GetString(FlagDestAddress)
			destAddress, err := sdk.AccAddressFromBech32(destAddressStr)
			if err != nil {
				return err
			}

			withdrawAmountStr := viper.GetString(FlagWithdrawAmount)
			withdrawAmount, err := sdk.ParseCoins(withdrawAmountStr)
			if err != nil {
				return err
			}

			msg := types.NewMsgSvcWithdrawTax(fromAddr, destAddress, withdrawAmount)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(FsServiceWithdrawTax)
	cmd.MarkFlagRequired(FlagDestAddress)
	cmd.MarkFlagRequired(FlagWithdrawAmount)

	return cmd
}
