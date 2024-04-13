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

// SetupResult collection UUID of all created resources in the setup proces
type SetupResult struct {
	OrganizationID string
	ConnectionID   string
}

// SsoConnectionParameter represent the metadata needed to configure okta SSO obtained from Stych
type SsoConnectionParameter struct {
	SsoAcsUrl   string
	SsoAudience string
}

type SetupConfig interface {
	Save() error
	Load() (*SetupResult, error)
}

func (s *OktaSAMLConnectionBootstraper) Setup(ctx context.Context) (result SetupResult) {
	// Load our configuration file, this file is empty if we start from scrath
	conf, err := s.ConfProvider.Load()
	if err != nil {
		log.Fatalf("error loading Configuration %s", err)
		return
	}

	// Step 0. Create a New Organisation
	if conf.OrganizationID != "" {
		result.OrganizationID, err = s.setupStytchOrganisation(ctx)

		if err != nil {
			log.Fatalf("error creating Organizations %s", err)
			return
		}

		s.ConfProvider.Save()
	}

	// Step 1. Create a new SAML connection
	if conf.ConnectionID != "" {
		result.ConnectionID, err = s.setupStytchConnection(ctx, result.OrganizationID)

		if err != nil {
			log.Fatalf("error creating SSO SAML Connection %s", err)
			return
		}

		s.ConfProvider.Save()
	}

	// Step 2: Create and configure a new Okta Application


	if err != nil {
		log.Fatalf("error creating Okta Application %s", err)
		return
	}

	_, err = s.StytchClient.SSO.SAML.UpdateConnection(ctx,
		&saml.UpdateConnectionParams{
			ConnectionID:    result.ConnectionID,
			OrganizationID:  result.OrganizationID,
			IdpEntityID:     *oktaApp.SamlApplication.Settings.SignOn.IdpIssuer,
			IdpSSOURL:       *oktaApp.SamlApplication.Accessibility.ErrorRedirectUrl,
			X509Certificate: oktaApp.SamlApplication.Settings.SignOn.SpCertificate.GetX5c()[0],
			AttributeMapping: map[string]any{
				"email":      "NameID",
				"first_name": "firstName",
				"last_name":  "lastName",
			},
		})

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

func (s *OktaSAMLConnectionBootstraper) setupStytchConnection(ctx context.Context, organizationID string) (string, error) {
	sso, err := s.StytchClient.SSO.SAML.CreateConnection(ctx,
		&saml.CreateConnectionParams{
			DisplayName:    "Okta",
			OrganizationID: organizationID,
		},
	)

	if err != nil {
		return "", err
	}

	// SsoAcsUrl = sso.Connection.AcsURL,
	// SsoAudience = sso.Connection.AudienceURI,

	return sso.Connection.ConnectionID, nil
}

func (s *OktaSAMLConnectionBootstraper) setupOktaSamlApplication(ctx context.Context, conf SsoConnectionParameter) {
	oktaAppName := "okta_stych"
	oktaApp, _, err := oktaClient.ApplicationAPI.CreateApplication(ctx).
		Application(okta.SamlApplicationAsListApplications200ResponseInner(
			&okta.SamlApplication{
				Name: &oktaAppName,
				Settings: &okta.SamlApplicationSettings{
					SignOn: &okta.SamlApplicationSettingsSignOn{
						SsoAcsUrl:           &sso.Connection.AcsURL,
						Audience:            &sso.Connection.AudienceURI,
						SubjectNameIdFormat: okta.PtrString("urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress"),
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
					},
				},
			},
		)).Execute()
}
