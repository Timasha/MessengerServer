package requests

type RegistrationRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
