package cmd

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/spf13/cobra"

	"polymarket-cli/internal/config"
	"polymarket-cli/pkg/relayer"
	"polymarket-cli/pkg/relayer/transactions"
)

var (
	privateKey string
	txType     string
)

var redeemCmd = &cobra.Command{
	Use:   "redeem [condition-id]",
	Short: "Redeem positions for a condition",
	Long:  `Redeem positions for a given condition ID.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Error: condition ID is required")
			return
		}

		conditionID, err := hexutil.Decode(args[0])
		if err != nil {
			fmt.Printf("Error: invalid condition ID: %v\n", err)
			return
		}

		if len(privateKey) == 0 {
			fmt.Println("Error: private key is required")
			return
		}

		if len(config.AppCfg.Builder.APIKey) == 0 {
			fmt.Println("Error: builder API key not configured")
			return
		}

		result, err := executeRedeem(conditionID)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		jsonData, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			fmt.Printf("Error formatting output: %v\n", err)
			return
		}

		fmt.Println(string(jsonData))
	},
}

func init() {
	rootCmd.AddCommand(redeemCmd)

	redeemCmd.Flags().StringVar(&privateKey, "private-key", "", "Private key for signing (required)")
	redeemCmd.Flags().StringVar(&txType, "tx-type", "SAFE", "Transaction type (SAFE or PROXY)")
	redeemCmd.MarkFlagRequired("private-key")
}

func executeRedeem(conditionID []byte) (*relayer.ExecuteResponse, error) {
	creds := &relayer.BuilderCreds{
		Key:        config.AppCfg.Builder.APIKey,
		Secret:     config.AppCfg.Builder.APISecret,
		Passphrase: config.AppCfg.Builder.Passphrase,
	}

	relayerTxType := relayer.RelayerTxTypeSAFE
	if txType == "PROXY" {
		relayerTxType = relayer.RelayerTxTypePROXY
	}

	client, err := relayer.NewClient(creds, relayerTxType, nil, &privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create relayer client: %w", err)
	}

	params := transactions.RedeemParams{
		ConditionalTokens:  relayer.CTF_ADDRESS,
		CollateralToken:    relayer.USDC_ADDRESS,
		ParentCollectionID: common.Hash{},
		ConditionID:        common.BytesToHash(conditionID),
		IndexSets: []*big.Int{
			big.NewInt(1),
			big.NewInt(2),
		},
	}

	tx, err := transactions.BuildRedeemTransaction(params)
	if err != nil {
		return nil, fmt.Errorf("failed to build transaction: %w", err)
	}

	return client.Execute([]*transactions.Transaction{tx}, "Redeem positions")
}
