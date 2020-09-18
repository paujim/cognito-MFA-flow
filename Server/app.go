package main

import (
	"log"

	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/gin-gonic/gin"
)

type IfaceApp interface {
	InitiateAuth(input *cognitoidentityprovider.InitiateAuthInput) (*cognitoidentityprovider.InitiateAuthOutput, error)
	RespondToAuthChallenge(input *cognitoidentityprovider.RespondToAuthChallengeInput) (*cognitoidentityprovider.RespondToAuthChallengeOutput, error)
	AssociateSoftwareToken(input *cognitoidentityprovider.AssociateSoftwareTokenInput) (*cognitoidentityprovider.AssociateSoftwareTokenOutput, error)
	SetUserMFAPreference(input *cognitoidentityprovider.SetUserMFAPreferenceInput) (*cognitoidentityprovider.SetUserMFAPreferenceOutput, error)
	VerifySoftwareToken(input *cognitoidentityprovider.VerifySoftwareTokenInput) (*cognitoidentityprovider.VerifySoftwareTokenOutput, error)
}

type App struct {
	CognitoClient IfaceApp
	UserPoolID    string
	AppClientID   string
	Router        *gin.Engine
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, UPDATE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
		} else {
			c.Next()
		}
	}
}

func createApp(userPoolID, appClientID string, cognitoClient IfaceApp) *App {

	app := &App{
		CognitoClient: cognitoClient,
		UserPoolID:    userPoolID,
		AppClientID:   appClientID,
		Router:        gin.Default(),
	}
	app.Router.Use(corsMiddleware())

	log.Printf("Cold start")
	app.addPingRoutes()
	app.addTokenRoutes()
	app.addMFARoutes()

	return app
}
