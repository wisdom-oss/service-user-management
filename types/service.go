package types

import (
	"context"
	"errors"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"

	"microservice/internal/db"
)

type Service struct {
	ID              string   `json:"id" db:"id"`
	Name            string   `json:"name" db:"name"`
	Description     *string  `json:"description" db:"description"`
	SupportedScopes []string `json:"supportedScopes" db:"supported_scope_levels"`
}

var ErrUnknownService = errors.New("unknown service")

func (s Service) LoadFromDB(identifier any) (err error) {
	var query string
	switch identifier.(type) {
	case ExternalIdentifier:
		query, err = db.Queries.Raw("get-service-by-external-id")
		break
	case InternalIdentifier:
		query, err = db.Queries.Raw("get-service-by-internal-id")
		break
	default:
		err = errors.New("invalid identifier type")
	}
	if err != nil {
		return err
	}

	err = pgxscan.Get(context.Background(), db.Pool, &s, query, identifier.(string))
	if err != nil {
		if err == pgx.ErrNoRows {
			return errors.Join(ErrUnknownService, err)
		}
		return err
	}
	return nil
}
