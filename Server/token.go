package main

import (
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/gin-gonic/gin"
)

type TokenRequest struct {
	Username *string `form:"username" json:"username" binding:"required"`
	Password *string `form:"password" json:"password" binding:"required"`
}
type TokenUpdateRequest struct {
	Username *string `form:"username" json:"username" binding:"required"`
	Password *string `form:"password" json:"password" binding:"required"`
	Session  *string `form:"session" json:"session" binding:"required"`
}
type TokenCodeRequest struct {
	Username *string `form:"username" json:"username" binding:"required"`
	Code     *string `form:"code" json:"code" binding:"required"`
	Session  *string `form:"session" json:"session" binding:"required"`
}

type TokenResponse struct {
	Message      *string `json:"message"`
	Session      *string `json:"sesstion"`
	AccessToken  *string `json:"accessToken"`
	RefreshToken *string `json:"refreshToken"`
}

func successfulResponse(c *gin.Context, result *cognito.AuthenticationResultType) {
	c.JSON(http.StatusOK, TokenResponse{
		Message:      aws.String("Success"),
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
	})
}

func (app *App) addTokenRoutes(rg *gin.RouterGroup) {
	token := rg.Group("/token")

	token.POST("/", func(c *gin.Context) {

		var request TokenRequest
		if err := c.ShouldBind(&request); err != nil {
			c.JSON(http.StatusBadRequest, TokenResponse{Message: aws.String("Missing required parameter")})
			return
		}
		res, err := app.CognitoClient.InitiateAuth(&cognito.InitiateAuthInput{
			AuthFlow: aws.String("USER_PASSWORD_AUTH"),
			AuthParameters: map[string]*string{
				"USERNAME": request.Username,
				"PASSWORD": request.Password,
			},
			ClientId: aws.String(app.AppClientID),
		})
		if err != nil {
			c.JSON(http.StatusUnauthorized, TokenResponse{Message: aws.String("Incorrect username or password")})
			log.Printf(err.Error())
			return
		}

		if res.ChallengeName != nil && *res.ChallengeName == "NEW_PASSWORD_REQUIRED" {
			c.JSON(http.StatusOK, TokenResponse{
				Message: aws.String("New password required"),
				Session: res.Session,
			})
			return
		}
		if res.ChallengeName != nil && *res.ChallengeName == "SOFTWARE_TOKEN_MFA" {
			c.JSON(http.StatusOK, TokenResponse{
				Message: aws.String("MFA required"),
				Session: res.Session,
			})
			return
		}
		successfulResponse(c, res.AuthenticationResult)
	})
	token.POST("/update", func(c *gin.Context) {

		var request TokenUpdateRequest
		if err := c.ShouldBind(&request); err != nil {
			c.JSON(http.StatusBadRequest, TokenResponse{Message: aws.String("Missing required parameter")})
			return
		}
		res, err := app.CognitoClient.RespondToAuthChallenge(&cognitoidentityprovider.RespondToAuthChallengeInput{
			Session:       request.Session,
			ChallengeName: aws.String("NEW_PASSWORD_REQUIRED"),
			ClientId:      aws.String(app.AppClientID),
			ChallengeResponses: map[string]*string{
				"USERNAME":     request.Username,
				"NEW_PASSWORD": request.Password,
			},
		})
		if err != nil {
			c.JSON(http.StatusUnauthorized, TokenResponse{Message: aws.String("Unable to update password")})
			log.Printf(err.Error())
			return
		}
		successfulResponse(c, res.AuthenticationResult)

	})
	token.POST("/code", func(c *gin.Context) {

		var request TokenCodeRequest
		if err := c.ShouldBind(&request); err != nil {
			c.JSON(http.StatusBadRequest, TokenResponse{Message: aws.String("Missing required parameter")})
			return
		}

		res, err := app.CognitoClient.RespondToAuthChallenge(&cognitoidentityprovider.RespondToAuthChallengeInput{
			Session:       request.Session,
			ChallengeName: aws.String("SOFTWARE_TOKEN_MFA"),
			ClientId:      aws.String(app.AppClientID),
			ChallengeResponses: map[string]*string{
				"USERNAME":                request.Username,
				"SOFTWARE_TOKEN_MFA_CODE": request.Code,
			},
		})
		if err != nil {
			c.JSON(http.StatusUnauthorized, TokenResponse{Message: aws.String("Unable to validate code")})
			log.Printf(err.Error())
			return
		}

		successfulResponse(c, res.AuthenticationResult)
	})
}
