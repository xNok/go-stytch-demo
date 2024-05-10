package setup

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/okta/okta-sdk-golang/v4/okta"
	"github.com/xNok/go-stytch-demo/pkg/config"
)

func (s *OktaSAMLConnectionBootstraper) setupOktaSamlApplication(ctx context.Context, oktaConf *config.OktaSetupInput, stytchConf *config.StychSsoParameters) (string, error) {
	samlApp := okta.NewSamlApplication()
	samlApp.Label = &oktaConf.SAMLAppLabel
	samlApp.SignOnMode = okta.PtrString("SAML_2_0")
	samlApp.Visibility = okta.NewApplicationVisibility()
	samlApp.Settings = okta.NewSamlApplicationSettingsWithDefaults()
	samlApp.Settings.SignOn = &okta.SamlApplicationSettingsSignOn{
		// Data coming from Stytch
		SsoAcsUrl:   &stytchConf.AcsUrl,
		Audience:    &stytchConf.Audience,
		Recipient:   &stytchConf.AcsUrl,
		Destination: &stytchConf.AcsUrl,
		// Attributes you want to send to Stytch
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
			}, {
				FilterType:  okta.PtrString("REGEX"),
				FilterValue: okta.PtrString(".*billing.*"),
				Name:        okta.PtrString("groups"),
				Namespace:   okta.PtrString("urn:oasis:names:tc:SAML:2.0:attrname-format:unspecified"),
				Type:        okta.PtrString("GROUP"),
			},
		},
		// Default value required in the request
		DefaultRelayState:     okta.PtrString(""),
		IdpIssuer:             okta.PtrString("http://www.okta.com/${org.externalKey}"),
		SubjectNameIdTemplate: okta.PtrString("${user.userName}"),
		SubjectNameIdFormat:   okta.PtrString("urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress"),
		ResponseSigned:        okta.PtrBool(true),
		AssertionSigned:       okta.PtrBool(true),
		SignatureAlgorithm:    okta.PtrString("RSA_SHA256"),
		DigestAlgorithm:       okta.PtrString("SHA256"),
		HonorForceAuthn:       okta.PtrBool(true),
		AuthnContextClassRef:  okta.PtrString("urn:oasis:names:tc:SAML:2.0:ac:classes:PasswordProtectedTransport"),
		SpIssuer:              nil,
	}

	oktaApp, _, err := s.OktaClient.ApplicationAPI.CreateApplication(ctx).Application(
		okta.ListApplications200ResponseInner{
			SamlApplication: samlApp,
		},
	).Execute()

	if err != nil {
		return "", err
	}

	return *oktaApp.SamlApplication.Id, nil
}

func (s *OktaSAMLConnectionBootstraper) getOktaSamlApplicationMetada(ctx context.Context, oktaAppID string) (*config.OktaSsoParameters, error) {
	// Fetch the SAML metdata we need to configure Stych
	// The okta SDK is broken it does set the Content-Type as application/xml
	// metadata, _, err := s.OktaClient.ApplicationSSOAPI.PreviewSAMLmetadataForApplication(ctx, oktaAppID).Execute()

	metadata, err := previewSAMLmetadataForApplication(ctx, s.OktaClient, oktaAppID)

	if err != nil {
		return nil, err
	}

	// Parse the SAML metadata XML
	SAML, err := parseXML(metadata)

	if err != nil {
		return nil, err
	}

	result := &config.OktaSsoParameters{
		IdpEntityID:     SAML.EntityID,
		IdpSSOURL:       SAML.IDPSSODescriptor.SingleSignOnServices[0].Location,
		X509Certificate: SAML.IDPSSODescriptor.KeyDescriptors[0].KeyInfo.X509Data.X509Certificate,
	}

	return result, nil
}

// previewSAMLmetadataForApplicationm replace the oktaSDK that doesn't send the right header
// We are forced to used this until they fix their SDK
func previewSAMLmetadataForApplication(ctx context.Context, oktaClient *okta.APIClient, appId string) (string, error) {
	client := &http.Client{}

	url := "https://" + oktaClient.GetConfig().Host + fmt.Sprintf("/api/v1/apps/%s/sso/saml/metadata", appId)

	var key string
	if auth, ok := oktaClient.GetConfig().Context.Value(okta.ContextAPIKeys).(map[string]okta.APIKey); ok {
		if apiKey, ok := auth["API_Token"]; ok {
			if apiKey.Prefix != "" {
				key = apiKey.Prefix + " " + apiKey.Key
			} else {
				key = apiKey.Key
			}
		}
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)

	if err != nil {
		return "", err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/xml")
	req.Header.Set("Accept", "application/xml")
	req.Header.Set("Authorization", key)

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending request to API endpoint: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	return string(responseBody), nil
}
