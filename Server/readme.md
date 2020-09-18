# A lambda backend
A gin server that uses aws cognito for authentication

## Creating a deployment package on Windows
To create a .zip that will work on AWS Lambda using Windows, install the build-lambda-zip tool.

```
go.exe get -u github.com/aws/aws-lambda-go/cmd/build-lambda-zip
```

Then run the following:

```
$env:GOOS = "linux"
$env:CGO_ENABLED = "0"
$env:GOARCH = "amd64"
go build -o main .
~\Go\Bin\build-lambda-zip.exe -output main.zip main
```
