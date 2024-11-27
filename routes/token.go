package routes

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwe"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/thanhpk/randstr"
	"golang.org/x/oauth2"

	"microservice/interfaces"
	"microservice/internal/db"
	apiErrors "microservice/internal/errors"
	"microservice/oidc"
	"microservice/resources"
	"microservice/types"
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

var TokenAudiences = []string{"user-management", "wisdom"}
var RefreshTokenAudiences = []string{"wisdom"}

const TokenIssuer = "user-management"

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
		user = checkClientCredentials(c, tokenRequest) // todo: fully implement client credentials
	case "authorization_code":
		user = exchangeAuthorizationCode(c, tokenRequest)
	case "refresh_token":
		user = issueFromRefreshToken(c, tokenRequest)
	}

	if c.IsAborted() {
		return
	}

	if user == nil {
		c.Abort()
		c.Status(500)
		return
	}

	if !user.IsActive() {
		c.Abort()
		apiErrors.ErrUserDisabled.Emit(c)
		return
	}

	var permissions []string
	for system, scopes := range user.Permissions() {
		for _, scope := range scopes {
			scopeString := fmt.Sprintf("%s:%s", system, scope)
			permissions = append(permissions, scopeString)
		}
	}

	if user.IsAdministrator() {
		permissions = append(permissions, "*:*")
	}

	tokenBuilder := jwt.NewBuilder()
	tokenBuilder.Expiration(time.Now().Add(time.Minute * 15))
	tokenBuilder.NotBefore(time.Now())
	tokenBuilder.Subject(user.GetID())
	tokenBuilder.Audience(TokenAudiences)
	tokenBuilder.Issuer(TokenIssuer)
	tokenBuilder.JwtID(randstr.Base62(256))
	tokenBuilder.Claim("scopes", permissions)

	token, err := tokenBuilder.Build()
	if err != nil {
		c.Abort()
		_ = c.Error(err)
		return
	}

	serializer := jwt.NewSerializer()
	serializer.Sign(jwt.WithKey(resources.PrivateSigningKey.Algorithm(), resources.PrivateSigningKey))
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
	refreshTokenBuilder.Issuer(TokenIssuer)
	refreshTokenBuilder.Audience(RefreshTokenAudiences)
	refreshTokenBuilder.Claim("scopes", permissions)
	refreshTokenBuilder.JwtID(randstr.Base62(128))
	refreshToken, err := refreshTokenBuilder.Build()
	if err != nil {
		c.Abort()
		_ = c.Error(err)
		return
	}

	serializer.Encrypt(jwt.WithKey(jwa.ECDH_ES, resources.PublicEncryptionKey))
	serializedRefreshToken, err := serializer.Serialize(refreshToken)
	if err != nil {
		c.Abort()
		_ = c.Error(err)
		return
	}

	res := types.TokenResponse{
		AccessToken:  string(serializedToken),
		ExpiresIn:    int(math.Ceil(token.Expiration().Sub(time.Now()).Seconds())),
		TokenType:    "Bearer",
		RefreshToken: string(serializedRefreshToken),
	}

	query, err := db.Queries.Raw("register-refresh-token")
	if err != nil {
		fmt.Println(err)
		_ = c.Error(err)
		res.RefreshToken = ""
		goto output
	}

	_, err = db.Pool.Exec(c, query, refreshToken.JwtID(), refreshToken.Expiration())
	if err != nil {
		res.RefreshToken = ""
		fmt.Println(err)
		_ = c.Error(err)
		goto output
	}

output:
	c.JSON(200, res)
}

func checkClientCredentials(c *gin.Context, tokenRequest TokenRequest) interfaces.PermissionableObject {
	clientID := strings.TrimSpace(tokenRequest.ClientID)
	clientSecret := strings.TrimSpace(tokenRequest.ClientSecret)

	if clientID == "" || clientSecret == "" {
		c.Abort()
		apiErrors.ErrMissingParameter.Emit(c)
		return nil
	}

	query, err := db.Queries.Raw("get-client")
	if err != nil {
		c.Abort()
		_ = c.Error(err)
		return nil
	}

	var client types.Client
	err = pgxscan.Get(c, db.Pool, &client, query, clientID)
	if err != nil {
		if err == pgx.ErrNoRows {
			c.Abort()
			apiErrors.ErrInvalidClientCredentials.Emit(c)
			return nil
		}
		c.Abort()
		_ = c.Error(err)
		return nil
	}

	err = client.ReadPermissions(clientID, clientSecret)
	if err != nil {
		if errors.Is(err, types.ErrInvalidSubject) {
			c.Abort()
			apiErrors.ErrInvalidClientCredentials.Emit(c)
		}
		c.Abort()
		_ = c.Error(err)
		return nil
	}
	return client
}

