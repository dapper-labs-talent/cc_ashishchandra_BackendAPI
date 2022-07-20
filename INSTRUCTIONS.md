# Instructions for the User registration API

## Prerequisites

The following applications/libraries should be installed on the machine where this application is deployed, tested and evaluated.

1. docker (version 20 or higher) and docker-compose (version 1.2 or higher) [Get Docker](https://docs.docker.com/get-docker/) and [Install Docker Compose](https://docs.docker.com/compose/install/)
2. Go 1.17 or higher [Download and install](https://go.dev/doc/install)
3. ssh-keygen (This is present on almost all Linux/Unix flavors)
4. psql client (version 12 or higher) [Download Postgres](https://www.postgresql.org/download/)
5. git (version 2+) [This is standard on Linux but not on the Windows platform so you can install Git Bash] [Download GitBash](https://git-scm.com/downloads)

## Getting the code from GitHub

In a terminal window (preferably Linux), type

1. `git clone` https://github.com/dapper-labs-talent/cc_ashishchandra_BackendAPI.git
2. cd cc_ashishchandra_BackendAPI

## Starting the Server

The following command builds the docker image for the server and also downloads the Postgres 12 image from DockerHub


`/fakepath/cc_ashishchandra_BackendAPI$ ./start.sh`

If everything worked, it should echo out a message at the end of the script like this one:

**Successfully started all containers. Application is ready on port 3000**

Now you can invoke the APIs as listed in the README.md

## Testing the APIs

This application uses the Go testing framework to execute unit test cases.

To execute all the tests:

Ensure that the postgres container is up and running or the tests will not work.

To check:
`docker ps | grep postgres`

If the postgres container is running, you should see output somewhat like this:
`37034a0bf34e   postgres:12   "docker-entrypoint.s…"   2 minutes ago   Up 2 minutes   0.0.0.0:5432->5432/tcp, :::5432->5432/tcp   postgres`

If the Postgres container is not running, you can start it with:
`/fakepath/cc_ashishchandra_BackendAPI$ docker-compose up -d postgres`

Then run the tests:

`/fakepath/cc_ashishchandra_BackendAPI$ POSTGRES_USER=dapperlabs POSTGRES_PASSWORD=dapperlabs123 POSTGRES_PORT=5432 POSTGRES_DB=dapperlabs POSTGRES_HOST=localhost go test`