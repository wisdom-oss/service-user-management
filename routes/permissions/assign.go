package permissions

import (
	"github.com/gin-gonic/gin"

	"microservice/internal/errors"
)

// Assign only creates new permission assignments and doesn't modify existing
// ones
func Assign(c *gin.Context) {
	var parameters struct {
		UserID      string `json:"user" binding:"required"`
		Assignments []struct {
			Service string `json:"service" binding:"required"`
			Scope   string `json:"scope" binding:"required"`
		} `json:"assignments" binding:"required"`
	}
	err := c.BindJSON(&parameters)
	if err != nil {
		res := errors.ErrMissingParameter
		res.Errors = []error{err}
		res.Emit(c)
		return
	}

}
