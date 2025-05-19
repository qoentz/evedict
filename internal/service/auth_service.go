package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
)

type AuthService struct {
	AuthSecret string
}

func NewAuthService(authSecret string) *AuthService {
	return &AuthService{
		AuthSecret: authSecret,
	}
}

func (s *AuthService) ValidateToken(token string) bool {
	expected := s.signToken("authenticated", s.AuthSecret)
	return hmac.Equal([]byte(token), []byte(expected))
}

func (s *AuthService) signToken(value, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(value))
	return hex.EncodeToString(h.Sum(nil))
}

func (s *AuthService) IssueToken(w http.ResponseWriter) {
	token := s.signToken("authenticated", s.AuthSecret)

	http.SetCookie(w, &http.Cookie{
		Name:     "auth",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,                // ⬅️ MUST be false on localhost, true in prod
		SameSite: http.SameSiteLaxMode, // ⬅️ Safe default for most apps
		MaxAge:   24 * 3600,
	})
}
