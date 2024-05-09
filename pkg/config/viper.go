package config

import (
	"github.com/spf13/viper"
)

func NewSetupResult(v *viper.Viper) (*SetupResult, error) {
	var C SetupResult
	err := v.Unmarshal(&C)
	return &C, err
}

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
