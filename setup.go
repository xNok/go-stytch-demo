package main

import (
	"context"
	"log"
	"os"

	"github.com/okta/okta-sdk-golang/okta"
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

func setup(ctx context.Context, stytchClient *b2bstytchapi.API, oktaClient *okta.Client) (result SetupResult) {
	var c OktaSAMLConf
	if err := envconfig.Process(ctx, &c); err != nil {
		log.Fatalf("error loading required env varaibles %s", err)
		return
	}

	org, err := stytchClient.Organizations.Create(ctx, &organizations.CreateParams{
		OrganizationName: "Example Org Inc.",
		OrganizationSlug: "example-org",
	})
	result.OrganizationID = org.Organization.OrganizationID

	if err != nil {
		log.Fatalf("error creating Organizations %s", err)
		return
	}

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

	x509, err := os.ReadFile(c.X509CertificateFilePath)
	if err != nil {
		log.Fatalf("error reading SSO X509 Certificate %s", err)
		return
	}

	_, err = stytchClient.SSO.SAML.UpdateConnection(ctx,
		&saml.UpdateConnectionParams{
			ConnectionID:    result.ConnectionID,
			OrganizationID:  result.OrganizationID,
			IdpEntityID:     c.IdentityProviderIssuer,
			IdpSSOURL:       c.IdentityProviderSingleSignOnUrl,
			X509Certificate: string(x509[:]),
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
