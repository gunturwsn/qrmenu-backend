package security

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTMaker struct{ secret []byte; ttl time.Duration }

func NewJWT(secret string, ttlMin int) *JWTMaker {
	return &JWTMaker{secret: []byte(secret), ttl: time.Duration(ttlMin) * time.Minute}
}

// Ganti signature supaya bisa kirim tenantID juga.
func (j *JWTMaker) SignAdmin(adminID string, email string, tenantID string) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"sub":       adminID,
		"email":     email,
		"role":      "admin",
		"tenant_id": tenantID,               // ✅ penting
		"iat":       now.Unix(),
		"exp":       now.Add(j.ttl).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secret)
}
