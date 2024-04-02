package main

import (
	"context"
	"log"

	"github.com/okta/okta-sdk-golang/v4/okta"
	"github.com/sethvargo/go-envconfig"
	"github.com/stytchauth/stytch-go/v12/stytch/b2b/b2bstytchapi"
)

type Conf struct {
	StytchConf StytchConf
	OktaConf   OktaConf
}

type StytchConf struct {
	ProjectID string `env:"STYTCH_PROJECT_ID"`
	Secret    string `env:"STYTCH_SECRET"`
}

type OktaConf struct {
	OrgUrl   string `env:"OKTA_ORG_URL"`
	APIToken string `env:"OKTA_API_TOKEN"`
}

func main() {
	ctx := context.Background()

	var c Conf
	if err := envconfig.Process(ctx, &c); err != nil {
		panic(err)
	}

	// Step 1: Instanciate stytch client
	stytchClient, err := b2bstytchapi.NewClient(
		c.StytchConf.ProjectID,
		c.StytchConf.Secret,
	)

	if err != nil {
		log.Fatalf("error instantiating API client %s", err)
	}

	// Step 2: Instanciate Okta client
	oktaConfig, err := okta.NewConfiguration(
		okta.WithOrgUrl(c.OktaConf.OrgUrl),
		okta.WithToken(c.OktaConf.APIToken),
	)

	if err != nil {
		log.Fatalf("error instantiating Okta API client %s", err)
	}

	oktaClient := okta.NewAPIClient(oktaConfig)

	setup(ctx, stytchClient, oktaClient)
}
