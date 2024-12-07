package apperrors

var (
	ErrInvalidDateLayout = Error{
		Message: "invalid date layout",
	}
	ErrInvalidDateFormat = Error{
		Message: "invalid date format",
	}
	ErrInvalidDateRange = Error{
		Message: "invalid date range",
	}
)

type Error struct {
	Message string
}

func (e Error) Error() string {
	return e.Message
}
