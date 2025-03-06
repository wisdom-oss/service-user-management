package routes

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/thanhpk/randstr"
	"golang.org/x/oauth2"

	"microservice/internal/db"
	"microservice/internal/errors"
	"microservice/oidc"
	"microservice/types"
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
	tokenParams := types.LoginParameters{}
	tokenParams.RedirectUri = parameters.RedirectUri
	tokenParams.CodeVerifier = randstr.Base62(128)

	params, _ := json.Marshal(tokenParams)
	err = db.Redis.Set(c, state, params, 5*time.Minute).Err()

	if err != nil {
		c.Abort()
		_ = c.Error(err)
		return
	}
	c.Redirect(http.StatusFound, oidc.ExternalProvider.AuthCodeURL(state, oauth2.S256ChallengeOption(tokenParams.CodeVerifier), oauth2.SetAuthURLParam("redirect_uri", parameters.RedirectUri)))
}
