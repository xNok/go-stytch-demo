package config

type ClientsConf struct {
	StytchConf *StytchConf `mapstructure:"STYTCH"`
	OktaConf   *OktaConf   `mapstructure:"OKTA"`
}

type StytchConf struct {
	ProjectID   string `mapstructure:"PROJECT_ID"`
	Secret      string `mapstructure:"SECRET"`
	PublicToken string `mapstructure:"PROJECT_PUBLIC_ID"`
}

type OktaConf struct {
	OrgUrl   string `mapstructure:"ORG_URL"`
	APIToken string `mapstructure:"API_TOKEN"`
}
