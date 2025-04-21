package util

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
)

func VerifyHMAC(r *http.Request) ([]byte, bool) {
	sig := r.Header.Get("X-Agent-Signature")
	if sig == "" {
		return nil, false
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, false
	}
	r.Body.Close()

	r.Body = io.NopCloser(bytes.NewReader(body)) // Reset for JSON decode

	mac := hmac.New(sha256.New, []byte(SharedSecret))
	mac.Write(body)
	expected := mac.Sum(nil)

	received, err := hex.DecodeString(sig)
	if err != nil {
		return nil, false
	}

	return body, hmac.Equal(received, expected)
}
