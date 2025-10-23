package usecase

import (
	"errors"
	"qrmenu/internal/platform/security"
	"qrmenu/internal/repository"
)

type AuthUC struct {
	adminRepo repository.AdminRepository
	jwt       *security.JWTMaker
}

func NewAuthUC(r repository.AdminRepository, j *security.JWTMaker) *AuthUC {
	return &AuthUC{adminRepo: r, jwt: j}
}

func (u *AuthUC) Login(email, password string) (token string, tenantID string, err error) {
    a, err := u.adminRepo.FindActiveByEmail(email)
    if err != nil { return "", "", errors.New("invalid credentials") }
    if !security.CheckPassword(a.PasswordHash, password) {
        return "", "", errors.New("invalid credentials")
    }
    tok, err := u.jwt.SignAdmin(a.ID, a.Email, a.TenantID)  // ✅ kirim tenantID
    return tok, a.TenantID, err
}

// (opsional) blacklist/token revoke di production
func (u *AuthUC) Logout(rawToken string) error { return nil }
