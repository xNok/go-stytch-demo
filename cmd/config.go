/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stytchauth/stytch-go/v12/stytch/b2b/b2bstytchapi"
	"github.com/xNok/go-stytch-demo/pkg/config"
	"github.com/xNok/go-stytch-demo/pkg/rbac"
)

const (
	flagOrgImpAss     = "organisation-implicit-assignment"
	flagConImpAss     = "connection-implicit-assignment"
	flagConSAMLImpAss = "connection-saml-implicit-assignment"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Config applies new configuration to a Stytch connection",
	Long: `This command provide variaous flags that lets you test various scenarios.

Including:
* Automatic role assignment based on metadata
* Set up Stytch default resources and custom roles
* Set up authorization checks for custom resources
`,
	RunE: RunConfig,
}

func RunConfig(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	v := viper.GetViper()

	clientConf, err := config.NewClientConfig(v)
	if err != nil {
		return fmt.Errorf("error loading client configs, did you forget to set environement varaibles? %s", err)
	}

	// Step 1: Instanciate stytch client
	stytchClient, err := b2bstytchapi.NewClient(
		clientConf.StytchConf.ProjectID,
		clientConf.StytchConf.Secret,
	)
	if err != nil {
		return fmt.Errorf("error instantiating API client %s", err)
	}

	conf, err := config.NewSetupResult(v)
	if err != nil {
		log.Fatalf("error reading config. Did you complete the setup? %s", err)
	}

	stytchRBACConfig := &rbac.StytchRBACConfig{
		OrganizationID: conf.OrganizationID,
		ConnectionID:   conf.ConnectionID,
		Domain:         "devops-family.com",
	}

	if flag, _ := cmd.Flags().GetBool(flagOrgImpAss); flag {
		cmd.Println("ApplyOrganizationImplictAssignement")
		if err = rbac.ApplyOrganizationImplictAssignement(ctx, stytchClient, stytchRBACConfig); err != nil {
			return err
		}
	}

	if flag, _ := cmd.Flags().GetBool(flagConImpAss); flag {
		cmd.Println("ApplyOrganizationImplictAssignement")
		if err = rbac.ApplyConnectionImplictAssignement(ctx, stytchClient, stytchRBACConfig); err != nil {
			return err
		}
	}

	if flag, _ := cmd.Flags().GetBool(flagConSAMLImpAss); flag {
		cmd.Println("ApplyOrganizationImplictAssignement")
		if err = rbac.ApplyConnectionSAMLGroupImplictAssignement(ctx, stytchClient, stytchRBACConfig); err != nil {
			return err
		}
	}

	return nil
}

func init() {
	rootCmd.AddCommand(configCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	configCmd.Flags().BoolP(flagOrgImpAss, "o", false, "Setup Stytch Organisation implicit role assignement")
	configCmd.Flags().BoolP(flagConImpAss, "p", false, "Setup Stytch Connection implicit role assignement")
	configCmd.Flags().BoolP(flagConSAMLImpAss, "q", false, "Setup Stytch Connection SAML Group implicit role assignement")

}
