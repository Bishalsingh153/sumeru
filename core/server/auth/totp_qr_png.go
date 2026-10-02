package auth

import qrcode "github.com/skip2/go-qrcode"

// TOTPQRPng returns a PNG image for an otpauth URI.
func TOTPQRPng(otpauthURI string, size int) ([]byte, error) {
	if size <= 0 {
		size = 200
	}
	return qrcode.Encode(otpauthURI, qrcode.Medium, size)
}
