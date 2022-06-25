package http

import "github.com/gin-gonic/gin"

type Error struct {
	Code    int     `json:"code"`
	Details []gin.H `json:"details,omitempty"`
	Message string  `json:"message"`
}

func NewError(code int, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

func (e *Error) WithDetails(details ...gin.H) *Error {
	e.Details = append(e.Details, details...)
	return e
}

func (e *Error) WriteJSON(ctx *gin.Context) {
	ctx.JSON(e.Code, gin.H{"error": e})
}
