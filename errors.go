package googlemaps

import (
	"fmt"
)

// APIError represents a GoogleMaps API Error response
type APIError struct {
	ErrorMessage string        `json:"error_message"`
	Results      []ErrorDetail `json:"results"`
	Status       string        `json:"status"`
}

// ErrorDetail represents an individual item in an APIError.
type ErrorDetail struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func (e APIError) Error() string {
	if len(e.Results) > 0 {
		err := e.Results[0]
		return fmt.Sprintf("GoogleMaps: %d %v", err.Code, err.Message)
	}
	return ""
}

// Empty returns true if empty. Otherwise, at least 1 error message/code is
// present and false is returned.
func (e APIError) Empty() bool {
	if len(e.Results) == 0 {
		return true
	}
	return false
}

// relevantError returns any non-nil http-related error (creating the request,
// getting the response, decoding) if any. If the decoded apiError is non-zero
// the apiError is returned. Otherwise, no errors occurred, returns nil.
func relevantError(httpError error, apiError APIError) error {
	if httpError != nil {
		return httpError
	}
	if apiError.Empty() {
		return nil
	}
	return apiError
}
