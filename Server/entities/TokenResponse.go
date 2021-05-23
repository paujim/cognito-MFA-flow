package entities

type TokenResponse struct {
	Message      *string `json:"message"`
	Session      *string `json:"session"`
	AccessToken  *string `json:"accessToken"`
	RefreshToken *string `json:"refreshToken"`
}
