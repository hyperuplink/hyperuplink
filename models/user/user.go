package user

import (
	"bytes"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"io"
	"runtime"
	"strings"
	"time"

	"github.com/mrusme/hyperuplink/errs"
	uuid "github.com/vgarvardt/pgx-google-uuid/v5"
	"golang.org/x/crypto/argon2"
)

const (
	HashMem    = 64 * 1024
	HashIter   = 3
	HashSaltln = 16
	HashKeyln  = 32
)

type User struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Role     string    `json:"role"`

	Password            string    `json:"password"`
	PasswordResetToken  string    `json:"password_reset_token"`
	PasswordResetSentAt time.Time `json:"password_reset_sent_at"`

	Email                   string    `json:"email"`
	EmailUnconfirmed        string    `json:"email_unconfirmed"`
	EmailConfirmationToken  string    `json:"email_confirmation_token"`
	EmailConfirmationSentAt time.Time `json:"email_confirmation_sent_at"`
	EmailConfirmedAt        time.Time `json:"email_confirmed_at"`

	OTPEnabled  bool   `json:"otp_enabled"`
	OTPSecret   string `json:"otp_secret"`
	OTPTimestep int    `json:"otp_timestep"`

	SignInLastAt         time.Time `json:"sign_in_last_at"`
	SignInFailedAttempts int       `json:"sign_in_failed_attempts"`
	SignInLockedAt       time.Time `json:"sign_in_locked_at"`
	SignInUnlockToken    string    `json:"sign_in_unlock_token"`

	ProfilePicture string `json:"profile_picture"`
	Signature      string `json:"signature"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	BannedAt  time.Time `json:"banned_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

type PasswordParams struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

func (m *User) generateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (m *User) decodePassword() (
	params *PasswordParams,
	salt,
	key []byte,
	err error,
) {
	r := strings.NewReader(m.Password)

	_, err = fmt.Fscanf(r, "$argon2id$")
	if err != nil {
		return nil, nil, nil, errs.ErrHashVariantIncompatible
	}

	var version int
	_, err = fmt.Fscanf(r, "v=%d$", &version)
	if err != nil {
		return nil, nil, nil, err
	}
	if version != argon2.Version {
		return nil, nil, nil, errs.ErrHashVersionIncompatible
	}

	params = &PasswordParams{}
	_, err = fmt.Fscanf(r,
		"m=%d,t=%d,p=%d$",
		&params.Memory,
		&params.Iterations,
		&params.Parallelism,
	)
	if err != nil {
		return nil, nil, nil, err
	}

	rest, err := io.ReadAll(r)
	if err != nil {
		return nil, nil, nil, err
	}
	if bytes.ContainsAny(rest, "\r\n") {
		return nil, nil, nil, errs.ErrHashInvalid
	}

	var i int
	if i = bytes.IndexByte(rest, '$'); i == -1 {
		return nil, nil, nil, errs.ErrHashInvalid
	}

	b64Enc := base64.RawStdEncoding.Strict()

	salt = make([]byte, b64Enc.DecodedLen(i))
	_, err = b64Enc.Decode(salt, rest[:i])
	if err != nil {
		return nil, nil, nil, err
	}
	params.SaltLength = uint32(len(salt))

	key = make([]byte, b64Enc.DecodedLen(len(rest)-i-1))
	_, err = b64Enc.Decode(key, rest[i+1:])
	if err != nil {
		return nil, nil, nil, err
	}
	params.KeyLength = uint32(len(key))

	return params, salt, key, nil
}

func (m *User) SetPassword(password string) error {
	salt, err := m.generateRandomBytes(HashSaltln)
	if err != nil {
		return err
	}

	parallelism := uint8(runtime.NumCPU())

	key := argon2.IDKey(
		[]byte(password),
		salt,
		HashIter,
		HashMem,
		parallelism,
		HashKeyln,
	)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Key := base64.RawStdEncoding.EncodeToString(key)

	m.Password = fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		HashMem,
		HashIter,
		parallelism,
		b64Salt,
		b64Key,
	)
	return nil
}

func (m *User) CheckPassword(password string) (
	match bool,
	params *PasswordParams,
	err error,
) {
	params, salt, key, err := m.decodePassword()
	if err != nil {
		return false, nil, err
	}

	otherKey := argon2.IDKey(
		[]byte(password),
		salt,
		params.Iterations,
		params.Memory,
		params.Parallelism,
		params.KeyLength,
	)

	keyLen := int32(len(key))
	otherKeyLen := int32(len(otherKey))

	if subtle.ConstantTimeEq(keyLen, otherKeyLen) == 0 {
		return false, params, nil
	}
	if subtle.ConstantTimeCompare(key, otherKey) == 1 {
		return true, params, nil
	}
	return false, params, nil
}

func (m *User) ResetEmailConfirmationToken() error {
	const charset string = "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"0123456789"
	const length uint32 = 8

	b, err := m.generateRandomBytes(length)
	if err != nil {
		return err
	}

	for i := uint32(0); i < length; i++ {
		b[i] = charset[int(b[i])%len(charset)]
	}

	m.EmailConfirmationToken = string(b)

	return nil
}
