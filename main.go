package main

import (
	"context"

	"github.com/sethvargo/go-envconfig"
	"github.com/stytchauth/stytch-go/v12/stytch/consumer/stytchapi"
	"github.com/stytchauth/stytch-go/v12/stytch/consumer/users"
)

type StytchConf struct {
	ProjectID string `env:"STYTCH_PROJECT_ID"`
	Secret    string `env:"STYTCH_SECRET"`
}

func main() {
	ctx := context.Background()

	var c StytchConf
	if err := envconfig.Process(ctx, &c); err != nil {
		panic(err)
	}

	stytchAPIClient, err := stytchapi.NewClient(
		c.ProjectID,
		c.Secret,
		stytchapi.WithBaseURI("https://test.stytch.com"),
	)

	if err != nil {
		panic(err)
	}

	_, err = stytchAPIClient.Users.Search(
		ctx, &users.SearchParams{
			Limit: 10,
		},
	)

	if err != nil {
		panic(err)
	}

}
