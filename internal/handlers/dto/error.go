package dto

type ErrorDto struct {
	Message string `json:"message"` // Текст ошибки
}

var (
	UnauthorizedError     = ErrorDto{Message: "unauthorized"}
	MissingOrInvalidToken = ErrorDto{Message: "missing or invalid token"}
	ForbiddenError        = ErrorDto{Message: "forbidden: insufficient permissions"}
	NoUserIDProvided      = ErrorDto{Message: "no user_id provided"}
)
