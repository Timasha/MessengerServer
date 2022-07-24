package responses

// Объект ответа для сериализации в JSON. Используется в хэндлере авторизации.
type AuthorizeResponce struct {
	Token string `json:"token"`
	Err   string `json:"error"`
}
