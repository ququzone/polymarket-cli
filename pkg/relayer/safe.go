package relayer

import (
	"bytes"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"

	"polymarket-cli/pkg/relayer/transactions"
)

const (
	SAFEInitCodeHash     = "0x2bce2127ff07fb632d16c8347c4ebf501f4841168bed00d9e6ef715ddb6fcecf"
	PROXYInitCodeHashHex = "0xd21df8dc65880a8606f09fe0ce3df9b8869287ab0b058be05aa9e8af6330a00b"
	SafeFactory          = "0xaacFeEa03eb1561C4e67d661e40682Bd20E3541b"
	SafeMultisend        = "0xA238CBeb142c10Ef7Ad8442C6D1f9E89e07e7761"
	ZeroAddress          = "0x0000000000000000000000000000000000000000"
)

func DeriveProxyWallet(address, proxyFactory common.Address) common.Address {
	salt := crypto.Keccak256Hash(address.Bytes())
	return calculateCreate2Address(proxyFactory, salt, common.HexToHash(PROXYInitCodeHashHex))
}

func DeriveSafe(address, safeFactory common.Address) common.Address {
	result, _ := encodeAbiParameters(
		[]string{"address"},
		[]any{
			address,
		},
	)

	salt := crypto.Keccak256Hash(result)
	return calculateCreate2Address(safeFactory, salt, common.HexToHash(SAFEInitCodeHash))
}

func calculateCreate2Address(from common.Address, salt, initCodeHash common.Hash) common.Address {
	data := make([]byte, 1+20+32+32)
	data[0] = 0xff
	copy(data[1:21], from.Bytes())
	copy(data[21:53], salt.Bytes())
	copy(data[53:85], initCodeHash.Bytes())

	hash := crypto.Keccak256Hash(data)
	return common.BytesToAddress(hash.Bytes()[12:])
}

func encodeAbiParameters(types []string, values []any) ([]byte, error) {
	if len(types) != len(values) {
		return nil, fmt.Errorf("types and values length mismatch")
	}

	args := make(abi.Arguments, 0, len(types))

	for _, t := range types {
		abiType, err := abi.NewType(t, "", nil)
		if err != nil {
			return nil, err
		}
		args = append(args, abi.Argument{
			Type: abiType,
		})
	}

	return args.Pack(values...)
}

func aggregateTransaction(txs []*transactions.SafeTransaction, safeMultisend common.Address) (*transactions.SafeTransaction, error) {
	if len(txs) == 1 {
		return txs[0], nil
	}
	args := encodePackedTxs(txs)

	multisendABI := `
	[{
      "constant": false,
      "inputs": [
        {
          "internalType": "bytes",
          "name": "transactions",
          "type": "bytes"
        }
      ],
      "name": "multiSend",
      "outputs": [],
      "payable": false,
      "stateMutability": "nonpayable",
      "type": "function"
    }]
	`

	parsedABI, err := abi.JSON(strings.NewReader(multisendABI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse ABI: %w", err)
	}

	data, err := parsedABI.Pack("multiSend", args)
	if err != nil {
		return nil, fmt.Errorf("failed to pack arguments: %w", err)
	}

	return &transactions.SafeTransaction{
		To:        safeMultisend,
		Operation: 1,
		Value:     big.NewInt(0),
		Data:      data,
	}, nil
}

func uint256ToBytes(v *big.Int) []byte {
	b := v.Bytes()
	if len(b) > 32 {
		panic("uint256 overflow")
	}
	padded := make([]byte, 32)
	copy(padded[32-len(b):], b)
	return padded
}

func encodePackedTx(tx *transactions.SafeTransaction) []byte {
	var buf bytes.Buffer

	buf.WriteByte(tx.Operation)
	buf.Write(tx.To.Bytes())
	buf.Write(uint256ToBytes(tx.Value))
	buf.Write(uint256ToBytes(big.NewInt(int64(len(tx.Data)))))
	buf.Write(tx.Data)

	return buf.Bytes()
}

func encodePackedTxs(txs []*transactions.SafeTransaction) []byte {
	var out []byte
	for _, tx := range txs {
		b := encodePackedTx(tx)
		out = append(out, b...)
	}
	return out
}

type SafeTxData struct {
	To             common.Address
	Value          *big.Int
	Data           []byte
	Operation      uint8
	SafeTxGas      *big.Int
	BaseGas        *big.Int
	GasPrice       *big.Int
	GasToken       common.Address
	RefundReceiver common.Address
	Nonce          *big.Int
}

func CreateStructHash(
	chainId *big.Int,
	safe common.Address,
	to common.Address,
	value *big.Int,
	data []byte,
	operation uint8,
	safeTxGas *big.Int,
	baseGas *big.Int,
	gasPrice *big.Int,
	gasToken common.Address,
	refundReceiver common.Address,
	nonce *big.Int,
) ([]byte, error) {
	typeData := apitypes.TypedData{
		Types: apitypes.Types{
			"EIP712Domain": []apitypes.Type{
				{Name: "chainId", Type: "uint256"},
				{Name: "verifyingContract", Type: "address"},
			},
			"SafeTx": []apitypes.Type{
				{Name: "to", Type: "address"},
				{Name: "value", Type: "uint256"},
				{Name: "data", Type: "bytes"},
				{Name: "operation", Type: "uint8"},
				{Name: "safeTxGas", Type: "uint256"},
				{Name: "baseGas", Type: "uint256"},
				{Name: "gasPrice", Type: "uint256"},
				{Name: "gasToken", Type: "address"},
				{Name: "refundReceiver", Type: "address"},
				{Name: "nonce", Type: "uint256"},
			},
		},
		PrimaryType: "SafeTx",
		Domain: apitypes.TypedDataDomain{
			ChainId:           math.NewHexOrDecimal256(chainId.Int64()),
			VerifyingContract: safe.Hex(),
		},
		Message: apitypes.TypedDataMessage{
			"to":             to.Hex(),
			"value":          value,
			"data":           data,
			"operation":      fmt.Sprint(operation),
			"safeTxGas":      safeTxGas,
			"baseGas":        baseGas,
			"gasPrice":       gasPrice,
			"gasToken":       gasToken.Hex(),
			"refundReceiver": refundReceiver.Hex(),
			"nonce":          nonce,
		},
	}

	typedDataHash, err := typeData.HashStruct(typeData.PrimaryType, typeData.Message)
	if err != nil {
		return nil, err
	}
	domainSeparator, err := typeData.HashStruct("EIP712Domain", typeData.Domain.Map())
	if err != nil {
		return nil, err
	}

	rawData := fmt.Appendf(nil, "\x19\x01%s%s", string(domainSeparator), string(typedDataHash))
	return crypto.Keccak256(rawData), nil
}
