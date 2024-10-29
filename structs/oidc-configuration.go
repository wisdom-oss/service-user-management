package structs

type OIDCConfiguration struct {
	ClientID       string `json:"clientID"`
	IssuerURL      string `json:"issuerURL"`
	RedirectionURL string `json:"redirectionURL"`
}
