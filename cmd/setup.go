/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"log"

	"github.com/okta/okta-sdk-golang/v4/okta"
	"github.com/sethvargo/go-envconfig"
	"github.com/spf13/cobra"
	"github.com/stytchauth/stytch-go/v12/stytch/b2b/b2bstytchapi"
	"github.com/xNok/go-stytch-demo/pkg/setup"
)

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "A utility script to create the SAML connection between stycth and okta",
	Long: `This setup will create a new stych organisation and connection,
Then create a new okta application and and finally proceed with the SAML metadata exchange.`,
	RunE: RunSetup,
}

func RunSetup(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	var c setup.Conf
	if err := envconfig.Process(ctx, &c); err != nil {
		panic(err)
	}

	// Step 1: Instanciate stytch client
	stytchClient, err := b2bstytchapi.NewClient(
		c.StytchConf.ProjectID,
		c.StytchConf.Secret,
	)

	if err != nil {
		log.Fatalf("error instantiating API client %s", err)
	}

	// Step 2: Instanciate Okta client
	oktaConfig, err := okta.NewConfiguration(
		okta.WithOrgUrl(c.OktaConf.OrgUrl),
		okta.WithToken(c.OktaConf.APIToken),
	)

	if err != nil {
		log.Fatalf("error instantiating Okta API client %s", err)
	}

	oktaClient := okta.NewAPIClient(oktaConfig)

	bootstraper := setup.NewOktaSAMLConnectionBootstraper(stytchClient, oktaClient)
	_, err = bootstraper.Setup(ctx)

	return err
}

func init() {
	rootCmd.AddCommand(setupCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setupCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setupCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
