# How to set up RBAC in Go using Stytch

* [X] Set up Stych B2B SaaS Authentication with Okta as IDP
    * [Manual Process](https://stytch.com/docs/b2b/guides/sso/okta-saml)
    * [Automated With Go](./pkg/setup/setup.go) 
* [~] SSO Authentication Workflow with Styth and Go   
    * [](https://stytch.com/docs/b2b/guides/sso/backend)
* [] Set up RBAC with Stytch
    * [] Automatic role assignment based on metadata
    * [] Set up Stytch default resources and custom roles
    * [] Set up authorization checks for custom resources


## Set up Stych B2B SaaS Authentication with Okta as IDP

```
go run main.go setup
```

## Run local server


```
go run main.go serve
```

## Configure RBAC

```
go run main.go config [args]
```