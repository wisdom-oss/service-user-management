<div align="center">
<img height="150px" src="https://raw.githubusercontent.com/wisdom-oss/brand/main/svg/standalone_color.svg">
<h1>User Management</h1>
<h3>service-user-management</h3>
<p>👥 user and permission management using OpenID Connect</p>
<img src="https://img.shields.io/github/go-mod/go-version/wisdom-oss/service-user-management?style=for-the-badge" alt="Go Lang Version"/>
<a href="openapi.yaml">
<img src="https://img.shields.io/badge/Schema%20Version-3.0.0-6BA539?style=for-the-badge&logo=OpenAPI%20Initiative" alt="Open
API Schema Version"/></a>
</div>

> [!IMPORTANT]
> This microservice depends on an external OpenID Connect Provider

This microservice acts as a middle-man between the OpenID Connect Provider used
for authentication of users and the permission management inside the WISdoM
platform to minimize the amount of customization required in a OpenID Connect
Provider.
It accepts the authentication codes generated by the OpenID Connect Provider
and uses them to request an ID Token from the provider.
This ID token is then used to provide an access token which allows 
authenticating with backend services and allows dynamically showing and hiding 
entries in the frontend.
The user management service uses signed JWTs to ensure that no tampering can
happen on the client side to gain access to services without proper 
authorization.

## Configuration
The microservice requires access to a PostgreSQL database for storing the
external identifiers of users and for persisting permission information about
the users.
To connect the microservice to a database, please set the following environment
variables:
  - `PGUSER`
  - `PGPASSWORD`
  - `PGHOST`
  - `PGDATABASE`

and if necessary:
  - `PGPORT`

Furthermore, you need to specify the client id and secret for the OpenID Connect
provider as well as the issuer as shown in your provider using the following
environment variables:
  - `OIDC_CLIENT_ID`
  - `OIDC_CLIENT_SECRET`
  - `OIDC_ISSUER`

The required certificates are automatically generated during the initial startup
and stored in the microservice.
It is recommended to create a volume mount if using docker to persist the
certificates during container recreation to ensure updates to not break already
running sessions