package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thanhpk/randstr"
	"golang.org/x/oauth2"

	"microservice/internal/db"
	"microservice/internal/errors"
	"microservice/oidc"
)

func InitiateLogin(c *gin.Context) {
	var parameters struct {
		RedirectUri string `form:"redirect_uri" binding:"required"`
	}
	err := c.ShouldBindQuery(&parameters)
	if err != nil {
		c.Abort()
		res := errors.ErrMissingParameter
		res.Errors = []error{err}
		res.Emit(c)
	}
	// generate a new state for this login
	state := randstr.Base62(32)
	verifier := randstr.Base62(128)
	err = db.Redis.Set(c, state, verifier, 0).Err()

	if err != nil {
		c.Abort()
		_ = c.Error(err)
		return
	}
	c.Redirect(http.StatusFound, oidc.ExternalProvider.AuthCodeURL(state, oauth2.S256ChallengeOption(verifier), oauth2.SetAuthURLParam("redirect_uri", parameters.RedirectUri)))
}
