package errors

type NotAuthorizedError struct {
}

func (n NotAuthorizedError) Error() string {
	return "channel is not authorized"
}
