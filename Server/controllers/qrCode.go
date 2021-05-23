package controllers

import (
	"github.com/skip2/go-qrcode"
)

func generateGoogleAuthenticatorQRCode(secret, label, issuer string) ([]byte, error) {
	authLink := "otpauth://totp/" + label + "?secret=" + secret + "&issuer=" + issuer
	return qrcode.Encode(authLink, qrcode.Medium, 256)
}
