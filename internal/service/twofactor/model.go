package twofactor

type Verification struct {
	Code     string  `json:"code"`
	Secret   *string `json:"secret"`
	LastUsed *int    `json:"-"`
}

type Otp struct {
	Secret  string `json:"secret"`
	OtpAuth string `json:"otpauth"`
}
