package controllers

import (
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/gin-gonic/gin"
	"github.com/paujim/cognito-MFA-flow/Server/entities"
)

func successfulResponse(c *gin.Context, result *cognitoidentityprovider.AuthenticationResultType) {
	c.JSON(http.StatusOK, entities.TokenResponse{
		Message:      aws.String("Success"),
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
	})
}

func (app *App) addTokenRoutes() {
	token := app.Router.Group("/token")

	token.POST("/", func(c *gin.Context) {

		var request entities.TokenRequest
		if err := c.ShouldBind(&request); err != nil {
			c.JSON(http.StatusBadRequest, entities.TokenResponse{Message: aws.String("Missing required parameter")})
			return
		}
		res, err := app.cognitoAPI.InitiateAuth(&cognitoidentityprovider.InitiateAuthInput{
			AuthFlow: aws.String("USER_PASSWORD_AUTH"),
			AuthParameters: map[string]*string{
				"USERNAME": request.Username,
				"PASSWORD": request.Password,
			},
			ClientId: aws.String(app.appClientID),
		})
		if err != nil {
			c.JSON(http.StatusUnauthorized, entities.TokenResponse{Message: aws.String("Incorrect username or password")})
			log.Printf("error :%s", err.Error())
			return
		}

		if res.ChallengeName != nil && *res.ChallengeName == "NEW_PASSWORD_REQUIRED" {
			c.JSON(http.StatusOK, entities.TokenResponse{
				Message: aws.String("New password required"),
				Session: res.Session,
			})
			return
		}
		if res.ChallengeName != nil && *res.ChallengeName == "SOFTWARE_TOKEN_MFA" {
			c.JSON(http.StatusOK, entities.TokenResponse{
				Message: aws.String("MFA required"),
				Session: res.Session,
			})
			return
		}
		successfulResponse(c, res.AuthenticationResult)
	})
	token.POST("/update", func(c *gin.Context) {

		var request entities.TokenUpdateRequest
		if err := c.ShouldBind(&request); err != nil {
			c.JSON(http.StatusBadRequest, entities.TokenResponse{Message: aws.String("Missing required parameter")})
			return
		}
		res, err := app.cognitoAPI.RespondToAuthChallenge(&cognitoidentityprovider.RespondToAuthChallengeInput{
			Session:       request.Session,
			ChallengeName: aws.String("NEW_PASSWORD_REQUIRED"),
			ClientId:      aws.String(app.appClientID),
			ChallengeResponses: map[string]*string{
				"USERNAME":     request.Username,
				"NEW_PASSWORD": request.Password,
			},
		})
		if err != nil {
			c.JSON(http.StatusUnauthorized, entities.TokenResponse{Message: aws.String("Unable to update password")})
			log.Printf("error :%s", err.Error())
			return
		}
		successfulResponse(c, res.AuthenticationResult)

	})
	token.POST("/code", func(c *gin.Context) {

		var request entities.TokenCodeRequest
		if err := c.ShouldBind(&request); err != nil {
			c.JSON(http.StatusBadRequest, entities.TokenResponse{Message: aws.String("Missing required parameter")})
			return
		}

		res, err := app.cognitoAPI.RespondToAuthChallenge(&cognitoidentityprovider.RespondToAuthChallengeInput{
			Session:       request.Session,
			ChallengeName: aws.String("SOFTWARE_TOKEN_MFA"),
			ClientId:      aws.String(app.appClientID),
			ChallengeResponses: map[string]*string{
				"USERNAME":                request.Username,
				"SOFTWARE_TOKEN_MFA_CODE": request.Code,
			},
		})
		if err != nil {
			c.JSON(http.StatusUnauthorized, entities.TokenResponse{Message: aws.String("Unable to validate code")})
			log.Printf("error :%s", err.Error())
			return
		}

		successfulResponse(c, res.AuthenticationResult)
	})
}
