package authErrors

func IsNotFoundError(err error) bool {
	switch err.(type) {
	case NoSuchUser:
		return true
	default:
		return false
	}
}
