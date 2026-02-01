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

type SignatureParams struct {
	GasPrice        *string `json:"gasPrice"`
	RelayerFee      *string `json:"relayerFee"`
	GasLimit        *string `json:"gasLimit"`
	RelayHub        *string `json:"relayHub"`
	Relay           *string `json:"relay"`
	Operation       *string `json:"operation"`
	SafeTxnGas      *string `json:"safeTxnGas"`
	BaseGas         *string `json:"baseGas"`
	GasToken        *string `json:"gasToken"`
	RefundReceiver  *string `json:"refundReceiver"`
	PaymentToken    *string `json:"paymentToken"`
	Payment         *string `json:"payment"`
	PaymentReceiver *string `json:"paymentReceiver"`
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
