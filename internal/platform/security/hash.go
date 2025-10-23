package security

import (
	"crypto/subtle"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(pw string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	return string(b), err
}

func CheckPassword(hash, plain string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain))
	return err == nil
}

// subtleCompare untuk komparasi konstan
func subtleCompare(a, b string) bool { return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1 }
