package users

import (
	"net/http"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"

	"microservice/internal/db"
	"microservice/types"
)

func List(c *gin.Context) {
	query, err := db.Queries.Raw("get-users")
	if err != nil {
		c.Abort()
		_ = c.Error(err)
		return
	}

	var users []types.User
	err = pgxscan.Select(c, db.Pool, &users, query)
	if err != nil {
		if err == pgx.ErrNoRows {
			c.Status(http.StatusNoContent)
			return
		}
		c.Abort()
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, users)
}
