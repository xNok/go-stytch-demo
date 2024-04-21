/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"log"

	"github.com/sethvargo/go-envconfig"
	"github.com/spf13/cobra"
	"github.com/stytchauth/stytch-go/v12/stytch/b2b/b2bstytchapi"
	"github.com/xNok/go-stytch-demo/pkg/server"
	"github.com/xNok/go-stytch-demo/pkg/setup"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: RunServe,
}

func RunServe(cmd *cobra.Command, args []string) {
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

	// will be replace by viper soon
	confProvider := &setup.YAMLEntry{
		Path: "./setup.yaml",
	}

	// Load our configuration file, this file is empty if we start from scrath
	conf, err := confProvider.Load()
	if err != nil {
		log.Fatalf("error loading Configuration %s", err)
		return
	}

	server.Serve(stytchClient, &server.StytchServerConfig{
		OrganizationID: conf.OrganizationID,
	})

	return nil
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
