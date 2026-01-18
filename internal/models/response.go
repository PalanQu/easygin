package models

type APIResponse struct {
	Code      interface{} `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	RequestID string      `json:"request_id,omitempty"`
	Details   interface{} `json:"details,omitempty"`
}

func SuccessResponse(data interface{}, requestID string) *APIResponse {
	return &APIResponse{
		Code:      0,
		Message:   "success",
		Data:      data,
		RequestID: requestID,
	}
}

func ErrorResponse(code interface{}, message string, requestID string) *APIResponse {
	return &APIResponse{
		Code:      code,
		Message:   message,
		Data:      nil,
		RequestID: requestID,
	}
}

func ErrorResponseWithDetails(code interface{}, message string, requestID string, details interface{}) *APIResponse {
	return &APIResponse{
		Code:      code,
		Message:   message,
		Data:      nil,
		RequestID: requestID,
		Details:   details,
	}
}
