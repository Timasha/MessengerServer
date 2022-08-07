package requests

type AuthorizeRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
