package config

// SetupResult collection UUID of all created resources in the setup proces
type SetupResult struct {
	StytchResult `mapstructure:"stytch"`
	OktaResult   `mapstructure:"okta"`
}

type StytchResult struct {
	OrganizationID string              `mapstructure:"organization_id"`
	ConnectionID   string              `mapstructure:"connection_id"`
	SsoParameters  *StychSsoParameters `mapstructure:"sso_parameters"`
}

type OktaResult struct {
	ApplicationID string             `mapstructure:"application_id"`
	SsoParameters *OktaSsoParameters `mapstructure:"sso_parameters"`
}

// SsoStychParameters represent the metadata needed to configure okta SSO obtained from Stych
type StychSsoParameters struct {
	AcsUrl   string
	Audience string
}

// SsoOktaParameters  represent the metadata needed to configure okta SSO obtained from Okta
type OktaSsoParameters struct {
	IdpEntityID     string
	IdpSSOURL       string
	X509Certificate string
}
