package types

type LoginParameters struct {
	RedirectUri  string `json:"redirect_uri"`
	CodeVerifier string `json:"codeVerifier"`
}
