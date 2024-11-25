package utils

import (
	"context"
	"errors"
	"net/http"
	"reflect"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/go-jose/go-jose/v4/json"
	"golang.org/x/oauth2"

	"microservice/internal/db"
	oidc2 "microservice/oidc"
	"microservice/types"
)

var ErrNoUser = errors.New("no user with this id")

// GetUser retrieves a User object from the database
func GetUser[T types.InternalIdentifier | types.ExternalIdentifier](id T) (*types.User, error) {
	externalIDType := reflect.TypeOf(types.ExternalIdentifier(""))
	internalIDType := reflect.TypeOf(types.InternalIdentifier(""))
	parameterType := reflect.TypeOf(id)

	var rawQuery string
	var err error
	if parameterType.AssignableTo(externalIDType) {
		rawQuery, err = db.Queries.Raw("get-user-by-external-id")
		if err != nil {
			return nil, err
		}
	}
	if parameterType.AssignableTo(internalIDType) {
		rawQuery, err = db.Queries.Raw("get-user-by-internal-id")
		if err != nil {
			return nil, err
		}

	}

	var user types.User
	err = pgxscan.Get(context.Background(), db.Pool, &user, rawQuery, string(id))
	if err != nil {
		if pgxscan.NotFound(err) {
			return nil, ErrNoUser
		}
		return nil, err
	}
	return &user, nil
}

func CreateUser(subject string, token *oauth2.Token) (*types.User, error) {
	req, err := http.NewRequest("GET", oidc2.OidcProvider.UserInfoEndpoint(), nil)
	if err != nil {
		return nil, err
	}
	token.SetAuthHeader(req)
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	var desiredUserInformation struct {
		Username string `json:"preferred_username"`
		Name     string `json:"name"`
		Email    string `json:"email"`
	}

	err = json.NewDecoder(res.Body).Decode(&desiredUserInformation)
	if err != nil {
		return nil, err
	}

	query, err := db.Queries.Raw("create-user")
	if err != nil {
		return nil, err
	}

	_, err = db.Pool.Exec(context.Background(), query, subject, desiredUserInformation.Name, desiredUserInformation.Username, desiredUserInformation.Email)
	if err != nil {
		return nil, err
	}

	return GetUser(types.ExternalIdentifier(subject))
}
