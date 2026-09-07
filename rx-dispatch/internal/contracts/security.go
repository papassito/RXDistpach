package contracts

import "time"

// AuthenticateRequest is sent to RX Security to obtain a session/token.
type AuthenticateRequest struct {
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
	Scope        string `json:"scope"`
}

// AuthenticateResponse contains the granted token and permissions.
type AuthenticateResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expiresAt"`
	Role      string    `json:"role"`
	Scopes    []string  `json:"scopes"`
}

// AuthorizeRequest is checked by services to validate incoming tokens.
type AuthorizeRequest struct {
	Token          string `json:"token"`
	RequiredScope  string `json:"requiredScope"`
	ResourceAction string `json:"resourceAction"`
}

// AuthorizeResponse confirms if the operation is permitted.
type AuthorizeResponse struct {
	Authorized bool   `json:"authorized"`
	ClientID   string `json:"clientId"`
	Role       string `json:"role"`
	Reason     string `json:"reason,omitempty"`
}
