package auth

import (
	"fmt"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/google/uuid"
	"github.com/o1egl/paseto"
)

type PasetoVerifier struct {
	symmetricKey string
	paseto       *paseto.V2
}

func NewPasetoVerifier(symmetricKey string) (*PasetoVerifier, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("Invalid key size. Must be exactly %d chars", chacha20poly1305.KeySize)
	}

	return &PasetoVerifier{
		symmetricKey: symmetricKey,
		paseto:       paseto.NewV2(),
	}, nil

}

type Payload struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

func (p *Payload) Valid() bool {
	return time.Now().After(p.ExpiredAt)
}

func (p *PasetoVerifier) CreateToken(username string, duration time.Duration) (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	payload := &Payload{
		ID:        id,
		Username:  username,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}

	token, err := p.paseto.Encrypt([]byte(p.symmetricKey), payload, nil)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (p *PasetoVerifier) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}

	err := p.paseto.Decrypt(token, []byte(p.symmetricKey), payload, nil)
	if err != nil {
		return nil, err
	}

	if !payload.Valid() {
		return nil, fmt.Errorf("Token has expired")
	}

	return payload, nil
}
