package auth_test

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"sumeru/core/server/auth"
)

func TestVerifyIDToken_expired(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	n := base64.RawURLEncoding.EncodeToString(key.N.Bytes())
	e := base64.RawURLEncoding.EncodeToString([]byte{0x01, 0x00, 0x01})
	jwks := `{"keys":[{"kty":"RSA","kid":"t1","use":"sig","alg":"RS256","n":"` + n + `","e":"` + e + `"}]}`
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(jwks))
	}))
	defer srv.Close()
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"RS256","kid":"t1","typ":"JWT"}`))
	payload, _ := json.Marshal(map[string]interface{}{
		"iss": "https://idp.example",
		"sub": "user-1",
		"aud": "client-abc",
		"exp": time.Now().Add(-time.Hour).Unix(),
	})
	payloadB64 := base64.RawURLEncoding.EncodeToString(payload)
	signingInput := header + "." + payloadB64
	digest := sha256.Sum256([]byte(signingInput))
	sig, err := rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, digest[:])
	if err != nil {
		t.Fatal(err)
	}
	token := signingInput + "." + base64.RawURLEncoding.EncodeToString(sig)
	_, err = auth.VerifyIDToken(context.Background(), token, "https://idp.example", "client-abc", srv.URL)
	if err == nil {
		t.Fatal("expected expired error")
	}
}

func TestVerifyIDToken_wrongIssuer(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	n := base64.RawURLEncoding.EncodeToString(key.N.Bytes())
	e := base64.RawURLEncoding.EncodeToString([]byte{0x01, 0x00, 0x01})
	jwks := `{"keys":[{"kty":"RSA","kid":"t1","use":"sig","alg":"RS256","n":"` + n + `","e":"` + e + `"}]}`
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(jwks))
	}))
	defer srv.Close()
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"RS256","kid":"t1","typ":"JWT"}`))
	payload, _ := json.Marshal(map[string]interface{}{
		"iss": "https://wrong.example",
		"sub": "user-1",
		"aud": "client-abc",
		"exp": time.Now().Add(time.Hour).Unix(),
	})
	payloadB64 := base64.RawURLEncoding.EncodeToString(payload)
	signingInput := header + "." + payloadB64
	digest := sha256.Sum256([]byte(signingInput))
	sig, err := rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, digest[:])
	if err != nil {
		t.Fatal(err)
	}
	token := signingInput + "." + base64.RawURLEncoding.EncodeToString(sig)
	_, err = auth.VerifyIDToken(context.Background(), token, "https://idp.example", "client-abc", srv.URL)
	if err == nil {
		t.Fatal("expected issuer mismatch")
	}
}

func TestVerifyIDToken_wrongAudience(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	n := base64.RawURLEncoding.EncodeToString(key.N.Bytes())
	e := base64.RawURLEncoding.EncodeToString([]byte{0x01, 0x00, 0x01})
	jwks := `{"keys":[{"kty":"RSA","kid":"t1","use":"sig","alg":"RS256","n":"` + n + `","e":"` + e + `"}]}`
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(jwks))
	}))
	defer srv.Close()
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"RS256","kid":"t1","typ":"JWT"}`))
	payload, _ := json.Marshal(map[string]interface{}{
		"iss": "https://idp.example",
		"sub": "user-1",
		"aud": "other-client",
		"exp": time.Now().Add(time.Hour).Unix(),
	})
	payloadB64 := base64.RawURLEncoding.EncodeToString(payload)
	signingInput := header + "." + payloadB64
	digest := sha256.Sum256([]byte(signingInput))
	sig, err := rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, digest[:])
	if err != nil {
		t.Fatal(err)
	}
	token := signingInput + "." + base64.RawURLEncoding.EncodeToString(sig)
	_, err = auth.VerifyIDToken(context.Background(), token, "https://idp.example", "client-abc", srv.URL)
	if err == nil {
		t.Fatal("expected audience mismatch")
	}
}

func TestVerifyIDToken_rs256(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	n := base64.RawURLEncoding.EncodeToString(key.N.Bytes())
	e := base64.RawURLEncoding.EncodeToString([]byte{0x01, 0x00, 0x01})
	jwks := `{"keys":[{"kty":"RSA","kid":"t1","use":"sig","alg":"RS256","n":"` + n + `","e":"` + e + `"}]}`
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(jwks))
	}))
	defer srv.Close()

	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"RS256","kid":"t1","typ":"JWT"}`))
	payload, _ := json.Marshal(map[string]interface{}{
		"iss":   "https://idp.example",
		"sub":   "user-1",
		"aud":   []interface{}{"client-abc", "other"},
		"exp":   time.Now().Add(time.Hour).Unix(),
		"email": "a@example.com",
	})
	payloadB64 := base64.RawURLEncoding.EncodeToString(payload)
	signingInput := header + "." + payloadB64
	digest := sha256.Sum256([]byte(signingInput))
	sig, err := rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, digest[:])
	if err != nil {
		t.Fatal(err)
	}
	token := signingInput + "." + base64.RawURLEncoding.EncodeToString(sig)

	claims, err := auth.VerifyIDToken(context.Background(), token, "https://idp.example", "client-abc", srv.URL)
	if err != nil {
		t.Fatalf("verify: %v", err)
	}
	if claims["sub"] != "user-1" {
		t.Fatalf("sub=%v", claims["sub"])
	}
}

func TestParseIDTokenClaimsUnsafe(t *testing.T) {
	email, sub, claims := auth.ParseIDTokenClaimsUnsafe("")
	if email != "" || sub != "" || claims != nil {
		t.Fatal("empty token should yield nothing")
	}
	payload := base64.RawURLEncoding.EncodeToString([]byte(`{"sub":"u1","email":"a@example.com"}`))
	token := "e30." + payload + ".sig"
	email, sub, claims = auth.ParseIDTokenClaimsUnsafe(token)
	if email != "a@example.com" || sub != "u1" || claims == nil {
		t.Fatalf("email=%q sub=%q claims=%v", email, sub, claims)
	}
}

func TestVerifyIDToken_rejectsMalformed(t *testing.T) {
	_, err := auth.VerifyIDToken(context.Background(), "not.a.jwt", "", "", "http://127.0.0.1")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestVerifyIDToken_rejectsUnsupportedAlg(t *testing.T) {
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))
	payload := base64.RawURLEncoding.EncodeToString([]byte(`{"sub":"x"}`))
	token := header + "." + payload + ".sig"
	_, err := auth.VerifyIDToken(context.Background(), token, "", "", "http://127.0.0.1/jwks")
	if err == nil {
		t.Fatal("expected unsupported alg error")
	}
}

func TestTOTPOtpauthURI(t *testing.T) {
	uri := auth.TOTPOtpauthURI("Sumeru", "admin@example.com", "JBSWY3DPEHPK3PXP")
	if uri == "" || len(uri) < 20 {
		t.Fatalf("uri=%q", uri)
	}
}

func TestTOTPQRPng(t *testing.T) {
	png, err := auth.TOTPQRPng(auth.TOTPOtpauthURI("Sumeru", "u", "JBSWY3DPEHPK3PXP"), 120)
	if err != nil || len(png) < 32 {
		t.Fatalf("png len=%d err=%v", len(png), err)
	}
}
