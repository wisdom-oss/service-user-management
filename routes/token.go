package routes

import (
	"fmt"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/thanhpk/randstr"
	"golang.org/x/oauth2"

	"microservice/interfaces"
	"microservice/internal/db"
	"microservice/oidc"
	"microservice/structs"
	"microservice/utils"
)

type TokenRequest struct {
	GrantType    string `json:"grant_type" form:"grant_type" binding:"required"`
	Code         string `json:"code" form:"code"`
	State        string `json:"state" form:"state"`
	ClientID     string `json:"client_id" form:"client_id"`
	ClientSecret string `json:"client_secret" form:"client_secret"`
	RefreshToken string `json:"refresh_token" form:"refresh_token"`
}

func Token(c *gin.Context) {
	var tokenRequest TokenRequest
	if err := c.ShouldBind(&tokenRequest); err != nil {
		c.Abort()
		_ = c.Error(err)
		return
	}

	var user interfaces.PermissionableObject
	switch tokenRequest.GrantType {
	case "client_credentials":
		// TODO: Handle client credentials
		user = checkClientCredentials(c, tokenRequest)
	case "authorization_code":
		user = exchangeAuthorizationCode(c, tokenRequest)
	case "refresh_token":
		user = identifyUserFromRefreshToken(c, tokenRequest)
	}

	if user == nil {
		c.Abort()
		c.Status(500)
		return
	}

	var permissions []string
	for system, scopes := range user.Permissions() {
		for _, scope := range scopes {
			scopeString := fmt.Sprintf("%s:%s", system, scope)
			permissions = append(permissions, scopeString)
		}
	}

	tokenBuilder := jwt.NewBuilder()
	tokenBuilder.Expiration(time.Now().Add(time.Minute * 120))
	tokenBuilder.NotBefore(time.Now())
	tokenBuilder.Subject(user.GetID())
	tokenBuilder.Audience([]string{"wisdom"})
	tokenBuilder.Issuer("user-management")
	tokenBuilder.Claim("scopes", permissions)

	token, err := tokenBuilder.Build()
	if err != nil {
		c.Abort()
		_ = c.Error(err)
		return
	}

	serializer := jwt.NewSerializer()
	serializer.Sign(jwt.WithInsecureNoSignature())
	serializedToken, err := serializer.Serialize(token)
	if err != nil {
		c.Abort()
		_ = c.Error(err)
		return
	}

	refreshTokenBuilder := jwt.NewBuilder()
	refreshTokenBuilder.Expiration(time.Now().Add(time.Hour * 12))
	refreshTokenBuilder.NotBefore(time.Now())
	refreshTokenBuilder.Subject(user.GetID())
	refreshTokenBuilder.Issuer("user-management")
	refreshTokenBuilder.Audience([]string{"user-management"})
	refreshTokenBuilder.Claim("scopes", permissions)
	refreshTokenBuilder.JwtID(randstr.Base62(128))
	refreshToken, err := refreshTokenBuilder.Build()
	if err != nil {
		c.Abort()
		_ = c.Error(err)
		return
	}

	// TODO: reactivate refresh token if storing and validation are successful
	_, err = serializer.Serialize(refreshToken)
	if err != nil {
		c.Abort()
		_ = c.Error(err)
		return
	}

	res := structs.TokenResponse{
		AccessToken: string(serializedToken),
		ExpiresIn:   int(token.Expiration().Sub(time.Now()).Seconds()),
		TokenType:   "Bearer",
		//RefreshToken: string(serializedRefreshToken),
	}

	c.JSON(200, res)

}

func checkClientCredentials(c *gin.Context, tokenRequest TokenRequest) interfaces.PermissionableObject {
	return nil
}

func exchangeAuthorizationCode(c *gin.Context, tokenRequest TokenRequest) interfaces.PermissionableObject {
	// retrieve the verifier from the database
	var verifier string
	err := db.StateDB.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(tokenRequest.State))
		_ = item.Value(func(val []byte) error {
			verifier = string(val)
			return nil
		})
		return err
	})
	if err != nil {
		c.Abort()
		_ = c.Error(err)
		return nil
	}

	// now exchange the code for a token
	token, err := oidc.ExternalProvider.Exchange(c, tokenRequest.Code, oauth2.VerifierOption(verifier), oauth2.SetAuthURLParam("state", tokenRequest.State))
	if err != nil {
		c.Abort()
		_ = c.Error(err)
		return nil
	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		c.Abort()
		c.Status(500)
		return nil
	}

	idToken, err := oidc.TokenVerifier.Verify(c, rawIDToken)
	if err != nil {
		c.Abort()
		_ = c.Error(err)
		return nil
	}

	user, err := utils.GetUser(idToken.Subject)
	if err != nil {
		if err == utils.ErrNoUser {
			user, err = utils.CreateUser(token)
		}
	}

	return user
}
