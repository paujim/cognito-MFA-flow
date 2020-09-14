package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockedCognitoClient struct {
	mock.Mock
}

func (m *mockedCognitoClient) InitiateAuth(input *cognitoidentityprovider.InitiateAuthInput) (*cognitoidentityprovider.InitiateAuthOutput, error) {
	args := m.Called(input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*cognitoidentityprovider.InitiateAuthOutput), args.Error(1)
}

func (m *mockedCognitoClient) RespondToAuthChallenge(input *cognitoidentityprovider.RespondToAuthChallengeInput) (*cognitoidentityprovider.RespondToAuthChallengeOutput, error) {
	args := m.Called(input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*cognitoidentityprovider.RespondToAuthChallengeOutput), args.Error(1)
}

func (m *mockedCognitoClient) AssociateSoftwareToken(input *cognitoidentityprovider.AssociateSoftwareTokenInput) (*cognitoidentityprovider.AssociateSoftwareTokenOutput, error) {
	args := m.Called(input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*cognitoidentityprovider.AssociateSoftwareTokenOutput), args.Error(1)
}

func (m *mockedCognitoClient) SetUserMFAPreference(input *cognitoidentityprovider.SetUserMFAPreferenceInput) (*cognitoidentityprovider.SetUserMFAPreferenceOutput, error) {
	args := m.Called(input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*cognitoidentityprovider.SetUserMFAPreferenceOutput), args.Error(1)
}

func (m *mockedCognitoClient) VerifySoftwareToken(input *cognitoidentityprovider.VerifySoftwareTokenInput) (*cognitoidentityprovider.VerifySoftwareTokenOutput, error) {
	args := m.Called(input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*cognitoidentityprovider.VerifySoftwareTokenOutput), args.Error(1)
}

func TestPingURL(t *testing.T) {
	t.Run("Successful ping", func(t *testing.T) {
		assert := assert.New(t)
		mockClient := &mockedCognitoClient{}
		ts := httptest.NewServer(createApp("userPool", "client", mockClient).Router)
		defer ts.Close()
		resp, err := http.Get(fmt.Sprintf("%s/v1/ping", ts.URL))
		assert.NoError(err)
		assert.Equal(resp.StatusCode, 200)
	})
}

