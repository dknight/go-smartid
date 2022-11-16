package smartid

import (
	"fmt"
	"net/http"
)

const internalResponseError = -1

// Response is a basic response structure that holds code (usually HTTP
// status) and a message (usually HTTP status text). It is typically only
// used for HTTP code for responses.
type Response struct {
	// Code is HTTP status or internal error code.
	Code int `json:"code"`

	// Messages is HTTP status text or internal message.
	Message string `json:"message"`
}

// IsStatusOK checks if response has HTTP Status 200 (OK).
func (r *Response) IsStatusOK() bool {
	return r.Code == http.StatusOK
}

// ResponseError represents error for the response.
type ResponseError Response

func (e *ResponseError) Error() string {
	return fmt.Sprintf("Response error: %v %v", e.Code, e.Message)
}
