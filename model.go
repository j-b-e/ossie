package main

// type Cloud struct {
// 	Name    string
// 	AuthURL string `json:"-" env:"OS_AUTH_URL`

// 	Username string `json:"username,omitempty"`
// 	UserID   string `json:"-"`

// 	Password string `json:"password,omitempty"`
// 	Passcode string `json:"passcode,omitempty"`

// 	DomainID   string `json:"-"`
// 	DomainName string `json:"name,omitempty"`

// 	RegionName string
// 	RegionID   string

// 	TokenID string `json:"-"`

// 	ProjectID   string
// 	ProjectName string
// 	TenantName  string
// 	TenantID    string

// 	System  bool
// 	TrustID string

// 	CaCert string

// 	Interface string

//		IdentityApiVersion string
//	}

type Cloud struct {
	Name string
	Env  map[string]string
}

func GetCloud(name string, clouds []Cloud) Cloud {
	for _, v := range clouds {
		if v.Name == name {
			return v
		}
	}
	return Cloud{}
}
