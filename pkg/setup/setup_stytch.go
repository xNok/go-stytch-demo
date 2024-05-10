package setup

import (
	"context"

	"github.com/stytchauth/stytch-go/v12/stytch/b2b/organizations"
	"github.com/stytchauth/stytch-go/v12/stytch/b2b/sso/saml"
	"github.com/xNok/go-stytch-demo/pkg/config"
)

func (s *OktaSAMLConnectionBootstraper) setupStytchOrganisation(ctx context.Context, stytchConf *config.StytchSetupInput) (string, error) {
	org, err := s.StytchClient.Organizations.Create(ctx, &organizations.CreateParams{
		OrganizationName: stytchConf.OrganizationName,
		OrganizationSlug: stytchConf.OrganizationSlug,
	})

	if err != nil {
		return "", err
	}

	return org.Organization.OrganizationID, nil
}

func (s *OktaSAMLConnectionBootstraper) createStytchConnection(ctx context.Context, stytchConf *config.StytchSetupInput, organizationID string) (string, *config.StychSsoParameters, error) {
	sso, err := s.StytchClient.SSO.SAML.CreateConnection(ctx,
		&saml.CreateConnectionParams{
			DisplayName:    stytchConf.ConnectionDisplayName,
			OrganizationID: organizationID,
		},
	)

	if err != nil {
		return "", nil, err
	}

	return sso.Connection.ConnectionID, &config.StychSsoParameters{
		AcsUrl:   sso.Connection.AcsURL,
		Audience: sso.Connection.AudienceURI,
	}, nil
}

func (s *OktaSAMLConnectionBootstraper) updateStytchConnection(ctx context.Context, organizationID, connectionID string, conf *config.OktaSsoParameters) error {
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
				// This allows us to use implicit group assignements
				// ref: https://stytch.com/docs/b2b/guides/rbac/role-assignment
				"groups": "groups",
			},
		})

	return err
}
