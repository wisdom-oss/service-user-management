<!-- TODO: REMOVE BLOCK AFTER READING README -->
> [!CAUTION]
> This repository either is the template for new microservices or this
> repository has been generated from the template repository.
>
> Please read the Architecture section thoroughly to minimize risks of data
> leaking and unauthorized access.

<div align="center">
<img height="150px" src="https://raw.githubusercontent.com/wisdom-oss/brand/main/svg/standalone_color.svg">

<!-- TODO: Change Information here -->

<h1>Microservice Template/Example</h1>
<h3>service-example</h3>
<p>📐 A minimal working example for microservices in the WISdoM Architecture</p>

<!-- TODO: Change URL here to point to correct repository -->
<img src="https://img.shields.io/github/go-mod/go-version/wisdom-oss/microservice-template?style=for-the-badge" alt="Go Lang Version"/>
<a href="openapi.yaml">
<img src="https://img.shields.io/badge/Schema%20Version-3.0.0-6BA539?style=for-the-badge&logo=OpenAPI%20Initiative" alt="Open
API Schema Version"/></a>
</div>

<!-- TODO: Replace README.md contents with correct description -->


 
## Architecture
The template repository contains the basic code fragments to configure and
start up a new microservice.
It is delivered with two basic routes showcasing the features of the template
and the current packages included it.

This section explains some files and their function within the service to allow
a better understanding of the service

### [`init.go`](init.go) — The service initialization
The `init.go` file contains the code required to connect the service to the
PostgreSQL database used in the WISdoM system. It furthermore handles the
automatic configuration of the environment variables by automatically loading
environment variables stored in a `.env` file.

> [!CAUTION]
> Never commit a file containing secrets like usernames, password, 
> client _secrets_, client ids, connection urls that contain these values.

### [`main.go`](main.go) — Main Application
The `main.go` file contains the main part of the application. 
In this case it is the setup of the healthcheck and the router used to manage
the different routes and handlers.


### [`globals`](globals) — Globally available variables and connections
The package `globals` manages some variables that are used at places which would
require importing each other resulting in a circular import.

#### [`variables.go`](globals/variables.go) — Variables
The `variables.go` file contains globally used variables that are used at
multiple places in different packages.
Furthermore, the `variables.go` file uses the [embedding of values] during the 
build of the executable.

[embedding of values]: https://pkg.go.dev/embed

> [!IMPORTANT]
> To set the service's name and to allow a first build, please write the name of
> the service into the `service-name.sample` file and rename this file to
> `service-name`.
> Otherwise, the build process will fail!

#### [`connections.go`](globals/connections.go) — Connections
The `connections.go` file contains globally used variables that are used at
multiple places throughout the code.

### [`resources`](resources) — Resources
The `resources` folder contains resources needed for the service.
In its bare state, the service requires the listed files

#### [`environment.json`](resources/environment.json) — Environment Configuration
The `environment.json` file contains information about the required and
optional environment variables consumed by the service.
For optional variables you need to specify a default value which is populated
into the `globals.Environment` variable.
The optional values may only be of the type `string` as this reflects the
behavior of the `os.LookupEnv` function.


#### [`queries.sql`](resources/queries.sql) — SQL Queries
The `queries.sql` file contains all sql queries required for your service.
The queries are loaded during the initialization of the service and are managed
by the [`dotsql` package]

[`dotsql` package]: https://pkg.go.dev/github.com/qustavo/dotsql

### [`config`](config) — Default configurations
The `config` folder contains Go files which set default values for the 
executable as constants or other variables.
The different files are used in dependency of the [build tags] supplied during
the compilation of the service.

[build tags]: https://pkg.go.dev/cmd/go#hdr-Build_constraints

#### [`defaults.go`](config/defaults.go) - Local Development
The `defaults.go` file is always used if the build tag `docker` is not supplied
to the compiler.
This default configuration automatically disables the authentication and
authorization measures preconfigured to allow easy local development.
It also sets the location of some required files to the `resources` folder to
match the repository layout.

#### [`defaults.docker.go`](config/defaults.docker.go) - Deployment
> [!NOTE]
> This build tag is pre-set in the [Dockerfile](Dockerfile). Therefore, no
> intervention is required.

The `defaults.docker.go` file is used if the build tag `docker` has been 
supplied to the compiler.
This configuration automatically enables the authentication and
authorization measures preconfigured to secure the service in its deployed
state.
It also sets the location of some required files to lay directly next to the
executable to allow a flatter folder structure in the generated docker
image.

