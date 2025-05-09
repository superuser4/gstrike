package middleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"gstrike/pkg/util"
	"io"
	"net/http"
)

func AgentHMACAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		signature := r.Header.Get("X-Agent-Signature")
		if signature == "" {
			http.Error(w, "Missing signature", http.StatusUnauthorized)
			return
		}

		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading body", http.StatusBadRequest)
			return
		}
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // restore body

		mac := hmac.New(sha256.New, []byte(util.SharedSecret))
		mac.Write(bodyBytes)
		expectedMAC := hex.EncodeToString(mac.Sum(nil))

		if !hmac.Equal([]byte(signature), []byte(expectedMAC)) {
			http.Error(w, "Unauthorized (HMAC failed)", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
