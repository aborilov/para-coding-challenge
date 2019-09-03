package main

import "fmt"

type ErrNotFound struct {
	Email string
}

// Error implements the error interface for ErrNotFound.
func (e ErrNotFound) Error() string {
	return fmt.Sprintf("user with email %s not found ", e.Email)
}
