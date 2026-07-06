package user

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image/png"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

const (
	OTPPeriod         = 30
	OTPSecretSize     = 20
	OTPDigits         = otp.DigitsSix
	OTPAlgorithm      = otp.AlgorithmSHA1
	OTPIssuerFallback = "Hyperuplink"

	otpQRSize = 256
)

type OTPEnrollment struct {
	URL       string
	Secret    string
	QRDataURI string
}

func NewOTPEnrollment(issuer, account, existingURL string) (*OTPEnrollment, error) {
	if issuer == "" {
		issuer = OTPIssuerFallback
	}

	var key *otp.Key
	var err error

	if existingURL != "" {
		key, err = otp.NewKeyFromURL(existingURL)
	} else {
		key, err = totp.Generate(totp.GenerateOpts{
			Issuer:      issuer,
			AccountName: account,
			Period:      OTPPeriod,
			SecretSize:  OTPSecretSize,
			Digits:      OTPDigits,
			Algorithm:   OTPAlgorithm,
		})
	}
	if err != nil {
		return nil, err
	}

	qr, err := otpQRDataURI(key)
	if err != nil {
		return nil, err
	}

	return &OTPEnrollment{
		URL:       key.String(),
		Secret:    key.Secret(),
		QRDataURI: qr,
	}, nil
}

func otpQRDataURI(key *otp.Key) (string, error) {
	img, err := key.Image(otpQRSize, otpQRSize)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err = png.Encode(&buf, img); err != nil {
		return "", err
	}

	return fmt.Sprintf(
		"data:image/png;base64,%s",
		base64.StdEncoding.EncodeToString(buf.Bytes()),
	), nil
}

func OTPSecretFromURL(uri string) (string, error) {
	key, err := otp.NewKeyFromURL(uri)
	if err != nil {
		return "", err
	}

	return key.Secret(), nil
}

func ValidateOTP(secret, code string) bool {
	return totp.Validate(code, secret)
}

func (m *User) EnableOTP(secret string) {
	m.OTPSecret = secret
	m.OTPTimestep = OTPPeriod
	m.OTPEnabled = true
}

func (m *User) DisableOTP() {
	m.OTPSecret = ""
	m.OTPTimestep = 0
	m.OTPEnabled = false
}
