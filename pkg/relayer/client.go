package relayer

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"polymarket-cli/pkg/relayer/transactions"
)

const (
	DefaultRelayerURL = "https://relayer-v2.polymarket.com/"
	PolygonChainID    = 137
)

var (
	CTF_ADDRESS  = common.HexToAddress("0x4D97DCd97eC945f40cF65F87097ACe5EA0476045")
	USDC_ADDRESS = common.HexToAddress("0x2791Bca1f2de4661ED88A30C99A7a9449Aa84174")
)

type BuilderCreds struct {
	Key        string
	Secret     string
	Passphrase string
}

type Client struct {
	baseURL    string
	chainId    *big.Int
	httpClient *http.Client
	creds      *BuilderCreds
	privateKey *ecdsa.PrivateKey
	address    common.Address
	txType     RelayerTxType
}

func NewClient(creds *BuilderCreds, txType RelayerTxType, owner *string, privateKeyHex *string) (*Client, error) {
	client := &Client{
		baseURL: DefaultRelayerURL,
		chainId: big.NewInt(PolygonChainID),
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		creds:  creds,
		txType: txType,
	}

	if owner != nil {
		client.address = common.HexToAddress(*owner)
	}

	if privateKeyHex != nil {
		privateKey, err := crypto.HexToECDSA(*privateKeyHex)
		if err != nil {
			return nil, fmt.Errorf("invalid private key: %w", err)
		}
		client.privateKey = privateKey

		publicKey := privateKey.PublicKey
		client.address = crypto.PubkeyToAddress(publicKey)
	}

	return client, nil
}

type RelayerTxType string

const (
	RelayerTxTypeSAFE  RelayerTxType = "SAFE"
	RelayerTxTypePROXY RelayerTxType = "PROXY"
)

type ExecuteResponse struct {
	TransactionID   string `json:"transactionID"`
	State           string `json:"state"`
	Hash            string `json:"hash"`
	TransactionHash string `json:"transactionHash"`
}

type NonceResponse struct {
	Nonce string `json:"nonce"`
}

func (c *Client) Execute(txs []*transactions.Transaction, metadata string) (*ExecuteResponse, error) {
	request, err := c.buildTransactionRequest(txs, metadata)
	if err != nil {
		return nil, err
	}

	bodyBytes, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}
	bodyStr := string(bodyBytes)

	req, err := http.NewRequest("POST", c.baseURL+"submit", bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Access-Control-Allow-Credentials", "true")

	if c.creds != nil {
		timestamp := time.Now().Unix()
		signature, err := BuildHmacSignature(c.creds.Secret, timestamp, "POST", "/submit", &bodyStr)
		if err != nil {
			return nil, err
		}

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

func (c *Client) GetNonce(signerAddress string) (*string, error) {
	req, err := http.NewRequest("GET", c.baseURL+"nonce", nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("address", signerAddress)
	q.Add("type", string(c.txType))
	req.URL.RawQuery = q.Encode()

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get request: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	var result NonceResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result.Nonce, nil
}

func (c *Client) buildTransactionRequest(txs []*transactions.Transaction, metadata string) (*transactions.TransactionRequest, error) {
	switch c.txType {
	case RelayerTxTypeSAFE:
		stxs := make([]*transactions.SafeTransaction, len(txs))
		for i, tx := range txs {
			stxs[i] = &transactions.SafeTransaction{
				To:        tx.To,
				Operation: 0,
				Data:      tx.Data,
				Value:     tx.Value,
			}
		}
		return c.buildSafeTransactionRequest(stxs, metadata)
	default:
		return nil, errors.New("unsupport type")
	}
}

func (c *Client) buildSafeTransactionRequest(txs []*transactions.SafeTransaction, metadata string) (*transactions.TransactionRequest, error) {
	nonce, err := c.GetNonce(c.address.Hex())
	if err != nil {
		return nil, err
	}
	nonceBig, _ := new(big.Int).SetString(*nonce, 10)

	safeFactory := common.HexToAddress(SafeFactory)
	safeMultisend := common.HexToAddress(SafeMultisend)
	transaction, err := aggregateTransaction(txs, safeMultisend)
	if err != nil {
		return nil, err
	}
	safeTxnGas := "0"
	baseGas := "0"
	gasPrice := "0"
	gasToken := ZeroAddress
	refundReceiver := ZeroAddress
	operationStr := fmt.Sprint(transaction.Operation)

	safeAddress := DeriveSafe(c.address, safeFactory).Hex()

	structHash, err := CreateStructHash(
		c.chainId,
		common.HexToAddress(safeAddress),
		transaction.To,
		transaction.Value,
		transaction.Data,
		transaction.Operation,
		big.NewInt(0),
		big.NewInt(0),
		big.NewInt(0),
		common.HexToAddress(gasToken),
		common.HexToAddress(refundReceiver),
		nonceBig,
	)
	if err != nil {
		return nil, err
	}

	signature, err := c.SignMessage(structHash)
	if err != nil {
		return nil, err
	}

	return &transactions.TransactionRequest{
		From:        c.address.Hex(),
		To:          transaction.To.Hex(),
		ProxyWallet: &safeAddress,
		Data:        "0x" + hex.EncodeToString(transaction.Data),
		Nonce:       nonce,
		Signature:   signature,
		SignatureParams: transactions.SignatureParams{
			GasPrice:       &gasPrice,
			Operation:      &operationStr,
			SafeTxnGas:     &safeTxnGas,
			BaseGas:        &baseGas,
			GasToken:       &gasToken,
			RefundReceiver: &refundReceiver,
		},
		Type:     "SAFE",
		Metadata: &metadata,
	}, nil
}

func (c *Client) SignMessage(message []byte) (string, error) {
	hash := accounts.TextHash(message)

	sig, err := crypto.Sign(hash, c.privateKey)
	if err != nil {
		return "", err
	}

	v, err := normalizeV(sig[64])
	if err != nil {
		return "", err
	}
	sig[64] = v

	return "0x" + hex.EncodeToString(sig), nil
}

func normalizeV(v byte) (byte, error) {
	switch v {
	case 0, 1:
		return v + 31, nil
	case 27, 28:
		return v + 4, nil
	default:
		return 0, errors.New("Invalid signature")
	}
}
