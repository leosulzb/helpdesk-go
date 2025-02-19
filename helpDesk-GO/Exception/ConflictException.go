package Exception

import "fmt"

type ConflictException struct {
	Message string
	Uri     string
}

func (e *ConflictException) Error() string {
	return fmt.Sprintf("Conflict: %s", e.Message)
}

type ForbiddenException struct {
	Message string
	Uri     string
}

func (e *ForbiddenException) Error() string {
	return fmt.Sprintf("Forbidden: %s", e.Message)
}
