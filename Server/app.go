package main

import (
	"log"

	"github.com/aws/aws-sdk-go/aws/session"
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/gin-gonic/gin"
)

type App struct {
	CognitoClient *cognito.CognitoIdentityProvider
	UserPoolID    string
	AppClientID   string
	Router        *gin.Engine
}

func createApp(userPoolID, appClientID string) *App {
	sess := session.Must(session.NewSession())
	app := &App{
		CognitoClient: cognito.New(sess),
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
