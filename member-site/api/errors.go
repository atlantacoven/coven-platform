package api

import "fmt"

type HTTPResponseError interface {
	Error() string
	Status() int
}

type httpError struct {
	cause   error
	status  int
	message string
}

func (err httpError) Error() string {
	if err.message != "" {
		return fmt.Errorf("[%v] %v: %w", err.status, err.message, err.cause).Error()
	}
	return fmt.Errorf("[%v]: %w", err.status, err.cause).Error()
}

func (err httpError) Status() int {
	return err.status
}
