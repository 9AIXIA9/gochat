package api

var (
	ResponseSuccess            = NewResponse(CodeSuccess)
	ResponseServerError        = NewResponse(CodeServerError)
	ResponseTimeout            = NewResponse(CodeTimeout)
	ResponseInvalidToken       = NewResponse(CodeUnauthorized)
	ResponseInvalidParam       = NewResponse(CodeInvalidParam)
	ResponseServiceUnavailable = NewResponse(CodeServiceUnavailable)
	ResponseNotFound           = NewResponse(CodeNotFound)
	ResponseBusinessError      = NewResponse(CodeBusinessError)
)
