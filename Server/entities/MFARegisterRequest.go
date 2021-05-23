package entities

type MFARegisterRequest struct {
	AccessToken *string `form:"accessToken" json:"accessToken" binding:"required"`
}
