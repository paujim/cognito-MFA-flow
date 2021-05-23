package entities

type TokenCodeRequest struct {
	Username *string `form:"username" json:"username" binding:"required"`
	Code     *string `form:"code" json:"code" binding:"required"`
	Session  *string `form:"session" json:"session" binding:"required"`
}
