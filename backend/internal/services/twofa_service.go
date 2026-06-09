package services

import (
	"errors"
	"fmt"

	"github.com/pquerna/otp/totp"

	"ft_transcendence/backend/internal/repositories"
)

// SetupResponse holds the secret and the QR code URL for 2FA setup.
type SetupResponse struct {
	Secret    string `json:"secret"`
	QRCodeURL string `json:"qr_code_url"`
}

// TwoFAService handles two-factor authentication with TOTP.
type TwoFAService struct {
	userRepo repositories.UserRepository
	issuer   string
}

// NewTwoFAService creates a new TwoFAService.
func NewTwoFAService(repo repositories.UserRepository) *TwoFAService {
	return &TwoFAService{
		userRepo: repo,
		issuer:   "Transcendence",
	}
}

// GenerateSecret creates a new TOTP secret and stores it for the user. It fails if 2FA is already enabled.
func (s *TwoFAService) GenerateSecret(userID string) (*SetupResponse, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	if user.TwoFAEnabled {
		return nil, errors.New("2FA is already enabled for this account")
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      s.issuer,
		AccountName: user.Email,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate TOTP secret: %w", err)
	}

	secret := key.Secret()
	if err := s.userRepo.SetTwoFASecret(userID, secret); err != nil {
		return nil, fmt.Errorf("failed to store 2FA secret: %w", err)
	}

	return &SetupResponse{
		Secret:    secret,
		QRCodeURL: key.URL(),
	}, nil
}

// EnableTwoFA turns 2FA on for the user after the code is validated.
func (s *TwoFAService) EnableTwoFA(userID string, code string) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	if user.TwoFASecret == nil {
		return errors.New("no 2FA setup in progress; call /setup first")
	}

	if user.TwoFAEnabled {
		return errors.New("2FA is already enabled")
	}

	if !totp.Validate(code, *user.TwoFASecret) {
		return errors.New("invalid 2FA code")
	}

	return s.userRepo.SetTwoFAEnabled(userID, true)
}

// ValidateCode checks if the code is valid for the user. The user must have 2FA enabled.
func (s *TwoFAService) ValidateCode(userID, code string) (bool, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return false, fmt.Errorf("user not found: %w", err)
	}

	if !user.TwoFAEnabled || user.TwoFASecret == nil {
		return false, errors.New("2FA is not enabled for this user")
	}

	return totp.Validate(code, *user.TwoFASecret), nil
}

// DisableTwoFA turns 2FA off and clears the secret of the user.
func (s *TwoFAService) DisableTwoFA(userID string) error {
	return s.userRepo.ClearTwoFA(userID)
}
