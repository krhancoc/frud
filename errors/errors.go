//Package errors holds the error objects used by the Frud framework
package errors

import "fmt"

// DriverError are errors outputted by the driver object, it will encapsulate any errors
// outputted by their associated database.
type DriverError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// Error standard error to string function
func (i DriverError) Error() string {
	return fmt.Sprintf("%d error: %s", i.Status, i.Message)
}
