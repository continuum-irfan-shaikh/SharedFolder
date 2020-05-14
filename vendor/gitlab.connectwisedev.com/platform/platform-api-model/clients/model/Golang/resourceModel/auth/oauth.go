package auth

import "time"

// OAuthCodeReq represents request for authorization code grant type
type OAuthCodeReq struct {
	Code         string `json:"code"`
	ClientID     string `json:"client_id"      validate:"required"`
	RedirectURI  string `json:"redirect_uri"   validate:"uri"`
	Scope        string `json:"scope"`
	ResponseType string `json:"response_type"`
}

// OAuthRequest represents authorization request
type OAuthRequest struct {
	GrantType    string `json:"grant_type"    validate:"required,eq=authorization_code|eq=client_credentials|eq=refresh_token"`
	ClientID     string `json:"client_id"     validate:"required"`
	ClientSecret string `json:"client_secret" validate:"required"`
	Scope        string `json:"scope"`
	RefreshToken string `json:"refresh_token"`
}

// OAuthToken represents authorization token
type OAuthToken struct {
	ClientID     string `json:"client_id,omitempty"`
	PartnerID    string `json:"memberid,omitempty"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

// OAuthClientBaseConfig represents basic OAuth2 client configuration
type OAuthClientBaseConfig struct {
	DefaultScopes   []string `json:"default_scopes,omitempty"`
	Scopes          []string `json:"scopes,omitempty"`
	GrantTypes      []string `json:"grant_types,omitempty"`
	RedirectionUris []string `json:"redirection_uris,omitempty"`
	ClientName      string   `json:"client_name,omitempty"`
	Status          string   `json:"status,omitempty"              validate:"oneof=Active Inactive"`
}

// AdvancedOAuth2ClientConfig represents advanced OAuth2 client configuration
type AdvancedOAuth2ClientConfig struct {
	OAuthClientBaseConfig
	ClientID     string     `json:"client_id,omitempty"`
	ClientSecret string     `json:"client_secret,omitempty"`
	UpdatedAt    *time.Time `json:"updated_at,omitempty"`
	CreatedBy    string     `json:"created_by,omitempty"`
}

// OAuthIdentity oauth identity
type OAuthIdentity struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Scope        string `json:"scope"`
	GrantType    string `json:"grant_type"`
	AuthCode     string `json:"code"`
	RedirectURI  string `json:"redirect_uri"`
	RefreshToken string `json:"refresh_token"`
}
