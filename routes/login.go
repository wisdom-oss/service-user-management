package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thanhpk/randstr"
	"golang.org/x/oauth2"

	"microservice/internal/db"
	"microservice/oidc"
)

func InitiateLogin(c *gin.Context) {
	// generate a new state for this login
	state := randstr.Base62(32)
	verifier := randstr.Base62(128)
	err := db.Redis.Set(c, state, verifier, 0).Err()
	if err != nil {
		c.Abort()
		_ = c.Error(err)
		return
	}
	c.Redirect(http.StatusFound, oidc.ExternalProvider.AuthCodeURL(state, oauth2.S256ChallengeOption(verifier)))
}
