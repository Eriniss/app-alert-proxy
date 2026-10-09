package logto

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

// verifySignature는 body를 signing key로 HMAC-SHA256 한 값이
// 헤더로 받은 서명(hex)과 같은지 확인한다.
func verifySignature(key, body []byte, sigHex string) bool {
	got, err := hex.DecodeString(sigHex)
	if err != nil {
		return false
	}
	mac := hmac.New(sha256.New, key)
	mac.Write(body)
	return hmac.Equal(mac.Sum(nil), got)
}
