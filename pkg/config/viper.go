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
type ViperConfigProvider struct{}

func (c *ViperConfigProvider) Load() (*SetupConfig, error) {
	v := viper.GetViper()
	return NewSetupConfig(v)
}

func (c *ViperConfigProvider) Save() error {
	v := viper.GetViper()
	return v.SafeWriteConfig()
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
