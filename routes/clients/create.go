package clients

import (
	"fmt"
	"microservice/internal/db"
	apiErrors "microservice/internal/errors"
	"microservice/resources"
	"microservice/types"
	"microservice/utils"
	"net/http"
	"slices"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

func Create(c *gin.Context) {
	var parameters struct {
		Description  string   `json:"description" binding:"required"`
		ContactName  string   `json:"contactName" binding:"required"`
		ContactEmail string   `json:"contactEMail" binding:"required"`
		Scopes       []string `json:"scopes" binding:"required"`
	}

	err := c.BindJSON(&parameters)
	if err != nil {
		c.Abort()
		res := apiErrors.ErrMissingParameter
		res.Errors = []error{err}
		res.Emit(c)
		return
	}

	if slices.Contains(parameters.Scopes, "*:*") {
		c.Abort()
		apiErrors.ErrInvalidClientScopeRequested.Emit(c)
		return
	}

	userSubject := c.GetString("subject")
	user, err := utils.GetUser(types.InternalIdentifier(userSubject))
	if err != nil {
		c.Abort()
		_ = c.Error(err)
		return
	}

	userPermissions := make([]string, 0)
	for system, scopes := range user.Permissions() {
		for _, scope := range scopes {
			scopeString := fmt.Sprintf("%s:%s", system, scope)
			userPermissions = append(userPermissions, scopeString)
		}
	}

	for _, requestedScope := range parameters.Scopes {
		if !slices.Contains(userPermissions, requestedScope) {
			c.Abort()
			apiErrors.ErrPermissionMismatch.Emit(c)
			return
		}
	}

	query, err := db.Queries.Raw("get-services")
	if err != nil {
		panic(err)
	}
	var services []types.Service
	err = pgxscan.Select(c, db.Pool, &services, query)
	if err != nil {
		panic(err)
	}

	availableScopes := make([]string, 0)
	for _, service := range services {
		for _, scope := range service.SupportedScopes {
			availableScopes = append(availableScopes, fmt.Sprintf("%s:%s", service.Name, scope))
		}
	}

	for _, requestedScope := range parameters.Scopes {
		if !slices.Contains(availableScopes, requestedScope) {
			c.Abort()
			apiErrors.ErrInvalidClientScopeRequested.Emit(c)
			return
		}
	}

	query, err = db.Queries.Raw("create-client")
	if err != nil {
		c.Abort()
		_ = c.Error(err)
		return
	}

	var clientID string
	err = pgxscan.Get(c, db.Pool, &clientID, query, parameters.Description, parameters.ContactName, parameters.ContactEmail)
	if err != nil {
		c.Abort()
		_ = c.Error(err)
		return
	}

	b := jwt.NewBuilder()
	b.Issuer("user-management")
	b.IssuedAt(time.Now())
	b.Subject(clientID)
	b.Audience([]string{"user-management"})
	b.Claim("scopes", parameters.Scopes)

	clientToken, err := b.Build()
	if err != nil {
		c.Abort()
		_ = c.Error(err)
		return
	}

	s := jwt.NewSerializer()
	s.Sign(jwt.WithKey(resources.PrivateSigningKey.Algorithm(), resources.PrivateSigningKey))
	s.Encrypt(jwt.WithKey(jwa.ECDH_ES, resources.PublicEncryptionKey))

	clientSecret, err := s.Serialize(clientToken)
	if err != nil {
		c.Abort()
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"clientID":     clientID,
		"clientSecret": string(clientSecret),
	})

}
