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

func createApp(userPoolID, appClientID string, cognitoClient IfaceApp) *App {
	app := &App{
		CognitoClient: cognitoClient,
		UserPoolID:    userPoolID,
		AppClientID:   appClientID,
		Router:        gin.Default(),
	}
	log.Printf("Cold start")
	v1 := app.Router.Group("/v1")
	app.addPingRoutes(v1)
	app.addTokenRoutes(v1)
	app.addMFARoutes(v1)

	return app
}
