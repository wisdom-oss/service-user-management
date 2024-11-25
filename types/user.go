package types

import (
	"context"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/go-jose/go-jose/v4/json"

	"microservice/internal/db"
)

type User struct {
	ID                 string `json:"id" db:"id"`
	ExternalIdentifier string `json:"externalIdentifier" db:"external_identifier"`
	Name               string `json:"name" db:"name"`
	Email              string `json:"email" db:"email"`
	Username           string `json:"username" db:"username"`
	Disabled           bool   `json:"disabled" db:"disabled"`
	Administrator      bool   `json:"administrator" db:"is_admin"`
}

func (u User) GetID() string {
	return u.ID
}

func (u User) Permissions() map[string][]string {
	if u.Administrator {
		query, err := db.Queries.Raw("get-services")
		if err != nil {
			panic(err)
		}
		var services []Service
		err = pgxscan.Select(context.Background(), db.Pool, &services, query)
		if err != nil {
			panic(err)
		}

		permissions := make(map[string][]string)
		for _, service := range services {
			permissions[service.Name] = append(permissions[service.Name], service.SupportedScopes...)
		}
		return permissions
	}
	query, err := db.Queries.Raw("get-user-permissions")
	if err != nil {
		return nil
	}

	var permissionMappings []struct {
		Name  string `db:"name"`
		Level string `db:"level"`
	}
	err = pgxscan.Select(context.Background(), db.Pool, &permissionMappings, query, u.GetID())
	if err != nil {
		panic(err)
	}

	permissions := make(map[string][]string)

	for _, mapping := range permissionMappings {
		permissions[mapping.Name] = append(permissions[mapping.Name], mapping.Level)
	}

	return permissions
}

func (u User) IsAdministrator() bool {
	return u.Administrator
}

func (u User) IsActive() bool {
	return !u.Disabled
}

func (u User) MarshalJSON() ([]byte, error) {
	type output struct {
		ID                 string              `json:"id" db:"id"`
		ExternalIdentifier string              `json:"externalIdentifier" db:"external_identifier"`
		Name               string              `json:"name" db:"name"`
		Email              string              `json:"email" db:"email"`
		Username           string              `json:"username" db:"username"`
		Disabled           bool                `json:"disabled" db:"disabled"`
		Administrator      bool                `json:"administrator" db:"is_admin"`
		Permissions        map[string][]string `json:"permissions"`
	}
	o := output{
		ID:                 u.ID,
		ExternalIdentifier: u.ExternalIdentifier,
		Name:               u.Name,
		Email:              u.Email,
		Username:           u.Username,
		Disabled:           u.Disabled,
		Administrator:      u.Administrator,
		Permissions:        u.Permissions(),
	}
	return json.Marshal(o)
}
