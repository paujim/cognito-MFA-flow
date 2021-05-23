package entities

type MFAVerifyRequest struct {
	AccessToken *string `form:"accessToken" json:"accessToken" binding:"required"`
	Code        *string `form:"code" json:"code" binding:"required"`
	DeviceName  *string `form:"deviceName" json:"deviceName"`
}
