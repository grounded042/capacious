package utils

const recordNotFound = "record not found"

// GetCodeForError returns the error code related ot the error err passed in.
func GetCodeForError(err error) int {
	switch err.Error() {
	case recordNotFound:
		return 404
	}
	return 500
}
