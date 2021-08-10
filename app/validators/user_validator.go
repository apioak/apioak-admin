package validators

type UserRegister struct {
	RePassword string `json:"re_password" zh:"确认密码" en:"Confirm Password" binding:"required,eqfield=Password"`
	Password   string `json:"password" zh:"密码" en:"Password" binding:"required"`
	Email      string `json:"email" zh:"邮箱" en:"Email" binding:"required,email"`
	Name       string `json:"name" zh:"昵称" en:"User name" binding:"required"`
}
