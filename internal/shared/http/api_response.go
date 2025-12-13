package http

var (
	ResponseSuccess = NewApiResponse(CodeSuccess)
)

type ApiResponse struct {
	Code    BusinessCode `json:"code"`
	Message string       `json:"message"`
	Data    any          `json:"data,omitempty"`
}

func NewApiResponse(code BusinessCode) *ApiResponse {
	return &ApiResponse{
		Code:    code,
		Message: code.String(),
	}
}

func NewApiResponseWithMessage(code BusinessCode, message string) *ApiResponse {
	return &ApiResponse{
		Code:    code,
		Message: message,
	}
}

func NewApiResponseWithData(data any) *ApiResponse {
	return &ApiResponse{
		Code:    CodeSuccess,
		Message: CodeSuccess.String(),
		Data:    data,
	}
}
