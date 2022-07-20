package gtypes

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Type to manage the signup service. JSON bindings ensure all required fields are provided and the data inputs are in correct format
type Signup struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required"`
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
}

// Type to manage the login process. JSON bindings manage user input and format requirements
type Login struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Type that manages the update of a user record
type UserRecord struct {
	Email     string `json:"email" binding:"omitempty"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

// Function that logs a fatal error
func Fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// Function that evaluates a gin response and returns HTTP 500 if an error exists
func EvalError(c *gin.Context, err error) bool {
	if err != nil {
		log.Println(err)
		c.String(http.StatusInternalServerError, "%+v", err)
		return true
	}
	return false
}
