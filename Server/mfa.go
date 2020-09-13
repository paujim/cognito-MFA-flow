package main

import (
	"log"
	"net/http"

	b64 "encoding/base64"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/gin-gonic/gin"
)

type MFARegisterRequest struct {
	AccessToken *string `form:"accessToken" json:"accessToken" binding:"required"`
}

type MFAEnableRequest struct {
	AccessToken *string `form:"accessToken" json:"accessToken" binding:"required"`
}

type MFAVerifyRequest struct {
	AccessToken *string `form:"accessToken" json:"accessToken" binding:"required"`
	Code        *string `form:"code" json:"code" binding:"required"`
	DeviceName  *string `form:"deviceName" json:"deviceName"`
}

type MFAResponse struct {
	Message            *string `json:"message"`
	SecretCode         *string `json:"secret"`
	GoogleAutheticator *string `json:"googleAutheticator"`
}

func successfulMFAResponse(c *gin.Context, secretCode, googleAutheticator *string) {
	c.JSON(http.StatusOK, MFAResponse{
		Message:            aws.String("Success"),
		SecretCode:         secretCode,
		GoogleAutheticator: googleAutheticator,
	})
}

func (app *App) addMFARoutes(rg *gin.RouterGroup) {
	mfa := rg.Group("/mfa")

	mfa.POST("/register", func(c *gin.Context) {

		var request MFARegisterRequest
		if err := c.ShouldBind(&request); err != nil {
			c.JSON(http.StatusBadRequest, MFAResponse{Message: aws.String("Missing required parameter")})
			return
		}

		res, err := app.CognitoClient.AssociateSoftwareToken(
			&cognitoidentityprovider.AssociateSoftwareTokenInput{
				AccessToken: request.AccessToken,
			})
		if err != nil {
			c.JSON(http.StatusUnauthorized, MFAResponse{Message: aws.String("Unable to regiter device")})
			log.Printf(err.Error())
			return
		}
		var encoded *string
		if raw, err := generateGoogleAuthenticatorQRCode(*res.SecretCode, "PJ-Test", "PJ"); err == nil {
			encoded = aws.String(b64.StdEncoding.EncodeToString(raw))
		}

		successfulMFAResponse(c, res.SecretCode, encoded)
	})
	mfa.POST("/enable", func(c *gin.Context) {

		var request MFAEnableRequest
		if err := c.ShouldBind(&request); err != nil {
			c.JSON(http.StatusBadRequest, MFAResponse{Message: aws.String("Missing required parameter")})
			return
		}

		_, err := app.CognitoClient.SetUserMFAPreference(
			&cognitoidentityprovider.SetUserMFAPreferenceInput{
				AccessToken: request.AccessToken,
				SoftwareTokenMfaSettings: &cognitoidentityprovider.SoftwareTokenMfaSettingsType{
					Enabled:      aws.Bool(true),
					PreferredMfa: aws.Bool(true),
				},
			})
		if err != nil {
			log.Printf(err.Error())
		}

		successfulMFAResponse(c, nil, nil)
	})
	mfa.POST("/verify", func(c *gin.Context) {

		var request MFAVerifyRequest
		if err := c.ShouldBind(&request); err != nil {
			c.JSON(http.StatusBadRequest, MFAResponse{Message: aws.String("Missing required parameter")})
			return
		}

		res, err := app.CognitoClient.VerifySoftwareToken(
			&cognitoidentityprovider.VerifySoftwareTokenInput{
				AccessToken:        request.AccessToken,
				UserCode:           request.Code,
				FriendlyDeviceName: request.DeviceName,
			})
		if err != nil {
			c.JSON(http.StatusBadRequest, MFAResponse{Message: aws.String("Unable to verify code")})
			log.Printf(err.Error())
			return
		}
		if res.Status != nil && *res.Status == "ERROR" {
			c.JSON(http.StatusBadRequest, MFAResponse{Message: aws.String("Unable to verify code")})
			log.Printf(res.GoString())
			return
		}

		successfulMFAResponse(c, nil, nil)
	})

}
