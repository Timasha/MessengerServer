package responses

type RegistrationResponse struct {
	Err          string `json:"error"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
