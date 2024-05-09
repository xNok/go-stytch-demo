# How to set up RBAC in Go using Stytch

* [X] Set up Stych B2B SaaS Authentication with Okta as IDP
    * [Manual Process](https://stytch.com/docs/b2b/guides/sso/okta-saml)
    * [Automated With Go](./pkg/setup/main.go)
* [x] SSO Authentication Workflow with Styth and Go   
    * [Official Documentation](https://stytch.com/docs/b2b/guides/sso/backend)
    * [Automated With Go](./pkg/server/main.go)
* [X] Set up RBAC with Stytch
    * [Official Documentation](https://stytch.com/docs/b2b/guides/rbac/role-assignment)
    * [X] Automatic role assignment based on metadata
    * [X] Set up Stytch default resources and custom roles
    * [X] Set up authorization checks for custom resources


## Using the Demo CLI

This demo project can be used a CLI to quicly experiment Stytch. Simply `install` the module to get started.

```bash
go install github.com/xNok/go-stytch-demo@main
```

## Set up Stych B2B SaaS Authentication with Okta as IDP

To complete this tutorial, you will need the following

- A Stytch Account: [Sign Up to Stytch](https://stytch.com/dashboard/start-now)
    - Take note of your Project ID and create a new Secret [here (dashboard/api-keys?env=test)](https://stytch.com/dashboard/api-keys?env=test)
- An Okta Account: [Okta free trial](https://www.okta.com/free-trial)
    - Take note of your Organisation's URL (It's the URL you are redirected to once you log in your Okta trial Okta) and [Create an API token](https://developer.okta.com/docs/guides/create-an-api-token/main/).
- A working Go environment: [Get Started with Golang](https://go.dev/).
    - You can use [GitPod](gitpod.io), I have added a `.gitpod.yml` to make things easier.

First setup you environment by exporting and defining your `STYTCH` and `OKTA` credential as follow:

```bash
STYTCH_PROJECT_ID="project-test-xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
STYTCH_SECRET="secret-test-xx-xxxxxxxxxxxxxxxxxxxxxxxxxxxx-xxx="
STYTCH_PROJECT_PUBLIC_ID="public-token-test-xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
OKTA_API_TOKEN="xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
OKTA_ORG_URL=https://trial-xxxxxxx-admin.okta.com/"
```

Next call the setup command to bootstrap the SSO SAML configuration between Stytch and Okta

```bash
go-stytch-demo setup
```

## Run local server

Now you can test that everything is working by running the local server. Make sure you redurect url is properly setup (http://localhost:8010/authenticate). Then run the following command:

```
go-stytch-demo serve
```

Go to http://localhost:8010 to start the authentication workflow, you should be redirected to Okta for login then back to you application.

## Configure RBAC

Now lets play with a few different features. Keep the server running and open a new terminal.

```
go-stytch-demo config [args]
```