package services

import (
	"errors"
	"fmt"

	"github.com/Transcendence/models"
	"github.com/Transcendence/repositories"
	"github.com/pquerna/otp/totp"
)

type SetupResponse struct {
	Secret string `json:"secret"`
	QRCodeURL string `json="qr_code_url"`
}

type TwoFAService struct {
	userRepo repositories.UserRepository
	issuer string
}

func NewTwoFAService(repo repositories.UserRepository) *TwoFAService {
	return &TwoFAService{
		userRepo: repo,
		issuer: "Transcendence",
	}
}

func (s *TwoFAService) GenerateSecret(userID string) (*SetupResponse, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	if user.TwoFAEnabled {
		return nil, errors.New("2FA is already enabled for this account")
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer: s.issuer,
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
		Secret: secret,
		QRCodeURL: key.URL(),
	}, nil
}

func (s *TwoFAService) EnableTwoFA(userID string, code string) error {
	
}
