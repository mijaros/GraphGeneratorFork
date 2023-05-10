package api

import "strings"

func NewErr(httpError error, cause error) ErrorResponse {
	var causePtr *string
	if cause != nil {
		builder := strings.Builder{}
		builder.WriteString(cause.Error())
		causeStr := builder.String()
		causePtr = &causeStr
	}
	return ErrorResponse{
		Error: httpError.Error(),
		Cause: causePtr,
	}
}
