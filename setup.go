package main

import (
	"context"
	"log"

	"github.com/okta/okta-sdk-golang/v4/okta"
	"github.com/sethvargo/go-envconfig"
	"github.com/stytchauth/stytch-go/v12/stytch/b2b/b2bstytchapi"
	"github.com/stytchauth/stytch-go/v12/stytch/b2b/organizations"
	"github.com/stytchauth/stytch-go/v12/stytch/b2b/sso/saml"
)

// OktaSAMLConf contains the info you need to export from Okta's UI to configure SSO
// Those fields reflect the name of the field in the Okta UI.
type OktaSAMLConf struct {
	IdentityProviderIssuer          string `env:"IDENTITY_PROVIDER_ISSUER"`
	IdentityProviderSingleSignOnUrl string `env:"IDENTITY_PROVIDER_SINGLE_SIGN_ON_URL"`
	X509CertificateFilePath         string `env:"X_509_CERTIFICATE_FILE_PATH"`
}

// SetupResult collection UUID of all created resources in the setup proces
type SetupResult struct {
	OrganizationID string
	ConnectionID   string
}

func setup(ctx context.Context, stytchClient *b2bstytchapi.API, oktaClient *okta.APIClient) (result SetupResult) {
	var c OktaSAMLConf
	if err := envconfig.Process(ctx, &c); err != nil {
		log.Fatalf("error loading required env varaibles %s", err)
		return
	}

	// Step 0. Create a New Organisation
	org, err := stytchClient.Organizations.Create(ctx, &organizations.CreateParams{
		OrganizationName: "Example Org Inc.",
		OrganizationSlug: "example-org",
	})
	result.OrganizationID = org.Organization.OrganizationID

	if err != nil {
		log.Fatalf("error creating Organizations %s", err)
		return
	}

	// Step 1. Create a new SAML connection
	sso, err := stytchClient.SSO.SAML.CreateConnection(ctx,
		&saml.CreateConnectionParams{
			DisplayName:    "Okta",
			OrganizationID: result.OrganizationID,
		},
	)
	result.ConnectionID = sso.Connection.ConnectionID

	if err != nil {
		log.Fatalf("error creating SSO SAML Connection %s", err)
		return
	}

	// Step 2: Create and configure a new Okta Application
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

	if err != nil {
		log.Fatalf("error creating Okta Application %s", err)
		return
	}

	_, err = stytchClient.SSO.SAML.UpdateConnection(ctx,
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
