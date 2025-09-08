package core

import "net/http"

func ValidateRequestMethod(r *http.Request, allowedMethod string) *ApiError {
	if r.Method != allowedMethod {
		return &ApiError{
			Code:    http.StatusMethodNotAllowed,
			Message: ErrMethodNotAllowed.Error(),
			Err:     nil,
		}
	}
	return nil
}
