package setup

import (
	"context"
	"log"

	"github.com/okta/okta-sdk-golang/v4/okta"
	"github.com/stytchauth/stytch-go/v12/stytch/b2b/b2bstytchapi"
	"github.com/stytchauth/stytch-go/v12/stytch/b2b/organizations"
	"github.com/stytchauth/stytch-go/v12/stytch/b2b/sso/saml"
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
		ConfProvider: &YAMLEntry{
			Path: "./setup.yaml",
		},
	}
}

// SetupResult collection UUID of all created resources in the setup proces
type SetupResult struct {
	// Stych
	OrganizationID string
	ConnectionID   string
	// Okta
	ApplicationID string

	SsoStychParameters *SsoStychParameters
	SsoOktaParameters  *SsoOktaParameters
}

// SsoStychParameters represent the metadata needed to configure okta SSO obtained from Stych
type SsoStychParameters struct {
	SsoAcsUrl   string
	SsoAudience string
}

// SsoOktaParameters  represent the metadata needed to configure okta SSO obtained from Okta
type SsoOktaParameters struct {
	IdpEntityID     string
	IdpSSOURL       string
	X509Certificate string
}

// SetupConfig is a abstraction to help us retrive our configuration data
// For testing purposed thay can be stored in YAML file
// But in a live application we might rely on a config server or a vault
type SetupConfig interface {
	Save() error
	Load() (*SetupResult, error)
	Get() *SetupResult
}

// Setup will Perform the bootstraping oprations between Stych and Okta
// To ensure idempotency of this function, after each step is performed the resulting configuration is persisted
// This means that if run a second time only the steps not yet completed will be played again
func (s *OktaSAMLConnectionBootstraper) Setup(ctx context.Context) (conf *SetupResult) {
	// Load our configuration file, this file is empty if we start from scrath
	conf, err := s.ConfProvider.Load()
	if err != nil {
		log.Fatalf("error loading Configuration %s", err)
		return
	}

	// Step 0. Create a New Organisation
	if conf.OrganizationID == "" {
		conf.OrganizationID, err = s.setupStytchOrganisation(ctx)

		if err != nil {
			log.Fatalf("error creating Organizations %s", err)
			return
		}

		s.ConfProvider.Save()
	}

	// Step 1. Create a new SAML connection
	if conf.ConnectionID == "" {
		conf.ConnectionID, conf.SsoStychParameters, err = s.createStytchConnection(ctx, conf.OrganizationID)

		if err != nil {
			log.Fatalf("error creating SSO SAML Connection %s", err)
			return
		}

		s.ConfProvider.Save()
	}

	// Step 2: Create and configure a new Okta Application
	if conf.ApplicationID == "" {

		conf.SsoOktaParameters, err = s.setupOktaSamlApplication(ctx, conf.SsoStychParameters)

		if err != nil {
			log.Fatalf("error creating Okta Application %s", err)
			return
		}

		s.ConfProvider.Save()
	}

	err = s.updateStytchConnection(ctx, conf.OrganizationID, conf.ConnectionID, conf.SsoOktaParameters)

	if err != nil {
		log.Fatalf("error updating SSO SAML Connection %s", err)
		return
	}

	return
}

func (s *OktaSAMLConnectionBootstraper) setupStytchOrganisation(ctx context.Context) (string, error) {
	org, err := s.StytchClient.Organizations.Create(ctx, &organizations.CreateParams{
		OrganizationName: "Example Org Inc.",
		OrganizationSlug: "example-org",
	})

	if err != nil {
		return "", err
	}

	return org.Organization.OrganizationID, nil
}

func (s *OktaSAMLConnectionBootstraper) createStytchConnection(ctx context.Context, organizationID string) (string, *SsoStychParameters, error) {
	sso, err := s.StytchClient.SSO.SAML.CreateConnection(ctx,
		&saml.CreateConnectionParams{
			DisplayName:    "Okta",
			OrganizationID: organizationID,
		},
	)

	if err != nil {
		return "", nil, err
	}

	return sso.Connection.ConnectionID, &SsoStychParameters{
		SsoAcsUrl:   sso.Connection.AcsURL,
		SsoAudience: sso.Connection.AudienceURI,
	}, nil
}

func (s *OktaSAMLConnectionBootstraper) updateStytchConnection(ctx context.Context, organizationID, connectionID string, conf *SsoOktaParameters) error {
	_, err := s.StytchClient.SSO.SAML.UpdateConnection(ctx,
		&saml.UpdateConnectionParams{
			ConnectionID:    connectionID,
			OrganizationID:  organizationID,
			IdpEntityID:     conf.IdpEntityID,
			IdpSSOURL:       conf.IdpSSOURL,
			X509Certificate: conf.X509Certificate,
			AttributeMapping: map[string]any{
				"email":      "NameID",
				"first_name": "firstName",
				"last_name":  "lastName",
			},
		})

	if err != nil {
		return err
	}

	return nil
}

func (s *OktaSAMLConnectionBootstraper) setupOktaSamlApplication(ctx context.Context, conf *SsoStychParameters) (*SsoOktaParameters, error) {
	samlApp := okta.NewSamlApplication()
	samlApp.Label = okta.PtrString("Example SAML App")
	samlApp.SignOnMode = okta.PtrString("SAML_2_0")
	samlApp.Visibility = okta.NewApplicationVisibility()
	samlApp.Settings = okta.NewSamlApplicationSettingsWithDefaults()
	samlApp.Settings.SignOn = &okta.SamlApplicationSettingsSignOn{
		DefaultRelayState:     okta.PtrString(""),
		SsoAcsUrl:             &conf.SsoAcsUrl,
		IdpIssuer:             okta.PtrString("http://www.okta.com/${org.externalKey}"),
		Audience:              &conf.SsoAudience,
		Recipient:             &conf.SsoAcsUrl,
		Destination:           &conf.SsoAcsUrl,
		SubjectNameIdTemplate: okta.PtrString("${user.userName}"),
		SubjectNameIdFormat:   okta.PtrString("urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress"),
		ResponseSigned:        okta.PtrBool(true),
		AssertionSigned:       okta.PtrBool(true),
		SignatureAlgorithm:    okta.PtrString("RSA_SHA256"),
		DigestAlgorithm:       okta.PtrString("SHA256"),
		HonorForceAuthn:       okta.PtrBool(true),
		AuthnContextClassRef:  okta.PtrString("urn:oasis:names:tc:SAML:2.0:ac:classes:PasswordProtectedTransport"),
		SpIssuer:              nil,
		AttributeStatements: []okta.SamlAttributeStatement{
			{
				Type:      okta.PtrString("EXPRESSION"),
				Name:      okta.PtrString("firstName"),
				Namespace: okta.PtrString("urn:oasis:names:tc:SAML:2.0:attrname-format:basic"),
				Values: []string{
					"user.firstName",
				},
			}, {
				Type:      okta.PtrString("EXPRESSION"),
				Name:      okta.PtrString("lastName"),
				Namespace: okta.PtrString("urn:oasis:names:tc:SAML:2.0:attrname-format:basic"),
				Values: []string{
					"user.lastName",
				},
			},
		},
	}

	_, _, err := s.OktaClient.ApplicationAPI.CreateApplication(ctx).Application(
		okta.ListApplications200ResponseInner{
			SamlApplication: samlApp,
		},
	).Execute()

	if err != nil {
		return nil, err
	}

	return nil, nil
}
