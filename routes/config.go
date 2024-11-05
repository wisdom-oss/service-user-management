package routes

import (
	"github.com/gin-gonic/gin"

	"microservice/oidc"
)

func Configuration(c *gin.Context) {
	res := struct {
		ClientID    string `json:"client_id"`
		Issuer      string `json:"issuer"`
		RedirectUri string `json:"redirect_uri"`
	}{
		ClientID:    oidc.ExternalProvider.ClientID,
		Issuer:      oidc.ExternalIssuer,
		RedirectUri: oidc.ExternalProvider.RedirectURL,
	}
	c.JSON(200, res)
}