func TestTokenURL(t *testing.T) {

	t.Run("Successful token", func(t *testing.T) {
		assert := assert.New(t)
		mockClient := &mockedCognitoClient{}
		mockClient.On("InitiateAuth", mock.Anything).Return(&cognitoidentityprovider.InitiateAuthOutput{
			AuthenticationResult: &cognitoidentityprovider.AuthenticationResultType{
				AccessToken: aws.String("ACCESS_TOKEN"),
			},
		}, nil)

		ts := httptest.NewServer(createApp("userPool", "client", mockClient).Router)
		defer ts.Close()

		resp, err := http.Post(fmt.Sprintf("%s/v1/token/", ts.URL), "application/json", strings.NewReader(`{"username": "hugo","password": "hugo@password"}`))
		assert.NoError(err)

		assert.Equal(resp.StatusCode, 200)
		mockClient.AssertExpectations(t)
	})
	t.Run("Validating token request", func(t *testing.T) {
		assert := assert.New(t)
		mockClient := &mockedCognitoClient{}
		ts := httptest.NewServer(createApp("userPool", "client", mockClient).Router)
		defer ts.Close()

		resp, err := http.Post(fmt.Sprintf("%s/v1/token/", ts.URL), "application/json", strings.NewReader(`{"value1": "hugo","value2": "hugo@password"}`))
		assert.NoError(err)
		defer resp.Body.Close()
		target := &TokenResponse{}
		json.NewDecoder(resp.Body).Decode(target)

		assert.Equal(http.StatusBadRequest, resp.StatusCode)
		assert.Equal("Missing required parameter", *target.Message)

		mockClient.AssertExpectations(t)
	})
	t.Run("Error when authenticating token request", func(t *testing.T) {
		assert := assert.New(t)
		mockClient := &mockedCognitoClient{}
		mockClient.On("InitiateAuth", mock.Anything).Return(nil, errors.New("ERROR with the password"))

		ts := httptest.NewServer(createApp("userPool", "client", mockClient).Router)
		defer ts.Close()

		resp, err := http.Post(fmt.Sprintf("%s/v1/token/", ts.URL), "application/json", strings.NewReader(`{"username": "hugo","password": "hugo@password"}`))
		assert.NoError(err)
		assert.Equal(resp.StatusCode, 401)
		mockClient.AssertExpectations(t)
	})
	t.Run("Successful update", func(t *testing.T) {
		assert := assert.New(t)
		mockClient := &mockedCognitoClient{}
		mockClient.On("RespondToAuthChallenge", mock.Anything).Return(&cognitoidentityprovider.RespondToAuthChallengeOutput{
			AuthenticationResult: &cognitoidentityprovider.AuthenticationResultType{
				AccessToken: aws.String("ACCESS_TOKEN"),
			},
		}, nil)

		ts := httptest.NewServer(createApp("userPool", "client", mockClient).Router)
		defer ts.Close()

		resp, err := http.Post(fmt.Sprintf("%s/v1/token/update", ts.URL), "application/json", strings.NewReader(`{"username": "hugo","password": "hugo@password", "session": "SESSION"}`))
		assert.NoError(err)
		assert.Equal(resp.StatusCode, 200)
		mockClient.AssertExpectations(t)
	})
	t.Run("Validating update request", func(t *testing.T) {
		assert := assert.New(t)
		mockClient := &mockedCognitoClient{}
		ts := httptest.NewServer(createApp("userPool", "client", mockClient).Router)
		defer ts.Close()

		resp, err := http.Post(fmt.Sprintf("%s/v1/token/update", ts.URL), "application/json", strings.NewReader(`{"value1": "hugo"}`))
		assert.NoError(err)
		defer resp.Body.Close()
		target := &TokenResponse{}
		json.NewDecoder(resp.Body).Decode(target)

		assert.Equal(http.StatusBadRequest, resp.StatusCode)
		assert.Equal("Missing required parameter", *target.Message)
		mockClient.AssertExpectations(t)
	})
	t.Run("Successful code", func(t *testing.T) {
		assert := assert.New(t)
		mockClient := &mockedCognitoClient{}
		mockClient.On("RespondToAuthChallenge", mock.Anything).Return(&cognitoidentityprovider.RespondToAuthChallengeOutput{
			AuthenticationResult: &cognitoidentityprovider.AuthenticationResultType{
				AccessToken: aws.String("ACCESS_TOKEN"),
			},
		}, nil)

		ts := httptest.NewServer(createApp("userPool", "client", mockClient).Router)
		defer ts.Close()

		resp, err := http.Post(fmt.Sprintf("%s/v1/token/code", ts.URL), "application/json", strings.NewReader(`{"username": "hugo","code": "123456", "session": "SESSION"}`))
		assert.NoError(err)
		assert.Equal(resp.StatusCode, 200)
		mockClient.AssertExpectations(t)
	})
	t.Run("Validating code request", func(t *testing.T) {
		assert := assert.New(t)
		mockClient := &mockedCognitoClient{}
		ts := httptest.NewServer(createApp("userPool", "client", mockClient).Router)
		defer ts.Close()

		resp, err := http.Post(fmt.Sprintf("%s/v1/token/code", ts.URL), "application/json", strings.NewReader(`{"value1": "hugo"}`))
		assert.NoError(err)

		defer resp.Body.Close()
		target := &TokenResponse{}
		json.NewDecoder(resp.Body).Decode(target)

		assert.Equal(http.StatusBadRequest, resp.StatusCode)
		assert.Equal("Missing required parameter", *target.Message)
		mockClient.AssertExpectations(t)
	})
}

