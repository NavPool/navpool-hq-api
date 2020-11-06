package service

import (
	"encoding/base32"
	"errors"
	"github.com/NavPool/navpool-hq-api/internal/service/account"
	"github.com/NavPool/navpool-hq-api/internal/service/twofactor"
	"github.com/dgryski/dgoogauth"
	"math/rand"
	"strings"
	"time"
)

type TwoFactorService struct {
	accounts *AccountService
}

func NewTwoFactorService(accounts *AccountService) *TwoFactorService {
	return &TwoFactorService{accounts}
}

var (
	ErrTwoFactorAlreadyActive = errors.New("2FA already active on account")
	ErrTwoFactorInvalidCode   = errors.New("Authentication code is invalid")
	ErrTwoFactorNotActive     = errors.New("2FA is not active on account")
)

func (s *TwoFactorService) GetSecret(accountName string, user *account.User) (otp twofactor.Otp, err error) {
	if !user.TwoFactorExists() {
		user.TwoFactor = &account.TwoFactor{}
	}

	if user.TwoFactor.Secret == nil {
		random := make([]byte, 10)
		rand.Read(random)
		secret := strings.TrimRight(base32.StdEncoding.EncodeToString(random), "=")
		user.TwoFactor.Secret = &secret

		err := s.accounts.UpdateUser(user)
		if err != nil {
			return otp, err
		}
	}

	otp.Secret = *user.TwoFactor.Secret
	otp.OtpAuth = s.getOtpConfig(user.TwoFactor.Secret, user.TwoFactor.LastUsed).ProvisionURIWithIssuer(accountName, "NavPool")

	return
}

func (s *TwoFactorService) Enable(verification twofactor.Verification, user *account.User) (err error) {
	if user.TwoFactorActive() {
		return ErrTwoFactorAlreadyActive
	}

	success, lastUsed, err := s.Verify(user.TwoFactor.Secret, verification.Code, user.TwoFactor.LastUsed)
	if err != nil || success == false {
		return ErrTwoFactorInvalidCode
	}

	user.TwoFactor.Enable().LastUsed = &lastUsed
	err = s.accounts.UpdateUser(user)

	return
}

func (s *TwoFactorService) Disable(verification twofactor.Verification, user *account.User) (err error) {
	if !user.TwoFactorActive() {
		return ErrTwoFactorNotActive
	}

	verification.Secret = user.TwoFactor.Secret
	verification.LastUsed = user.TwoFactor.LastUsed

	success, lastUsed, err := s.Verify(user.TwoFactor.Secret, verification.Code, user.TwoFactor.LastUsed)
	if err != nil || success == false {
		return ErrTwoFactorInvalidCode
	}
	user.TwoFactor.Disable().LastUsed = &lastUsed

	err = s.accounts.UpdateUser(user)

	return
}

func (s *TwoFactorService) Verify(secret *string, code string, lastUsed *int) (bool, int, error) {
	t0 := int(time.Now().UTC().Unix() / 30)

	success, err := s.getOtpConfig(secret, lastUsed).Authenticate(code)

	return success, t0, err
}

func (s *TwoFactorService) getOtpConfig(secret *string, lastUsed *int) *dgoogauth.OTPConfig {
	var lastUsedList []int
	if lastUsed != nil {
		lastUsedList = append(lastUsedList, *lastUsed)
	}

	return &dgoogauth.OTPConfig{
		Secret:        *secret,
		WindowSize:    3,
		DisallowReuse: lastUsedList,
		HotpCounter:   0,
		UTC:           true,
	}
}
