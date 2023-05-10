# GraphGenerator

This application defines a REST server with implementation of random graph generator.
The server is implemented in language Go, so in order to run it, you will need 

* Go: [Installation guide](https://go.dev/doc/install)
* NodeJS: [Installation](https://nodejs.dev/en/learn/how-to-install-nodejs/)
* Poetry: [Installation guide](https://python-poetry.org/docs/) (for testing purposes)


## Building the application

To build the application you can run command

```bash
make build
```

## Running

To run the development mode you can run

```bash
go run ./cmd/generator -devel
```

To start the ui project in development you'll need to run

```bash
cd ui 
npm install
npm start
```

These steps are automated in the helper script `scripts/run_devel.sh`.

To run in the all-in-one session (e.g. the server also serves the ui artefacts),
you'll just need to run

```bash
./generator -devel
```

This can be done only after the build was performed.

### Configuration options

Here follows table with available configuration options for the generator application.

| Variable                         | Option                | Value      | Description                                                                                                                       |
|----------------------------------|-----------------------|------------|-----------------------------------------------------------------------------------------------------------------------------------|
| `GENERATOR_MAX_NODES`            | `maxNodes`            | Integer    | Sets maximal allowed number of nodes in graph request.                                                                            |
| `GENERATOR_MAX_BATCH_SIZE`       | `maxBatchSize`        | Integer    | Sets maximal allowed number of graphs in batch.                                                                                   |
| `GENERATOR_WORKERS`              | `workers`             | Integer    | Sets number of worker goroutines used for graph generation.                                                                       |
| `GENERATOR_DB_ROOT`              | `dbRoot`              | Path       | Sets where the root of database will be located.                                                                                  |
| `GENERATOR_MAINTENANCE_INTERVAL` | `maintenanceInterval` | Interval   | Sets the interval between two instances of database maintenance. For format see [here](#interval-format)                          |
| `GENERATOR_MAINTENANCE_HOUR`     | `maintenanceHour`     | Integer    | Sets the base hour to which schedule maintenance.                                                                                 |
| `GENERATOR_REQ_TTL`              | `ttl`                 | Interval   | Sets the liveness timespan of one request.                                                                                        |
| `GENERATOR_LOG_LEVEL`            | `logLevel`            | Log string | Defines logging level, is one of `TRACE`, `DEBUG`, `INFO`, `WARNING` or `ERROR`.                                                  |
| `GENERATOR_BIND_HOST`            | `bindHost`            | String     | Sets local interface to which will the server bind to.                                                                            |
| `GENERATOR_PORT`                 | `port`                | Integer    | Sets the expected port number on which will te  server listen for new connections.                                                |
| `GENERATOR_HOST`                 | `host`                | String     | Sets expected domain name of application.                                                                                         |
| `SECURE_MODE`                    | `secure`              | Boolean    | Binary switch which sets the server to provide secure cookies and set them with configured host, for env var boolean is expected. |
| -                                | `devel`               | -          | Binary option sets server with development default values.                                                                        |
| -                                | `test`                | -          | Binary option sets server with testing defaults.                                                                                  | 
| `GENERATOR_UI_LOCATION`          | -                     | Path       | Instructs server to provide UI from specified directory.                                                                          |


#### Interval format

Interval format is specified in [Go Duration](https://pkg.go.dev/time#ParseDuration).
It is number closely followed by unit of time, e.g. one hour is `1h`, 30 seconds is `30s`.

## Shipping

To ship this project as docker image, you need [Docker](https://docs.docker.com/get-started/) or [Podman](https://podman.io/get-started) installed,
then you can build this project into container image by running

```bash
docker build . --tag my_generator:latest
```

or in case of Podman:

```bash
podman build . --tag my_generator:latest
```

For convenience there is also `make` target for automation of this step.
Let's assume that variable `CONTAINRER_ENGINE` exists and is set to either `docker` or `podman`,
and variable `IMAGE_NAME` contains tag name under which should the image be built.

```bash
make dist CONT="${CONTAINER_ENGINE}" IMAGE="${IMAGE_NAME}"
```

## Running as a container

To run application from container image, you'll can run following command:
(assuming you're running the application with docker and image is build under tag `generator:latest`)

```bash
docker run -p 8080:8080 generator:latest -devel
```

The command will start application with port forwarded port 8080 of a container. 
This will start the application in development mode with ephemeral database - e.g. the
database will be stored only for the time of existence of created container.

To create persistent storage, create volume in advance

```bash
docker create volume generator
```

The command above will create volume named generator which can be later used in the option `-v` of docker run
command.

To pass configuration in reasonable manner, you can use EnvFile, in the file `env.txt.tpl` you can see template
of such file. To create configuration for generator container, just copy the file to `env.txt` and 
fill in the variables you want to set up, removing the ones you don't.

To run container with the configuration from env file with volume just run command

```bash
docker run -p 8080:8080 -v generator:/var/badger/db --env-file=./env.txt --name generator generator:latest
```

The command will create and start new container with configuration from env.txt, and database data will be persisted
in the created volume. If you're using Podman, just change `docker` command with `podman`.


## Cleaning

To remove building residues you can run 

```shell
make clean
```

If you need to remove all the downloaded dependencies, you can run:

```shell
make full-clean
```
