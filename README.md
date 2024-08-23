<div align="center">
<img height="150px" src="https://raw.githubusercontent.com/wisdom-oss/brand/main/svg/standalone_color.svg">
<h1>User Management</h1>
<h3>service-user-management</h3>
<p>👥 user and permission management using OpenID Connect</p>

<!-- TODO: Change URL here to point to correct repository -->
<img src="https://img.shields.io/github/go-mod/go-version/wisdom-oss/service-user-management?style=for-the-badge" alt="Go Lang Version"/>
<a href="openapi.yaml">
<img src="https://img.shields.io/badge/Schema%20Version-3.0.0-6BA539?style=for-the-badge&logo=OpenAPI%20Initiative" alt="Open
API Schema Version"/></a>
<br/>
<img height="28" src="https://jwt.io/img/badge-compatible.svg"/>
<img src="https://img.shields.io/badge/OpenID%20Connect-Compatible-F78C40?style=for-the-badge&logo=openid"/>
</div>

> [!CAUTION]
> This microservice is required for the WISdoM platform as it is the
> authorization management for the APIs.
> Excluding it from the deployment will result in an inaccessible API and
> platform

The user management service integrates an external OpenID Connect provider
allowing external authentication and keeping authorization in the scope of the
WISdoM platform.
This allows for a granular access control in the frontend that is independent
of the selected OpenID Connect provider as this service generates the access
tokens for the WISdoM platform.
