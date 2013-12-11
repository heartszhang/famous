package backend

type backend_error struct {
	Reason string `json:"reason,omitempty"`
	Code   int    `json:"code"`
}

func (this backend_error) Error() string {
	return this.Reason
}
func new_backenderror(code int, reason string) backend_error {
	return backend_error{Code: code, Reason: reason}
}
