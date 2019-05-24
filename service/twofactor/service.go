package twofactor

import (
	"encoding/base32"
	"github.com/NavPool/navpool-hq-api/logger"
	"github.com/NavPool/navpool-hq-api/service/account"
	"github.com/dgryski/dgoogauth"
	"math/rand"
	"strings"
	"time"
)

func GetSecret(accountName string, user account.User) (otp Otp, err error) {
	if !user.TwoFactorExists() {
		user.TwoFactor = &account.TwoFactor{}
	}

	if user.TwoFactor.Secret == nil {
		random := make([]byte, 10)
		rand.Read(random)
		secret := strings.TrimRight(base32.StdEncoding.EncodeToString(random), "=")
		user.TwoFactor.Secret = &secret

		err := account.UpdateUser(user)
		if err != nil {
			logger.LogError(err)
			return otp, err
		}
	}

	otp.Secret = *user.TwoFactor.Secret
	otp.OtpAuth = getOtpConfig(user.TwoFactor.Secret, user.TwoFactor.LastUsed).ProvisionURIWithIssuer(accountName, "NavPool")

	return
}

func Enable(verification Verification, user account.User) (err error) {
	if user.TwoFactorActive() {
		err = ErrTwoFactorAlreadyActive
		return
	}

	success, lastUsed, err := Verify(user.TwoFactor.Secret, verification.Code, user.TwoFactor.LastUsed)
	if err != nil || success == false {
		logger.LogError(err)
		err = ErrTwoFactorInvalidCode
		return
	}

	user.TwoFactor.Enable().LastUsed = &lastUsed

	err = account.UpdateUser(user)

	return
}

func Disable(verification Verification, user account.User) (err error) {
	if !user.TwoFactorActive() {
		err = ErrTwoFactorNotActive
		return
	}

	verification.Secret = user.TwoFactor.Secret
	verification.LastUsed = user.TwoFactor.LastUsed

	success, lastUsed, err := Verify(user.TwoFactor.Secret, verification.Code, user.TwoFactor.LastUsed)
	if err != nil || success == false {
		if err != nil {
			logger.LogError(err)
		}
		err = ErrTwoFactorInvalidCode
		return
	}
	user.TwoFactor.Disable().LastUsed = &lastUsed

	err = account.UpdateUser(user)

	return
}

func Verify(secret *string, code string, lastUsed *int) (bool, int, error) {
	t0 := int(time.Now().UTC().Unix() / 30)

	success, err := getOtpConfig(secret, lastUsed).Authenticate(code)

	return success, t0, err
}

func getOtpConfig(secret *string, lastUsed *int) *dgoogauth.OTPConfig {
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
