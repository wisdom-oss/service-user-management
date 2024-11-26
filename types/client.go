package types

import (
	"errors"
	"fmt"
	"microservice/resources"
	"strings"

	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwe"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

var ErrNoScopesSet = errors.New("no scope set on client token")
var ErrScopesWrongFormat = errors.New("scopes not in a supported format")
var ErrInvalidSubject = errors.New("invaild subject for supplied")

// Client represents an external app which has been granted permanent access
// to the WISdoM platform.
// As the app is always external, it can't be an administrator.
// Therefore the IsAdministrator function always returns false
type Client struct {
	ID      string `json:"clientID" db:"id"`
	Name    string `json:"name" db:"name"`
	Contact struct {
		Name  string `json:"name" db:"contact_name"`
		EMail string `json:"email" db:"contact_email"`
	} `json:"contact"`
	permissions map[string][]string `json:"-" db:"-"`
}

func (c Client) GetID() string {
	return c.ID
}

func (c Client) Permissions() map[string][]string {
	return c.permissions
}

func (c Client) IsActive() bool {
	return true
}

func (c Client) IsAdministrator() bool {
	return false
}

func (c *Client) ReadPermissions(clientID, clientSecret string) error {
	decryptedClientSecret, err := jwe.Decrypt(
		[]byte(clientSecret),
		jwe.WithKey(jwa.ECDH_ES, resources.PrivateEncryptionKey),
	)
	if err != nil {
		return err
	}

	// TODO: Continue development here
	clientToken, err := jwt.Parse(decryptedClientSecret,
		jwt.WithIssuer("user-management"),
		jwt.WithVerify(true),
		jwt.WithKey(resources.PublicSigningKey.Algorithm(), resources.PublicSigningKey),
	)

	if err != nil {
		return fmt.Errorf("unable to read client permissions: %w", err)
	}

	if clientToken.Subject() != clientID {
		return ErrInvalidSubject
	}

	iface, set := clientToken.PrivateClaims()["scopes"]
	if !set {
		return ErrNoScopesSet
	}

	scopes, ok := iface.([]string)
	if !ok {
		return ErrScopesWrongFormat
	}

	permissions := make(map[string][]string)
	for _, scope := range scopes {
		parts := strings.Split(scope, ":")
		service := parts[0]
		level := parts[1]

		if permissions[service] == nil {
			permissions[service] = make([]string, 0)
		}

		permissions[service] = append(permissions[service], level)
	}

	c.permissions = permissions
	return nil
}
