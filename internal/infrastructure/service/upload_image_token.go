package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"strings"
	"time"

	"github.com/velonyapp/asset/internal/application/port"
	"github.com/velonyapp/asset/internal/conf"
)

type UploadImageToken struct {
	c *conf.Service
}

func NewUploadImageToken(c *conf.Service) port.UploadImageToken {
	return &UploadImageToken{c: c}
}

func (t *UploadImageToken) Sign(payload port.UploadImageTokenPayload) (string, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	encodedPayload := base64.RawURLEncoding.EncodeToString(data)

	hash := hmac.New(sha256.New, []byte(t.c.UploadTokenSecret))

	if _, err := hash.Write([]byte(encodedPayload)); err != nil {
		return "", err
	}

	signature := base64.RawURLEncoding.EncodeToString(hash.Sum(nil))

	return encodedPayload + "." + signature, nil
}

func (t *UploadImageToken) Verify(value string) (port.UploadImageTokenPayload, error) {
	parts := strings.Split(value, ".")
	if len(parts) != 2 {
		return port.UploadImageTokenPayload{}, port.ErrInvalidUploadToken
	}

	encodedPayload := parts[0]

	signature, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return port.UploadImageTokenPayload{}, port.ErrInvalidUploadToken
	}

	hash := hmac.New(sha256.New, []byte(t.c.UploadTokenSecret))

	if _, err := hash.Write([]byte(encodedPayload)); err != nil {
		return port.UploadImageTokenPayload{}, err
	}

	if !hmac.Equal(signature, hash.Sum(nil)) {
		return port.UploadImageTokenPayload{}, port.ErrInvalidUploadToken
	}

	data, err := base64.RawURLEncoding.DecodeString(encodedPayload)
	if err != nil {
		return port.UploadImageTokenPayload{}, port.ErrInvalidUploadToken
	}

	var payload port.UploadImageTokenPayload

	if err := json.Unmarshal(data, &payload); err != nil {
		return port.UploadImageTokenPayload{}, port.ErrInvalidUploadToken
	}

	if time.Now().After(payload.ExpireTime) {
		return port.UploadImageTokenPayload{}, port.ErrExpiredUploadToken
	}

	return payload, nil
}
