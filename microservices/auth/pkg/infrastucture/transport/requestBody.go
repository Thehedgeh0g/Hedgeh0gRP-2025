package transport

type loginBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type registerBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
