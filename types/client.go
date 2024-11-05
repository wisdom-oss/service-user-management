package types

// Client represents an external app which has been granted permanent access
// to the WISdoM platform.
// As the app is always external, it can't be an administrator.
// Therefore the IsAdministrator function always returns false
type Client struct {
	ClientID        string `json:"clientID"`
	ClientSecret    string `json:"clientSecret"`
	ApplicationName string `json:"applicationName"`
}

func (c Client) GetID() string {
	return c.ClientID
}

func (c Client) Permissions() map[string][]string {
	// TODO: Fill with code which fetches the permissions
	return map[string][]string{}
}

func (c Client) IsAdministrator() bool {
	return false
}
