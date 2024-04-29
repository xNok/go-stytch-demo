# How to set up RBAC in Go using Stytch

* [X] Set up Stych B2B SaaS Authentication with Okta as IDP
    * [Manual Process](https://stytch.com/docs/b2b/guides/sso/okta-saml)
    * [Automated With Go](./pkg/setup/main.go)
* [x] SSO Authentication Workflow with Styth and Go   
    * [Official Documentation](https://stytch.com/docs/b2b/guides/sso/backend)
    * [Automated With Go](./pkg/server/main.go)
* [~] Set up RBAC with Stytch
    * [Official Documentation](https://stytch.com/docs/b2b/guides/rbac/role-assignment)
    * [X] Automatic role assignment based on metadata
    * [X] Set up Stytch default resources and custom roles
    * [] Set up authorization checks for custom resources


## Using the Demo CLI

This demo project can be used a CLI to quicly experiment Stytch. Simply `install` the module to get started.

```bash
go install github.com/xNok/go-stytch-demo@main
```


## Set up Stych B2B SaaS Authentication with Okta as IDP

```
go-stytch-demo setup
```

## Run local server


```
go-stytch-demo serve
```

## Configure RBAC

```
go-stytch-demo config [args]
```