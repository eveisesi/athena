package esi

type GenericError struct {
	Message string `json:"message"`
}

func (e GenericError) Error() string {
	return e.Message
}
