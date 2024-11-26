package clients

import (
	"microservice/internal/db"
	apiErrors "microservice/internal/errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Delete(c *gin.Context) {
	clientID := c.Param("clientID")

	_, err := uuid.Parse(clientID)
	if err != nil {
		c.Abort()
		apiErrors.ErrInvalidClientID.Emit(c)
		return
	}

	query, err := db.Queries.Raw("delete-client")
	if err != nil {
		c.Abort()
		_ = c.Error(err)
		return
	}

	_, err = db.Pool.Exec(c, query, clientID)
	if err != nil {
		c.Abort()
		_ = c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}
