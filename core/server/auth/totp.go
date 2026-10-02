package auth

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"strings"
	"time"
)

// GenerateTOTPSecret returns a base32 secret suitable for authenticator apps.
func GenerateTOTPSecret() (string, error) {
	buf := make([]byte, 20)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(buf), nil
}

// ValidateTOTP checks code against secret with ±1 step window.
func ValidateTOTP(secret, code string) bool {
	secret = strings.TrimSpace(strings.ToUpper(secret))
	code = strings.TrimSpace(code)
	if secret == "" || len(code) != 6 {
		return false
	}
	key, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(secret)
	if err != nil {
		return false
	}
	now := time.Now().UTC().Unix() / 30
	for _, step := range []int64{now - 1, now, now + 1} {
		if fmt.Sprintf("%06d", hotp(key, step)) == code {
			return true
		}
	}
	return false
}

func hotp(key []byte, counter int64) int32 {
	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], uint64(counter))
	mac := hmac.New(sha1.New, key)
	_, _ = mac.Write(buf[:])
	sum := mac.Sum(nil)
	offset := sum[len(sum)-1] & 0x0f
	truncated := binary.BigEndian.Uint32(sum[offset:offset+4]) & 0x7fffffff
	return int32(truncated % 1000000)
}
