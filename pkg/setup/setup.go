package setup

import (
	"context"
	"log"

	"github.com/okta/okta-sdk-golang/v4/okta"
	"github.com/stytchauth/stytch-go/v12/stytch/b2b/b2bstytchapi"
	"github.com/xNok/go-stytch-demo/pkg/config"
)

// Set up an Okta SAML Connection
// ref: https://stytch.com/docs/b2b/guides/sso/okta-saml
type OktaSAMLConnectionBootstraper struct {
	// Clients
	StytchClient *b2bstytchapi.API
	OktaClient   *okta.APIClient
	// Persistent config (Those will be needed in the )
	ConfProvider SetupConfig
}

func NewOktaSAMLConnectionBootstraper(stytch *b2bstytchapi.API, okta *okta.APIClient) *OktaSAMLConnectionBootstraper {
	return &OktaSAMLConnectionBootstraper{
		StytchClient: stytch,
		OktaClient:   okta,
		ConfProvider: &config.ViperConfigProvider{},
	}
}

// SetupConfig is a abstraction to help us retrive our configuration data
// For testing purposed thay can be stored in YAML file
// But in a live application we might rely on a config server or a vault
type SetupConfig interface {
	Save() error
	Load() (*config.SetupConfig, error)
}

// Setup will Perform the bootstraping oprations between Stych and Okta
// To ensure idempotency of this function, after each step is performed the resulting configuration is persisted
// This means that if run a second time only the steps not yet completed will be played again
func (s *OktaSAMLConnectionBootstraper) Setup(ctx context.Context) (err error) {
	// Load our configuration file, this file is empty if we start from scrath
	conf, err := s.ConfProvider.Load()
	if err != nil {
		log.Fatalf("error loading Configuration %s", err)
		return err
	}

	// Step 0. Create a New Organisation
	if conf.StytchResult.OrganizationID == "" {
		conf.StytchResult.OrganizationID, err = s.setupStytchOrganisation(ctx, &conf.StytchSetupInput)

		if err != nil {
			log.Fatalf("error creating Organizations %s", err)
			return err
		}

		s.ConfProvider.Save()
	}

	// Step 1. Create a new SAML connection
	if conf.StytchResult.ConnectionID == "" {
		conf.StytchResult.ConnectionID, conf.StytchResult.SsoParameters, err = s.createStytchConnection(ctx,
			&conf.StytchSetupInput, conf.StytchResult.OrganizationID)

		if err != nil {
			log.Fatalf("error creating SSO SAML Connection %s", err)
			return
		}

		s.ConfProvider.Save()
	}

	// Step 2: Create and configure a new Okta Application
	if conf.OktaResult.ApplicationID == "" {
		conf.OktaResult.ApplicationID, err = s.setupOktaSamlApplication(ctx, &conf.OktaSetupInput, conf.StytchResult.SsoParameters)

		if err != nil {
			log.Fatalf("error creating Okta Application %s", err)
			return
		}

		s.ConfProvider.Save()
	}

	// Step 3: Fetch Okta SAML Metdata
	conf.OktaResult.SsoParameters, err = s.getOktaSamlApplicationMetada(ctx, conf.ApplicationID)
	if err != nil {
		log.Fatalf("error fetch Okta Application SSO metadata %s", err)
		return
	}

	// Step 4: Update Stych SSO Connactions
	err = s.updateStytchConnection(ctx, conf.OrganizationID, conf.ConnectionID, conf.OktaResult.SsoParameters)
	if err != nil {
		log.Fatalf("error updating SSO SAML Connection %s", err)
		return
	}

	return
}
