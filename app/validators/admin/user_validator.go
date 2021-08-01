package admin

type Register struct {
	RePassword string `json:"re_password" zh:"确认密码" binding:"required,eqfield=Password"`
	Password string `json:"password" zh:"密码" binding:"required"`
	Email string `json:"email" zh:"邮箱" binding:"required,email"`
	Name string `json:"name" zh:"昵称" binding:"required"`
}

