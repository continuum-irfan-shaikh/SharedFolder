package am

// OAuthTokenInfo am token info
type OAuthTokenInfo struct {
	AccessToken string   `json:"access_token"`
	GrantType   string   `json:"grant_type"`
	Realm       string   `json:"realm"`
	TokenType   string   `json:"token_type"`
	ClientID    string   `json:"client_id"`
	PartnerID   string   `json:"memberid,omitempty"`
	UserID      string   `json:"uid,omitempty"`
	ExpiresIn   int      `json:"expires_in"`
	Scope       []string `json:"scope"`
	Mail        string   `json:"mail,omitempty"`
}
