package usecase

import (
	"errors"
	"qrmenu/internal/domain"
	"qrmenu/internal/platform/security"
	"qrmenu/internal/repository"
	"strings"
)

type SetupUC struct {
	admins  repository.AdminRepository
	tenants repository.TenantRepository
	jwt     *security.JWTMaker
}

func NewSetupUC(a repository.AdminRepository, t repository.TenantRepository, jwt *security.JWTMaker) *SetupUC {
	return &SetupUC{admins: a, tenants: t, jwt: jwt}
}

type SetupTenantRequest struct {
	Email      string
	Password   string
	TenantCode string
	TenantName string
}

func (u *SetupUC) IsInitialized() (bool, error) {
	n, err := u.admins.Count()
	return n > 0, err
}

func (u *SetupUC) IsTenantInitialized(code string) (bool, error) {
	code = strings.TrimSpace(code)
	t, err := u.tenants.FindByCode(code)
	if err != nil {
		// kalau tenant belum ada → jelas belum initialized
		return false, nil
	}
	n, err := u.admins.CountActiveByTenant(t.ID)
	if err != nil { return false, err }
	return n > 0, nil
}

// Setup admin pertama untuk tenant tertentu.
// - Jika tenant belum ada → buat tenant + admin pertama.
// - Jika tenant ada & belum punya admin aktif → buat admin pertama.
// - Jika sudah ada admin aktif → error.
func (u *SetupUC) SetupFirstAdminForTenant(req SetupTenantRequest) (string, error) {
	code := strings.TrimSpace(req.TenantCode)
	if code == "" || req.Email == "" || req.Password == "" {
		return "", errors.New("invalid request")
	}

	// Cari/buat tenant
	t, err := u.tenants.FindByCode(code)
	if err != nil {
		// tenant belum ada → buat baru
		t = &domain.Tenant{Code: code, Name: strings.TrimSpace(req.TenantName)}
		if t.Name == "" { t.Name = code }
		if err := u.tenants.Create(t); err != nil {
			return "", err
		}
	}

	// Cek sudah punya admin aktif?
	n, err := u.admins.CountActiveByTenant(t.ID)
	if err != nil { return "", err }
	if n > 0 {
		return "", errors.New("tenant already initialized")
	}

	// Buat admin pertama
	hash, err := security.HashPassword(req.Password)
	if err != nil { return "", err }

	admin := &domain.AdminUser{
		TenantID:     t.ID,
		Email:        strings.ToLower(strings.TrimSpace(req.Email)),
		PasswordHash: hash,
		Name:         "Owner",
		Role:         "owner",
		IsActive:     true,
	}
	if err := u.admins.CreateForTenant(admin); err != nil {
		return "", err
	}

	// Sign JWT admin
	token, err := u.jwt.SignAdmin(admin.ID, admin.Email, admin.TenantID) // jika SignAdmin kamu hanya butuh (int,email)
	return token, err
}