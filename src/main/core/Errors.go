package core

type AppError struct {
	Message string
}

func (a AppError) Error() string {
	return a.Message
}


