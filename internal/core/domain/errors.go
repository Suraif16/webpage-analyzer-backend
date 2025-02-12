package domain

type APIError struct {
	StatusCode  int    `json:"statusCode"`
	Message     string `json:"message"`
	Description string `json:"description,omitempty"`
}

func (e *APIError) Error() string {
	return e.Message
}

var (
	ErrInvalidURL = &APIError{
		StatusCode:  400,
		Message:     "Bad Request",
		Description: "The provided URL is invalid or malformed",
	}
	ErrPageNotFound = &APIError{
		StatusCode:  404,
		Message:     "Not Found",
		Description: "The requested page could not be found",
	}
	ErrPageNotAccessible = &APIError{
		StatusCode:  503,
		Message:     "Service Unavailable",
		Description: "The target server is not responding or is unavailable",
	}
	ErrTimeout = &APIError{
		StatusCode:  504,
		Message:     "Gateway Timeout",
		Description: "The request timed out while trying to access the page",
	}
	ErrInternalServer = &APIError{
		StatusCode:  500,
		Message:     "Internal Server Error",
		Description: "An unexpected error occurred while processing the request",
	}
)
