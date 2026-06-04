package model

// LoginRequest contiene las credenciales para un inicio de sesión básico.
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse devuelve el estado de un inicio de sesión exitoso.
type LoginResponse struct {
	Username string `json:"username"`
	Token    string `json:"token"`
}
