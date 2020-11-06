package account

type Login struct {
	Username  string `form:"username" json:"username" binding:"required"`
	Password  string `form:"password" json:"password" binding:"required"`
	TwoFactor string `form:"twoFactor" json:"twoFactor"`
}

type Register struct {
	Username        string `form:"username" json:"username" binding:"required"`
	Password        string `form:"password" json:"password" binding:"required"`
	PasswordConfirm string `form:"passwordConfirm" json:"passwordConfirm" binding:"required"`
}
