package authErrors

func IsNotFoundError(err error) bool {
	switch err.(type) {
	case NoSuchUser:
		return true
	default:
		return false
	}
}

func IsIncorrectPasswordError(err error) bool {
	switch err.(type) {
	case IncorrectPassword:
		return true
	default:
		return false
	}
}
