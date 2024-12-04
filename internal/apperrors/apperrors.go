package apperrors

type Error struct {
	Err     error
	Message string
}

func (e Error) Error() string {
	return e.Message
}
