package config

import (
	"github.com/spf13/viper"
)

// NewSetupConfig is a wrapper that will load both SetupResult and SetupInput
func NewSetupConfig(v *viper.Viper) (*SetupConfig, error) {
	in, err := NewSetupInput(v)
	if err != nil {
		return nil, err
	}

	out, err := NewSetupResult(v)
	if err != nil {
		return nil, err
	}

	return &SetupConfig{
		in, out,
	}, nil
}

// NewSetupResult reads the existing configuration files with viper
// This files would only exist of the setup phase is completed
// An incpmplete setup would also genrate the file with only the value optained
func NewSetupResult(v *viper.Viper) (*SetupResult, error) {
	var C SetupResult
	err := v.Unmarshal(&C)
	return &C, err
}

// NewSetupInput offers the option to customise the inputs for this tutoral
// Those have default values but can be overriden via the config file
func NewSetupInput(v *viper.Viper) (*SetupInput, error) {
	var C SetupInput

	v.SetDefault("okta.SAMLAppLabel", "Example SAML App")
	v.SetDefault("stytch.OrganizationName", "Example SAML App")
	v.SetDefault("stytch.OrganizationSlug", "example-saml-app")
	v.SetDefault("stytch.ConnectionDisplayName", "Okta")

	err := v.Unmarshal(&C)
	return &C, err
}

// ViperConfigProvider implement the setup.ConfigProvider interface
// This allows us to use viper to lead and persiste setup data
type ViperConfigProvider struct {
	data *SetupConfig
}

func (c *ViperConfigProvider) Load() (*SetupConfig, error) {
	v := viper.GetViper()
	data, err := NewSetupConfig(v)
	c.data = data
	return data, err
}

func (c *ViperConfigProvider) Save() error {
	v := viper.GetViper()

	// This remove secrets from the config written
	v.Set("stytch.project_id", "<redacted>")
	v.Set("stytch.secret", "<redacted>")
	v.Set("stytch.project_public_id", "<redacted>")
	v.Set("okta.org_url", "<redacted>")
	v.Set("okta.api_token", "<redacted>")

	// Update the config
	v.Set("stytch.organization_id", c.data.StytchResult.OrganizationID)
	v.Set("stytch.connection_id", c.data.StytchResult.ConnectionID)
	v.Set("okta.application_id", c.data.OktaResult.ApplicationID)
	if c.data.StytchResult.SsoParameters != nil {
		v.Set("stytch.sso_parameters.acsurl", c.data.StytchResult.SsoParameters.AcsUrl)
		v.Set("stytch.sso_parameters.audience", c.data.StytchResult.SsoParameters.Audience)
	}

	v.SafeWriteConfig()
	return v.WriteConfig()
}

// NewClientConfig uses Viper to read secret from environement varaibles
// those secrets will be used to configure Stytch and Okta clients
func NewClientConfig(v *viper.Viper) (*ClientsConf, error) {
	var C ClientsConf

	v.BindEnv("stytch.project_id")
	v.BindEnv("stytch.secret")
	v.BindEnv("stytch.project_public_id")
	v.BindEnv("okta.org_url")
	v.BindEnv("okta.api_token")

	err := v.Unmarshal(&C)
	return &C, err
}
