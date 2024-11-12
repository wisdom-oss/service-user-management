package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/jwx/v2"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwe"
	"github.com/lestrrat-go/jwx/v2/jwt"

	"microservice/internal/db"
	"microservice/internal/errors"
	"microservice/resources"
)

func RevokeToken(c *gin.Context) {
	var parameters struct {
		Token string `form:"token" binding:"required"`
	}

	if err := c.ShouldBind(&parameters); err != nil {
		c.Abort()
		res := errors.ErrMissingParameter
		res.Errors = []error{err}
		res.Emit(c)
		return
	}

	tokenType := jwx.GuessFormat([]byte(parameters.Token))
	switch tokenType {
	case jwx.JWE:
		decryptedRefreshToken, err := jwe.Decrypt(
			[]byte(parameters.Token),
			jwe.WithKey(jwa.ECDH_ES, resources.PrivateSigningKey),
		)
		if err != nil {
			c.Status(200)
			return
		}

		token, err := jwt.Parse(decryptedRefreshToken,
			jwt.WithIssuer("user-management"),
			jwt.WithVerify(true),
			jwt.WithKey(resources.PublicSigningKey.Algorithm(), resources.PublicSigningKey),
		)
		if err != nil {
			c.Status(200)
			return
		}

		query, err := db.Queries.Raw("revoke-refresh-token")
		if err != nil {
			c.Status(200)
			return
		}

		_, err = db.Pool.Exec(c, query, token.JwtID())
		if err != nil {
			c.Status(200)
			return
		}

		c.Status(200)
		return
	default:
		c.Status(200)
	}

}
