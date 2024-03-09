# LEETA TECHNOLOGY CORE BACKEND

This repository is responsible for handling all backend related requests for ALL LEETA Technology API implementations

## Dependencies

- Docker (Running containerized version): [Follow this link for getting started with docker and docker installation ](https://docs.docker.com/get-started/).
- MongoDB: [Follow this link for getting started with MongoDB](https://www.mongodb.com/)

## Running & Debugging

### Run containerized version
To run the service in a docker container, Docker should be installed  Then run the make command

```shell
make all 
```

This will execute the following make commands. See Makefile for command details

```text
generate_keys check_docker check_mongodb create_user check_database generate_docs run_app
```

#### MongoDB

To stop the running mongoDB container

```shell
make stop-mongo
```

### Run non-containerized version 
Alternatively, you can choose to start the dependencies individually and run the `Go` service with local `.env` file for a better experience during development.
- Run `MongoDB` in a docker container, 
- and start the go service with the command:

```shell
go run ./cmd/main.go -c ../testlocal.env
```

### Application Health

The application runs on the default port`:3000`

[Click here to check application Health](http://localhost:3000/health)