package apperrors

var (
	ErrInvalidDateLayout = UserError{
		Message: "invalid date layout",
	}
	ErrInvalidDateFormat = UserError{
		Message: "invalid date format",
	}
	ErrInvalidDateRange = UserError{
		Message: "invalid date range",
	}

	ErrInvalidFilters = UserError{
		Message: "invalid filters",
	}

	ErrBadDetailsRequest = UserError{
		Message: "details for song not found",
	}
	ErrDetailsServerBadResponse = AppError{
		Message: "details server send bad response",
	}
)

type UserError struct {
	Message string
}

func (e UserError) Error() string {
	return e.Message
}

type AppError struct {
	Message string
}

func (e AppError) Error() string {
	return e.Message
}
