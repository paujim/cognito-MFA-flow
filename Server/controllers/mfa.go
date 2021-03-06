package controllers

import (
	"log"
	"net/http"

	b64 "encoding/base64"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/gin-gonic/gin"
	"github.com/paujim/cognito-MFA-flow/Server/entities"
)

const (
	GoogleAutheticatorLabel  = "PJ-Test"
	GoogleAutheticatorIssuer = "PJ"
)

func successfulMFAResponse(c *gin.Context, secretCode, googleAutheticator *string) {
	c.JSON(http.StatusOK, entities.MFAResponse{
		Message:            aws.String("Success"),
		SecretCode:         secretCode,
		GoogleAutheticator: googleAutheticator,
	})
}

func (app *App) addMFARoutes() {
	mfa := app.Router.Group("/mfa")

	mfa.POST("/register", func(c *gin.Context) {

		var request entities.MFARegisterRequest
		if err := c.ShouldBind(&request); err != nil {
			c.JSON(http.StatusBadRequest, entities.MFAResponse{Message: aws.String("Missing required parameter")})
			return
		}

		res, err := app.cognitoAPI.AssociateSoftwareToken(
			&cognitoidentityprovider.AssociateSoftwareTokenInput{
				AccessToken: request.AccessToken,
			})
		if err != nil {
			c.JSON(http.StatusUnauthorized, entities.MFAResponse{Message: aws.String("Unable to regiter device")})
			log.Printf("error: %s", err.Error())
			return
		}
		var encoded *string
		if raw, err := generateGoogleAuthenticatorQRCode(*res.SecretCode, GoogleAutheticatorLabel, GoogleAutheticatorIssuer); err == nil {
			encoded = aws.String(b64.StdEncoding.EncodeToString(raw))
		}

		successfulMFAResponse(c, res.SecretCode, encoded)
	})
	mfa.POST("/enable", func(c *gin.Context) {

		var request entities.MFAEnableRequest
		if err := c.ShouldBind(&request); err != nil {
			c.JSON(http.StatusBadRequest, entities.MFAResponse{Message: aws.String("Missing required parameter")})
			return
		}

		_, err := app.cognitoAPI.SetUserMFAPreference(
			&cognitoidentityprovider.SetUserMFAPreferenceInput{
				AccessToken: request.AccessToken,
				SoftwareTokenMfaSettings: &cognitoidentityprovider.SoftwareTokenMfaSettingsType{
					Enabled:      aws.Bool(true),
					PreferredMfa: aws.Bool(true),
				},
			})
		if err != nil {
			log.Printf("error :%s", err.Error())
		}

		successfulMFAResponse(c, nil, nil)
	})
	mfa.POST("/verify", func(c *gin.Context) {

		var request entities.MFAVerifyRequest
		if err := c.ShouldBind(&request); err != nil {
			c.JSON(http.StatusBadRequest, entities.MFAResponse{Message: aws.String("Missing required parameter")})
			return
		}

		res, err := app.cognitoAPI.VerifySoftwareToken(
			&cognitoidentityprovider.VerifySoftwareTokenInput{
				AccessToken:        request.AccessToken,
				UserCode:           request.Code,
				FriendlyDeviceName: request.DeviceName,
			})
		if err != nil {
			c.JSON(http.StatusBadRequest, entities.MFAResponse{Message: aws.String("Unable to verify code")})
			log.Printf("error :%s", err.Error())
			return
		}
		if res.Status != nil && *res.Status == "ERROR" {
			c.JSON(http.StatusBadRequest, entities.MFAResponse{Message: aws.String("Unable to verify code")})
			log.Printf("resp: %s", res.GoString())
			return
		}

		successfulMFAResponse(c, nil, nil)
	})

}
