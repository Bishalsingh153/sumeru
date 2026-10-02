package auth

import (
	"fmt"
	"net/url"
	"strings"
)

// TOTPOtpauthURI builds a standard otpauth URI for authenticator apps.
func TOTPOtpauthURI(issuer, accountName, secret string) string {
	issuer = strings.TrimSpace(issuer)
	accountName = strings.TrimSpace(accountName)
	secret = strings.TrimSpace(secret)
	label := url.PathEscape(issuer + ":" + accountName)
	if issuer == "" {
		label = url.PathEscape(accountName)
	}
	q := url.Values{}
	q.Set("secret", secret)
	if issuer != "" {
		q.Set("issuer", issuer)
	}
	return fmt.Sprintf("otpauth://totp/%s?%s", label, q.Encode())
}
