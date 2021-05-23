package entities

type MFAEnableRequest struct {
	AccessToken *string `form:"accessToken" json:"accessToken" binding:"required"`
}
