package caedral

import "fmt"

// APIError is returned when the Caedral API responds with an error status.
type APIError struct {
	Message    string
	StatusCode int
	Type       string
	RawBody    any
}

// CaedralAPIError is an alias for APIError (Caedral API error response).
type CaedralAPIError = APIError

func (e *APIError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return fmt.Sprintf("caedral api error (status %d)", e.StatusCode)
}

// NewAPIError builds an APIError from an HTTP status and parsed JSON body.
func NewAPIError(statusCode int, body any) *APIError {
	err := &APIError{StatusCode: statusCode, RawBody: body}

	switch v := body.(type) {
	case map[string]any:
		if nested, ok := v["error"].(map[string]any); ok {
			if msg, ok := nested["message"].(string); ok {
				err.Message = msg
			}
			if t, ok := nested["type"].(string); ok {
				err.Type = t
			}
			if code, ok := nested["code"].(float64); ok && code != 0 {
				err.StatusCode = int(code)
			}
		} else if msg, ok := v["message"].(string); ok {
			err.Message = msg
		}
	case string:
		err.Message = v
	}

	if err.Message == "" {
		err.Message = fmt.Sprintf("request failed with status %d", statusCode)
	}
	if err.Type == "" {
		err.Type = "unknown"
	}

	return err
}
