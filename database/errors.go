package database

import "fmt"

type DriverError struct {
	Status  int    `json:"status`
	Message string `json:"message"`
}

func (i DriverError) Error() string {
	return fmt.Sprintf("%i error: %s", i.Status, i.Message)
}