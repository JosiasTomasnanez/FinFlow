package lab

// AllowedErrorStatus solo 5xx: el stat de Grafana usa code=~"5..".
func AllowedErrorStatus(code int) int {
	switch code {
	case 500, 502, 503:
		return code
	default:
		return 500
	}
}
