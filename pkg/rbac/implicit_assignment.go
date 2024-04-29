package rbac

import (
	"context"

	"github.com/stytchauth/stytch-go/v12/stytch/b2b/b2bstytchapi"
	"github.com/stytchauth/stytch-go/v12/stytch/b2b/organizations"
	"github.com/stytchauth/stytch-go/v12/stytch/b2b/sso"
	"github.com/stytchauth/stytch-go/v12/stytch/b2b/sso/saml"
)

// ref https://stytch.com/docs/b2b/guides/rbac/role-assignment#implicit-assignment

type StytchRBACConfig struct {
	OrganizationID string
	ConnectionID   string
	Domain         string
}

// By email domain: everyone with the stytch.com email domain gets the “developer” Role.
func ApplyOrganizationImplictAssignement(ctx context.Context, stytchClient *b2bstytchapi.API, conf *StytchRBACConfig) error {
	_, err := stytchClient.Organizations.Update(ctx, &organizations.UpdateParams{
		OrganizationID: conf.OrganizationID,
		RBACEmailImplicitRoleAssignments: []*organizations.EmailImplicitRoleAssignment{
			{
				Domain: conf.Domain,
				RoleID: "developer",
			},
		},
	})

	if err != nil {
		return err
	}

	return nil
}

// By SSO Connection: everyone who authenticates via a specific SSO Connection gets the “employee” Role.
func ApplyConnectionImplictAssignement(ctx context.Context, stytchClient *b2bstytchapi.API, conf *StytchRBACConfig) error {
	_, err := stytchClient.SSO.SAML.UpdateConnection(ctx, &saml.UpdateConnectionParams{
		OrganizationID: conf.OrganizationID,
		ConnectionID:   conf.ConnectionID,
		SAMLConnectionImplicitRoleAssignments: []*sso.SAMLConnectionImplicitRoleAssignment{
			{
				RoleID: "employee",
			},
		},
	})

	if err != nil {
		return err
	}

	return nil
}

// By SSO Connection IdP Group: everyone who authenticates via a specific SSO Connection and is a part of the “engineering” IdP group gets the “developer” Role.
func ApplyConnectionSAMLGroupImplictAssignement(ctx context.Context, stytchClient *b2bstytchapi.API, conf *StytchRBACConfig) error {
	_, err := stytchClient.SSO.SAML.UpdateConnection(ctx, &saml.UpdateConnectionParams{
		OrganizationID: conf.OrganizationID,
		ConnectionID:   conf.ConnectionID,
		SAMLGroupImplicitRoleAssignments: []*sso.SAMLGroupImplicitRoleAssignment{
			{
				RoleID: "billing",
				Group:  "billing",
			},
		},
	})

	if err != nil {
		return err
	}

	return nil
}
