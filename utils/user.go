package utils

import (
	"context"
	"errors"
	"net/http"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/go-jose/go-jose/v4/json"
	"golang.org/x/oauth2"

	"microservice/internal/db"
	oidc2 "microservice/oidc"
	"microservice/structs"
)

var ErrNoUser = errors.New("no user with this id")

// GetUser retrieves a User object from the database
func GetUser(externalID string) (*structs.User, error) {
	rawQuery, err := db.Queries.Raw("get-user-by-external-id")
	if err != nil {
		return nil, err
	}

	var user structs.User
	err = pgxscan.Get(context.Background(), db.Pool, &user, rawQuery, externalID)
	if err != nil {
		if pgxscan.NotFound(err) {
			return nil, ErrNoUser
		}
		return nil, err
	}
	return &user, nil
}

func CreateUser(token *oauth2.Token) (*structs.User, error) {
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

	return nil, nil
}
