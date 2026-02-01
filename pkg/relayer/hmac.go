package relayer

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"strconv"
	"strings"
)

func BuildHmacSignature(
	secret string,
	timestamp int64,
	method string,
	requestPath string,
	body *string,
) string {
	message := strconv.FormatInt(timestamp, 10) + method + requestPath
	if body != nil {
		message += *body
	}

	base64Secret, err := base64.StdEncoding.DecodeString(secret)
	if err != nil {
		base64Secret = []byte(secret)
	}

	h := hmac.New(sha256.New, base64Secret)
	h.Write([]byte(message))
	sig := h.Sum(nil)

	sigBase64 := base64.StdEncoding.EncodeToString(sig)

	sigURLSafe := strings.ReplaceAll(sigBase64, "+", "-")
	sigURLSafe = strings.ReplaceAll(sigURLSafe, "/", "_")
	return sigURLSafe
}
