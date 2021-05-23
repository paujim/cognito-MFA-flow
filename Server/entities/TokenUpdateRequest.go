package entities

type TokenUpdateRequest struct {
	Username *string `form:"username" json:"username" binding:"required"`
	Password *string `form:"password" json:"password" binding:"required"`
	Session  *string `form:"session" json:"session" binding:"required"`
}
