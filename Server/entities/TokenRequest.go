package entities

type TokenRequest struct {
	Username *string `form:"username" json:"username" binding:"required"`
	Password *string `form:"password" json:"password" binding:"required"`
}
