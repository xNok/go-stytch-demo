#!/bin/bash

curl --request POST \
	--url https://test.stytch.com/v1/b2b/organizations \
	-u "$STYTCH_PROJECT_ID:$STYTCH_SECRET" \
	-H 'Content-Type: application/json' \
	-d '{
		"organization_name": "Example Org Inc.",
		"organization_slug": "example-org"
	}'

curl --request POST \
	--url https://test.stytch.com/v1/b2b/sso/saml/${STYTCH_ORGANIZATION_ID} \
	-u "$STYTCH_PROJECT_ID:$STYTCH_SECRET" \
	-H 'Content-Type: application/json' \
	-d '{
	  "display_name": "Okta"
	}'

curl --request PUT \
	--url https://test.stytch.com/v1/b2b/sso/saml/${STYTCH_ORGANIZATION_ID}/connections/${STYTCH_CONNECTION_ID} \
	-u "$STYTCH_PROJECT_ID:$STYTCH_SECRET" \
	-H 'Content-Type: application/json' \
    -d "$(cat <<EOF
{
    "idp_entity_id": "${IDENTITY_PROVIDER_ISSUER}",
    "idp_sso_url": "${IDENTITY_PROVIDER_SINGLE_SIGN_ON_URL}",
    "x509_certificate": "${X_509_CERTIFICATE//$'\n'/\\\\n}",
    "attribute_mapping": {
        "email": "NameID",
        "first_name": "firstName",
        "last_name": "lastName"
    }
}
EOF
)"

curl -v -X POST \
	-H "Accept: application/json" \
	-H "Content-Type: application/json" \
	-H "Authorization: SSWS ${OKTA_API_TOKEN}" \
	-d '{
	"name": "okta_org2org",
	"label": "Sample Okta Org2Org App",
	"signOnMode": "SAML_2_0",
	"settings": {
		"app": {
		"acsUrl": "https://example.okta.com/sso/saml2/exampleid",
		"audRestriction": "https://www.okta.com/saml2/service-provider/exampleid",
		"baseUrl": "https://example.okta.com"
		}
	}
	}' "${OKTA_ORG_URL}/api/v1/apps"
