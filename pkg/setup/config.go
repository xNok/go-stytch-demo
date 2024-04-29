package setup

type Conf struct {
	StytchConf StytchConf
	OktaConf   OktaConf
}

type StytchConf struct {
	ProjectID   string `env:"STYTCH_PROJECT_ID"`
	Secret      string `env:"STYTCH_SECRET"`
	PublicToken string `env:"STYTCH_PROJECT_PUBLIC_ID"`
}

type OktaConf struct {
	OrgUrl   string `env:"OKTA_ORG_URL"`
	APIToken string `env:"OKTA_API_TOKEN"`
}
