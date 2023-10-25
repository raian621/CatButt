# CatBook

A social media website made specifically for cat pictures: No humans allowed!

---

## Prereq Dependencies:

Download:
- [Golang](https://go.dev/doc/install)
- [Air](https://github.com/cosmtrek/air#via-go-install-recommended)
- [Node JS](https://docs.npmjs.com/downloading-and-installing-node-js-and-npm)

The included `Makefile` should be able to install package dependencies automatically, but if you'd like to install them manually then run:

```sh
$ cd site; npm i          # install site dependencies
$ cd server; go get ./... # install server dependencies
```

## Setting Up Development Environment

A `Makefile` is included in the root of the project directory is included that can set up the development environment:

To set up the develoment environment, simply run:

```sh
# set up docker test database:
$ make db       # builds docker database container; only needs to run once
$ make start-db # starts docker database container

# install source dependencies for the site and the server
$ make deps
```

Before running the project in development mode, you must set all the non-optional environment variables either by adding a `.env` file to both the `server` and `site` folders or by defining them on the current terminal using the `export` command:
- [Server environment variables](#server-env-vars)
- [Site environment variables](#site-env-vars)

## Running the Project in Development Mode

Again, we can use the included `Makefile` to run the `site` and `server`:
```sh
$ make dev-server # start the server in development mode (server recompiles on save)
$ make dev-site   # start the site in development mode (hot reload)
```

---
## `site`
### Environment Variables<a id='site-env-vars'></a>
Variable       | Optional | Default Value | Use
---------------|----------|---------------|--------------------------
`VITE_API_URL` | No       | n/a           | URL that the backend server service is hosted on

## `server`

### Environment Variables:<a id='server-env-vars'></a>

Variable  | Optional | Default Value | Use
----------|----------|---------------|--------------------------
`DB_PROV` | No       | n/a           | Database provider (only possible value is currently `postgres`)
`DB_USER` | No       | n/a           | Database user for logging into database provider
`DB_NAME` | No       | n/a           | Name of the database you wish to connect to
`DB_PASS` | No       | n/a           | Password of the account used to log into the database provider
`DB_HOST` | No       | n/a           | Hostname of the host that the database provider is on
`DB_PORT` | No       | n/a           | Port number the database provider is listening on
`DB_SSL`  | No       | n/a           | `disable` if the database provider is not using SSL, otherwise `enable`
`HOST`    | Yes      | `localhost`   | Hostname to host server service on
`PORT`    | Yes      | `8080`        | Port to host the server service on
