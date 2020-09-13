package main

import (
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

func TestServer(t *testing.T) {
	t.Run("Successful ping", func(t *testing.T) {
		assert := assert.New(t)
		mockClient := &mockedCognitoClient{}
		ts := httptest.NewServer(createApp("userPool", "client", mockClient).Router)
		defer ts.Close()
		resp, err := http.Get(fmt.Sprintf("%s/v1/ping", ts.URL))
		assert.NoError(err)
		assert.Equal(resp.StatusCode, 200)
	})

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
	t.Run("Validating username and password", func(t *testing.T) {
		assert := assert.New(t)
		mockClient := &mockedCognitoClient{}
		ts := httptest.NewServer(createApp("userPool", "client", mockClient).Router)
		defer ts.Close()

		resp, err := http.Post(fmt.Sprintf("%s/v1/token/", ts.URL), "application/json", strings.NewReader(`{"value1": "hugo","value2": "hugo@password"}`))
		assert.NoError(err)
		assert.Equal(resp.StatusCode, 400)
		mockClient.AssertExpectations(t)
	})
	t.Run("Error when authenticating", func(t *testing.T) {
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
}
