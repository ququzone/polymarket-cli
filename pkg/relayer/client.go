package relayer

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

const (
	DefaultRelayerURL = "https://relayer-v2.polymarket.com/"
	polygonChainID    = 137
)

type BuilderCreds struct {
	Key        string
	Secret     string
	Passphrase string
}

type Client struct {
	baseURL    string
	httpClient *http.Client
	creds      *BuilderCreds
	privateKey *ecdsa.PrivateKey
	address    common.Address
	txType     RelayerTxType
}

func NewClient(creds *BuilderCreds, privateKeyHex string) (*Client, error) {
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, fmt.Errorf("invalid private key: %w", err)
	}

	address := crypto.PubkeyToAddress(privateKey.PublicKey)

	return &Client{
		baseURL: DefaultRelayerURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		creds:      creds,
		privateKey: privateKey,
		address:    address,
	}, nil
}

type Transaction struct {
	To    string `json:"to"`
	Data  string `json:"data"`
	Value string `json:"value"`
}

type RelayerTxType string

const (
	RelayerTxTypeSAFE  RelayerTxType = "SAFE"
	RelayerTxTypePROXY RelayerTxType = "PROXY"
)

type executeRequest struct {
	Transactions []Transaction `json:"transactions"`
	TxType       RelayerTxType `json:"txType"`
	Metadata     string        `json:"metadata"`
}

type ExecuteResponse struct {
	TransactionID   string `json:"transactionID"`
	State           string `json:"state"`
	Hash            string `json:"hash"`
	TransactionHash string `json:"transactionHash"`
}

func (c *Client) Execute(txs []Transaction, metadata string) (*ExecuteResponse, error) {
	if c.txType == "" {
		c.txType = RelayerTxTypeSAFE
	}

	reqBody := executeRequest{
		Transactions: txs,
		TxType:       c.txType,
		Metadata:     metadata,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}
	bodyStr := string(bodyBytes)

	req, err := http.NewRequest("POST", c.baseURL+"submit", bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	if c.creds != nil {
		timestamp := time.Now().UnixMilli()
		signature := BuildHmacSignature(c.creds.Secret, timestamp, "POST", "/submit", &bodyStr)

		req.Header.Set("POLY_BUILDER_API_KEY", c.creds.Key)
		req.Header.Set("POLY_BUILDER_TIMESTAMP", strconv.FormatInt(timestamp, 10))
		req.Header.Set("POLY_BUILDER_PASSPHRASE", c.creds.Passphrase)
		req.Header.Set("POLY_BUILDER_SIGNATURE", signature)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	var result ExecuteResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}
