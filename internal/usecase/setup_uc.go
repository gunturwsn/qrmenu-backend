package usecase

import (
	"errors"
	"strings"

	"qrmenu/internal/domain"
	"qrmenu/internal/platform/logging"
	"qrmenu/internal/platform/security"
	"qrmenu/internal/repository"
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
	if err != nil {
		logging.UsecaseError("Setup.IsInitialized", "failed to count admins", "count_admins_failed", err)
		return false, err
	}
	initialized := n > 0
	logging.UsecaseInfo("Setup.IsInitialized", "result", "count_result", "initialized", initialized)
	return initialized, nil
}

func (u *SetupUC) IsTenantInitialized(code string) (bool, error) {
	code = strings.TrimSpace(code)
	t, err := u.tenants.FindByCode(code)
	if err != nil {
		// kalau tenant belum ada → jelas belum initialized
		logging.UsecaseInfo("Setup.IsTenantInitialized", "tenant not found", "tenant_not_found", "tenant_code", code)
		return false, nil
	}
	n, err := u.admins.CountActiveByTenant(t.ID)
	if err != nil {
		logging.UsecaseError("Setup.IsTenantInitialized", "failed to count tenant admins", "count_tenant_admins_failed", err, "tenant_code", code)
		return false, err
	}
	initialized := n > 0
	logging.UsecaseInfo("Setup.IsTenantInitialized", "result", "tenant_status", "tenant_code", code, "initialized", initialized)
	return initialized, nil
}

// Setup admin pertama untuk tenant tertentu.
// - Jika tenant belum ada → buat tenant + admin pertama.
// - Jika tenant ada & belum punya admin aktif → buat admin pertama.
// - Jika sudah ada admin aktif → error.
func (u *SetupUC) SetupFirstAdminForTenant(req SetupTenantRequest) (string, error) {
	code := strings.TrimSpace(req.TenantCode)
	if code == "" || req.Email == "" || req.Password == "" {
		err := errors.New("invalid request")
		logging.UsecaseError("Setup.SetupFirstAdminForTenant", "missing required fields", "invalid_request", err, "tenant_code", code, "email", req.Email)
		return "", err
	}

	// Cari/buat tenant
	t, err := u.tenants.FindByCode(code)
	if err != nil {
		// tenant belum ada → buat baru
		t = &domain.Tenant{Code: code, Name: strings.TrimSpace(req.TenantName)}
		if t.Name == "" {
			t.Name = code
		}
		if err := u.tenants.Create(t); err != nil {
			logging.UsecaseError("Setup.SetupFirstAdminForTenant", "failed to create tenant", "tenant_create_failed", err, "tenant_code", code)
			return "", err
		}
		logging.UsecaseInfo("Setup.SetupFirstAdminForTenant", "tenant created", "tenant_created", "tenant_code", code)
	}

	// Cek sudah punya admin aktif?
	n, err := u.admins.CountActiveByTenant(t.ID)
	if err != nil {
		logging.UsecaseError("Setup.SetupFirstAdminForTenant", "failed to count tenant admins", "count_tenant_admins_failed", err, "tenant_id", t.ID)
		return "", err
	}
	if n > 0 {
		err := errors.New("tenant already initialized")
		logging.UsecaseError("Setup.SetupFirstAdminForTenant", "tenant already initialized", "tenant_already_initialized", err, "tenant_code", code)
		return "", err
	}

	// Buat admin pertama
	hash, err := security.HashPassword(req.Password)
	if err != nil {
		logging.UsecaseError("Setup.SetupFirstAdminForTenant", "failed to hash password", "hash_failed", err, "tenant_code", code, "email", req.Email)
		return "", err
	}

	admin := &domain.AdminUser{
		TenantID:     t.ID,
		Email:        strings.ToLower(strings.TrimSpace(req.Email)),
		PasswordHash: hash,
		Name:         "Owner",
		Role:         "owner",
		IsActive:     true,
	}
	if err := u.admins.CreateForTenant(admin); err != nil {
		logging.UsecaseError("Setup.SetupFirstAdminForTenant", "failed to create admin", "admin_create_failed", err, "tenant_id", t.ID, "email", admin.Email)
		return "", err
	}

	// Sign JWT admin
	token, err := u.jwt.SignAdmin(admin.ID, admin.Email, admin.TenantID) // adjust parameters if your signer only needs (id,email)
	if err != nil {
		logging.UsecaseError("Setup.SetupFirstAdminForTenant", "failed to sign token", "token_sign_failed", err, "tenant_id", admin.TenantID, "email", admin.Email)
		return "", err
	}
	logging.UsecaseInfo("Setup.SetupFirstAdminForTenant", "admin created", "admin_created", "tenant_id", admin.TenantID, "email", admin.Email)
	return token, nil
}
