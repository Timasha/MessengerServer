package responses

// Объект ответа для сериализации в JSON. Используется в хэндлере авторизации.
type AuthorizationResponse struct {
	Err          string `json:"error"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
