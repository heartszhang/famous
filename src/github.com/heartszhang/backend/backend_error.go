package backend

type BackendError struct {
	Reason string `json:"reason,omitempty"`
	Code   int    `json:"code"`
}
