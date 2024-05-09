/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stytchauth/stytch-go/v12/stytch/b2b/b2bstytchapi"
	"github.com/xNok/go-stytch-demo/pkg/config"
	"github.com/xNok/go-stytch-demo/pkg/server"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve start a HTTP server for this application",
	Long: `This server implement the Stytch Backend Integration of SSO
see https://stytch.com/docs/b2b/guides/sso/backend`,
	Run: RunServe,
}

func RunServe(cmd *cobra.Command, args []string) {
	v := viper.GetViper()

	clientConf, err := config.NewClientConfig(v)
	if err != nil {
		log.Fatalf("error loading client configs, did you forget to set environement varaibles? %s", err)
	}

	// Step 1: Instanciate stytch client
	stytchClient, err := b2bstytchapi.NewClient(
		clientConf.StytchConf.ProjectID,
		clientConf.StytchConf.Secret,
	)
	if err != nil {
		log.Fatalf("error instantiating API client %s", err)
	}

	conf, err := config.NewSetupResult(v)
	if err != nil {
		log.Fatalf("error reading config. Did you complete the setup? %s", err)
	}

	server.Serve(stytchClient, &server.StytchServerConfig{
		OrganizationID: conf.OrganizationID,
		ConnectionID:   conf.ConnectionID,
		PublicToken:    clientConf.StytchConf.PublicToken,
	})
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
