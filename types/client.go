package types

import (
	"microservice/resources"

	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwe"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

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
}

func (c Client) GetID() string {
	return c.ID
}

func (c Client) Permissions() map[string][]string {
	// TODO: Fill with code which fetches the permissions
	return map[string][]string{}
}

func (c Client) IsAdministrator() bool {
	return false
}

func (c Client) ReadPermissions(secret string) error {
	decryptedClientSecret, err := jwe.Decrypt(
		[]byte(secret),
		jwe.WithKey(jwa.ECDH_ES, resources.PrivateEncryptionKey),
	)
	if err != nil {
		return err
	}

	// TODO: Continue development here
	_, err = jwt.Parse(decryptedClientSecret,
		jwt.WithIssuer("user-management"),
		jwt.WithVerify(true),
		jwt.WithKey(resources.PublicSigningKey.Algorithm(), resources.PublicSigningKey),
	)

	return nil
}
