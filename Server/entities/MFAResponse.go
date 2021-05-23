package entities

type MFAResponse struct {
	Message            *string `json:"message"`
	SecretCode         *string `json:"secret"`
	GoogleAutheticator *string `json:"googleAutheticator"`
}
