package main

type CloudModel struct {
	Name    string
	AuthURL string `json:"-" env:"OS_AUTH_URL`

	Username string `json:"username,omitempty"`
	UserID   string `json:"-"`

	Password string `json:"password,omitempty"`
	Passcode string `json:"passcode,omitempty"`

	DomainID   string `json:"-"`
	DomainName string `json:"name,omitempty"`

	TokenID string `json:"-"`

	ProjectID   string
	ProjectName string

	System  bool
	TrustID string
}

type Clouds struct {
	name string
}
