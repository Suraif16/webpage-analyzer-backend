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
        Message:     "Invalid URL Format",
        Description: "The URL provided is invalid or malformed. Please ensure it starts with http:// or https://.",
    }

    ErrPageNotFound = &APIError{
        StatusCode:  404,
        Message:     "Resource Not Found",
        Description: "The requested resource could not be found. Please check the URL.",
    }

    ErrDNSResolutionFailed = &APIError{
        StatusCode:  502,
        Message:     "DNS Resolution Failed",
        Description: "The domain could not be resolved. Please check if the URL is correct.",
    }
	ErrPageNotAccessible = &APIError{
		StatusCode:  503,
		Message:     "Service Unavailable",
		Description: "The target server is not responding or is temporarily unavailable",
	}
	
    ErrTimeout = &APIError{
        StatusCode:  504,
        Message:     "Request Timeout",
        Description: "The website took too long to respond. Please try again later.",
    }

    ErrInternalServer = &APIError{
		StatusCode:  500,
		Message:     "Internal Server Error",
		Description: "An unexpected error occurred. Please try again later.",
	}
)
