package usecase

import (
	"errors"

	"qrmenu/internal/platform/logging"
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
	logging.UsecaseInfo("Auth.Login", "attempt", "auth_attempt", "email", email)

	a, err := u.adminRepo.FindActiveByEmail(email)
	if err != nil {
		logging.UsecaseError("Auth.Login", "admin lookup failed", "admin_lookup_failed", err, "email", email)
		return "", "", errors.New("invalid credentials")
	}
	if !security.CheckPassword(a.PasswordHash, password) {
		invalid := errors.New("invalid credentials")
		logging.UsecaseError("Auth.Login", "password mismatch", "password_mismatch", invalid, "email", email)
		return "", "", invalid
	}
	tok, err := u.jwt.SignAdmin(a.ID, a.Email, a.TenantID) // include tenant ID in the JWT
	if err != nil {
		logging.UsecaseError("Auth.Login", "failed to sign token", "token_sign_failed", err, "email", email, "tenant_id", a.TenantID)
		return "", "", err
	}
	logging.UsecaseInfo("Auth.Login", "token issued", "token_issued", "email", email, "tenant_id", a.TenantID)
	return tok, a.TenantID, nil
}

// Optional: implement token blacklist/revocation for production setups.
func (u *AuthUC) Logout(rawToken string) error { return nil }
