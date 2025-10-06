package core

import "net/http"

func ValidateRequestMethod(r *http.Request, allowedMethod string) *ServerError {
	if r.Method != allowedMethod {
		return &ServerError{
			Code:    http.StatusMethodNotAllowed,
			Message: ErrMethodNotAllowed.Error(),
			Err:     nil,
		}
	}
	return nil
}
