package main

import (
	"context"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
)

var ginLambda *ginadapter.GinLambda
var app *App

func init() {
	sess := session.Must(session.NewSession())
	app = createApp(os.Getenv("USER_POOL_ID"), os.Getenv("CLIENT_ID"), cognito.New(sess))
	ginLambda = ginadapter.New(app.Router)
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != "" {
		// running in the cloud
		lambda.Start(Handler)
	} else {
		// running locally
		app.Router.Run()
	}
}
