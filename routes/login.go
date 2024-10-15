package routes

import (
	"net/http"

	"github.com/dgraph-io/badger/v4"
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
	go func(state, verifier string) {
		_ = db.StateDB.Update(func(txn *badger.Txn) error {
			return txn.Set([]byte(state), []byte(verifier))
		})
	}(state, verifier)
	c.Redirect(http.StatusFound, oidc.ExternalProvider.AuthCodeURL(state, oauth2.S256ChallengeOption(verifier)))
}
