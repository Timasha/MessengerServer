package responses

type CheckTokenResponse struct {
	Err   string `json:"error"`
	Login string `json:"login"`
}
