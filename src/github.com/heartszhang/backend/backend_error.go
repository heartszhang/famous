package backend

type backend_error struct {
	Reason string `json:"reason,omitempty"`
	Code   int    `json:"code"`
}

func (this backend_error) Error() string{
	return this.Reason
}
