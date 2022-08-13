package responses

type RefreshResponse struct {
	Err          string `json:"error"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