func TestMFAURL(t *testing.T) {

	t.Run("Success registering MFA", func(t *testing.T) {
		assert := assert.New(t)
		mockClient := &mockedCognitoClient{}
		mockClient.On("AssociateSoftwareToken", mock.Anything).Return(&cognitoidentityprovider.AssociateSoftwareTokenOutput{
			SecretCode: aws.String("SECRET_CODE_GENERATED"),
		}, nil)
		ts := httptest.NewServer(createApp("userPool", "client", mockClient).Router)
		defer ts.Close()

		resp, err := http.Post(fmt.Sprintf("%s/v1/mfa/register", ts.URL), "application/json", strings.NewReader(`{"accessToken": "ACCESS_TOKEN"}`))
		assert.NoError(err)
		assert.Equal(http.StatusOK, resp.StatusCode)
		mockClient.AssertExpectations(t)
	})
	t.Run("Validating registering MFA", func(t *testing.T) {
		assert := assert.New(t)
		mockClient := &mockedCognitoClient{}
		ts := httptest.NewServer(createApp("userPool", "client", mockClient).Router)
		defer ts.Close()

		resp, err := http.Post(fmt.Sprintf("%s/v1/mfa/register", ts.URL), "application/json", strings.NewReader(`{"value": "misssing"}`))
		assert.NoError(err)

		defer resp.Body.Close()
		target := &MFAResponse{}
		json.NewDecoder(resp.Body).Decode(target)

		assert.Equal(http.StatusBadRequest, resp.StatusCode)
		assert.Equal("Missing required parameter", *target.Message)
		mockClient.AssertExpectations(t)
	})

	t.Run("Success enabling MFA", func(t *testing.T) {
		assert := assert.New(t)
		mockClient := &mockedCognitoClient{}
		mockClient.On("SetUserMFAPreference", mock.Anything).Return(&cognitoidentityprovider.SetUserMFAPreferenceOutput{}, nil)
		ts := httptest.NewServer(createApp("userPool", "client", mockClient).Router)
		defer ts.Close()

		resp, err := http.Post(fmt.Sprintf("%s/v1/mfa/enable", ts.URL), "application/json", strings.NewReader(`{"accessToken": "ACCESS_TOKEN"}`))
		assert.NoError(err)
		assert.Equal(http.StatusOK, resp.StatusCode)
		mockClient.AssertExpectations(t)
	})
	t.Run("Validating enabling MFA", func(t *testing.T) {
		assert := assert.New(t)
		mockClient := &mockedCognitoClient{}
		ts := httptest.NewServer(createApp("userPool", "client", mockClient).Router)
		defer ts.Close()

		resp, err := http.Post(fmt.Sprintf("%s/v1/mfa/enable", ts.URL), "application/json", strings.NewReader(`{"value": "misssing"}`))
		assert.NoError(err)

		defer resp.Body.Close()
		target := &MFAResponse{}
		json.NewDecoder(resp.Body).Decode(target)

		assert.Equal(http.StatusBadRequest, resp.StatusCode)
		assert.Equal("Missing required parameter", *target.Message)
		mockClient.AssertExpectations(t)
	})

	t.Run("Verifing MFA successfully", func(t *testing.T) {
		assert := assert.New(t)
		mockClient := &mockedCognitoClient{}
		mockClient.On("VerifySoftwareToken", mock.Anything).Return(&cognitoidentityprovider.VerifySoftwareTokenOutput{Status: aws.String("SUCCESS")}, nil)
		ts := httptest.NewServer(createApp("userPool", "client", mockClient).Router)
		defer ts.Close()

		resp, err := http.Post(fmt.Sprintf("%s/v1/mfa/verify", ts.URL), "application/json", strings.NewReader(`{"accessToken": "ACCESS_TOKEN", "code": "123456"}`))
		assert.NoError(err)
		assert.Equal(http.StatusOK, resp.StatusCode)
		mockClient.AssertExpectations(t)
	})
	t.Run("Failing to verify MFA", func(t *testing.T) {
		assert := assert.New(t)
		mockClient := &mockedCognitoClient{}
		mockClient.On("VerifySoftwareToken", mock.Anything).Return(nil, errors.New("WRONG CODE"))
		ts := httptest.NewServer(createApp("userPool", "client", mockClient).Router)
		defer ts.Close()

		resp, err := http.Post(fmt.Sprintf("%s/v1/mfa/verify", ts.URL), "application/json", strings.NewReader(`{"accessToken": "ACCESS_TOKEN", "code": "123456"}`))
		assert.NoError(err)

		defer resp.Body.Close()
		target := &MFAResponse{}
		json.NewDecoder(resp.Body).Decode(target)

		assert.Equal(http.StatusBadRequest, resp.StatusCode)
		assert.Equal("Unable to verify code", *target.Message)

		mockClient.AssertExpectations(t)
	})
	t.Run("Validating verifing MFA", func(t *testing.T) {
		assert := assert.New(t)
		mockClient := &mockedCognitoClient{}
		ts := httptest.NewServer(createApp("userPool", "client", mockClient).Router)
		defer ts.Close()

		resp, err := http.Post(fmt.Sprintf("%s/v1/mfa/verify", ts.URL), "application/json", strings.NewReader(`{"value": "misssing"}`))
		assert.NoError(err)

		defer resp.Body.Close()
		target := &MFAResponse{}
		json.NewDecoder(resp.Body).Decode(target)

		assert.Equal(http.StatusBadRequest, resp.StatusCode)
		assert.Equal("Missing required parameter", *target.Message)
		mockClient.AssertExpectations(t)
	})
}
