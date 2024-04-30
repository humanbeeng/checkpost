package core

import (
	"fmt"
	"strconv"
	"time"

	"github.com/aead/chacha20poly1305"
	db "github.com/humanbeeng/checkpost/server/db/sqlc"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/o1egl/paseto"
)

type PasetoVerifier struct {
	symmetricKey string
	paseto       *paseto.V2
}

func NewPasetoVerifier(symmetricKey string) (*PasetoVerifier, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size. Must be exactly %d chars", chacha20poly1305.KeySize)
	}

	return &PasetoVerifier{
		symmetricKey: symmetricKey,
		paseto:       paseto.NewV2(),
	}, nil
}

type CreateTokenArgs struct {
	Username string
	UserId   int64
	Plan     db.Plan
	Role     string
}

func (p *PasetoVerifier) CreateToken(args CreateTokenArgs, duration time.Duration) (string, error) {
	id, err := gonanoid.New()
	if err != nil {
		return "", err
	}

	jt := paseto.JSONToken{
		Issuer:     "checkpost",
		Jti:        id,
		Subject:    strconv.FormatInt(args.UserId, 10),
		IssuedAt:   time.Now(),
		Expiration: time.Now().Add(duration),
	}
	jt.Set("username", args.Username)
	jt.Set("plan", string(args.Plan))
	jt.Set("role", args.Role)

	token, err := p.paseto.Encrypt([]byte(p.symmetricKey), jt, nil)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (p *PasetoVerifier) VerifyToken(token string) (paseto.JSONToken, error) {
	var jt paseto.JSONToken

	err := p.paseto.Decrypt(token, []byte(p.symmetricKey), &jt, nil)
	if err != nil {
		return jt, err
	}

	if time.Now().After(jt.Expiration) {
		return jt, fmt.Errorf("token has expired")
	}

	return jt, nil
}
