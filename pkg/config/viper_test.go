package config

import (
	"reflect"
	"strings"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestNewClientConfig(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(t *testing.T)
		want    *ClientsConf
		wantErr bool
	}{
		{
			name: "",
			setup: func(t *testing.T) {
				t.Setenv("STYTCH_PROJECT_ID", "1234")
				t.Setenv("STYTCH_SECRET", "12345")
			},
			want: &ClientsConf{
				StytchConf: &StytchConf{
					ProjectID: "1234",
					Secret:    "12345",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(t)
			v := setupViper(t)
			got, err := NewClientConfig(v)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClientConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewClientConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewSetupResult(t *testing.T) {
	tests := []struct {
		name       string
		configPath string
		want       *SetupResult
		wantErr    bool
	}{
		{
			name:       "",
			configPath: "setup",
			want: &SetupResult{
				StytchResult{
					OrganizationID: "organization-test-4f8867e0-d973-40b1-83ab-631ca1e8494d",
					ConnectionID:   "saml-connection-test-b5fc0295-fdde-4452-a210-56e5ab44ed59",
					SsoParameters: &StychSsoParameters{
						AcsUrl:   "https://test.stytch.com/v1/b2b/sso/callback/saml-connection-test-b5fc0295-fdde-4452-a210-56e5ab44ed59",
						Audience: "https://test.stytch.com/v1/b2b/sso/callback/saml-connection-test-b5fc0295-fdde-4452-a210-56e5ab44ed59",
					},
				},
				OktaResult{
					ApplicationID: "0oada53uqsswV59o9697",
					SsoParameters: nil,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := setupViper(t)
			v.SetConfigName(tt.configPath)
			got, err := NewSetupResult(v)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewSetupResult() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSetupResult() = %v, want %v", got, tt.want)
			}
		})
	}
}

func setupViper(t *testing.T) *viper.Viper {
	t.Helper()

	v := viper.New()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
	v.AddConfigPath("testdata")
	v.SetConfigType("yaml")
	v.SetConfigName("setup.yaml")

	err := v.ReadInConfig()
	require.NoError(t, err)

	return v
}
