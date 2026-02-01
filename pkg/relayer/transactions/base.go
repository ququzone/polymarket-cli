package transactions

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type Transaction struct {
	To    common.Address `json:"to"`
	Data  []byte         `json:"data"`
	Value *big.Int       `json:"value"`
}

type SafeTransaction struct {
	To        common.Address `json:"to"`
	Operation uint8          `json:"operation"`
	Data      []byte         `json:"data"`
	Value     *big.Int       `json:"value"`
}

type SignatureParams struct {
	GasPrice        *string `json:"gasPrice,omitempty"`
	RelayerFee      *string `json:"relayerFee,omitempty"`
	GasLimit        *string `json:"gasLimit,omitempty"`
	RelayHub        *string `json:"relayHub,omitempty"`
	Relay           *string `json:"relay,omitempty"`
	Operation       *string `json:"operation,omitempty"`
	SafeTxnGas      *string `json:"safeTxnGas,omitempty"`
	BaseGas         *string `json:"baseGas,omitempty"`
	GasToken        *string `json:"gasToken,omitempty"`
	RefundReceiver  *string `json:"refundReceiver,omitempty"`
	PaymentToken    *string `json:"paymentToken,omitempty"`
	Payment         *string `json:"payment,omitempty"`
	PaymentReceiver *string `json:"paymentReceiver,omitempty"`
}

type TransactionRequest struct {
	Type            string          `json:"type"`
	From            string          `json:"from"`
	To              string          `json:"to"`
	ProxyWallet     *string         `json:"proxyWallet,omitempty"`
	Data            string          `json:"data"`
	Nonce           *string         `json:"nonce,omitempty"`
	Signature       string          `json:"signature"`
	SignatureParams SignatureParams `json:"signatureParams"`
	Metadata        *string         `json:"metadata,omitempty"`
}