func exchangeAuthorizationCode(c *gin.Context, tokenRequest TokenRequest) interfaces.PermissionableObject {
	// retrieve the verifier from the database
	params, err := db.Redis.Get(c, tokenRequest.State).Bytes()
	if err != nil {
		c.Abort()
		_ = c.Error(err)
		return nil
	}

	tokenParams := types.LoginParameters{}
	json.Unmarshal(params, &tokenParams)

	if strings.TrimSpace(tokenRequest.Code) == "" || strings.TrimSpace(tokenRequest.State) == "" {
		c.Abort()
		apiErrors.ErrMissingParameter.Emit(c)
		return nil
	}

	// now exchange the code for a token
	token, err := oidc.ExternalProvider.Exchange(c, tokenRequest.Code, oauth2.VerifierOption(tokenParams.CodeVerifier), oauth2.SetAuthURLParam("state", tokenRequest.State), oauth2.SetAuthURLParam("redirect_uri", tokenParams.RedirectUri))
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

	var user *types.User
	user, err = utils.GetUser(types.ExternalIdentifier(idToken.Subject))
	if err != nil {
		if err == utils.ErrNoUser {
			newUser, err := utils.CreateUser(idToken.Subject, token)
			if err != nil {
				fmt.Println("error while creating user")
				c.Abort()
				_ = c.Error(err)
				return nil
			}
			return newUser
		}
		c.Abort()
		_ = c.Error(err)
		return nil
	}

	return user
}

// issueFromRefreshToken is the only function that issues tokens directly as
// a user can only gain access to the scopes already present while generating
// the refresh token
func issueFromRefreshToken(c *gin.Context, tokenRequest TokenRequest) interfaces.PermissionableObject {
	decryptedRefreshToken, err := jwe.Decrypt(
		[]byte(tokenRequest.RefreshToken),
		jwe.WithKey(jwa.ECDH_ES, resources.PrivateEncryptionKey),
	)
	if err != nil {
		c.Abort()
		apiErrors.ErrRefreshTokenInvalid.Emit(c)
		return nil
	}

	grantingRefreshToken, err := jwt.Parse(decryptedRefreshToken,
		jwt.WithIssuer("user-management"),
		jwt.WithVerify(true),
		jwt.WithKey(resources.PublicSigningKey.Algorithm(), resources.PublicSigningKey),
	)
	if err != nil {
		c.Abort()
		apiErrors.ErrRefreshTokenInvalid.Emit(c)
		return nil
	}

	query, err := db.Queries.Raw("check-for-refresh-token")
	if err != nil {
		c.Abort()
		_ = c.Error(err)
		return nil
	}

	var tokenAlive bool
	err = pgxscan.Get(c, db.Pool, &tokenAlive, query, grantingRefreshToken.JwtID())
	if err != nil {
		c.Abort()
		_ = c.Error(err)
		return nil
	}

	if !tokenAlive {
		c.Abort()
		apiErrors.ErrRefreshTokenInvalid.Emit(c)
		return nil
	}

	query, err = db.Queries.Raw("revoke-refresh-token")
	if err != nil {
		c.Abort()
		_ = c.Error(err)
	}

	_, err = db.Pool.Exec(c, query, grantingRefreshToken.JwtID())
	if err != nil {
		c.Abort()
		_ = c.Error(err)
	}

	var user *types.User
	user, err = utils.GetUser(types.InternalIdentifier(grantingRefreshToken.Subject()))
	if err != nil {
		if err == utils.ErrNoUser {
			return nil
		}
		c.Abort()
		_ = c.Error(err)
		return nil
	}

	return user

}
