package web

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"strings"
	"time"

	"sumeru/core/server/config"
)

func webMACSecret() []byte {
	secret := strings.TrimSpace(config.AppConfig.CSRFSecret)
	if secret == "" {
		secret = "dev-web-cookie-mac"
	}
	return []byte(secret)
}

func signWebCookiePayload(payload []byte) string {
	mac := hmac.New(sha256.New, webMACSecret())
	_, _ = mac.Write(payload)
	sig := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
	body := base64.RawURLEncoding.EncodeToString(payload)
	return body + "." + sig
}

func verifyWebCookiePayload(value string, out interface{}) bool {
	parts := strings.Split(strings.TrimSpace(value), ".")
	if len(parts) != 2 {
		return false
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return false
	}
	sig, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return false
	}
	mac := hmac.New(sha256.New, webMACSecret())
	_, _ = mac.Write(payload)
	if !hmac.Equal(sig, mac.Sum(nil)) {
		return false
	}
	return json.Unmarshal(payload, out) == nil
}

type signedUIDPayload struct {
	UID int   `json:"uid"`
	Exp int64 `json:"exp"`
}

func mintSignedUIDCookieValue(uid int, ttl time.Duration) string {
	payload, _ := json.Marshal(signedUIDPayload{
		UID: uid,
		Exp: time.Now().Add(ttl).Unix(),
	})
	return signWebCookiePayload(payload)
}

func parseSignedUIDCookie(value string) (uid int, ok bool) {
	var p signedUIDPayload
	if !verifyWebCookiePayload(value, &p) {
		return 0, false
	}
	if p.UID <= 0 || p.Exp <= 0 || time.Now().Unix() > p.Exp {
		return 0, false
	}
	return p.UID, true
}
