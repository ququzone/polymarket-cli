package transactions

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

type RedeemParams struct {
	ConditionalTokens  common.Address
	CollateralToken    common.Address
	ParentCollectionID common.Hash
	ConditionID        common.Hash
	IndexSets          []*big.Int
}

func BuildRedeemTransaction(params RedeemParams) (*Transaction, error) {
	data, err := encodeRedeemPositions(params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode redeem positions: %w", err)
	}

	return &Transaction{
		To:    params.ConditionalTokens,
		Data:  data,
		Value: big.NewInt(0),
	}, nil
}

func encodeRedeemPositions(params RedeemParams) ([]byte, error) {
	redeemABI := `[{"constant":false,"inputs":[{"name":"collateralToken","type":"address"},{"name":"parentCollectionId","type":"bytes32"},{"name":"conditionId","type":"bytes32"},{"name":"indexSets","type":"uint256[]"}],"name":"redeemPositions","outputs":[]}]`

	parsedABI, err := abi.JSON(strings.NewReader(redeemABI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse ABI: %w", err)
	}

	data, err := parsedABI.Pack("redeemPositions",
		params.CollateralToken,
		params.ParentCollectionID,
		params.ConditionID,
		params.IndexSets,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to pack arguments: %w", err)
	}

	return data, nil
}
