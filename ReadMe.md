# LEETA TECHNOLOGY CORE BACKEND

This repository is responsible for handling all backend related requests for ALL LEETA Technology API implementations 

## Dependencies

Docker is needed to run the backend service. [Follow this link for getting started with docker and docker installation ](https://docs.docker.com/get-started/).

```shell
make all 
```

This will execute the following make commands. See Makefile for command details

```text
generate_keys check_docker check_mongodb create_user check_database generate_docs run_app
```

### MongoDB

To stop the running mongoDB container

```shell
make stop-mongo
```


### Application Health

[Click here to check application Health](http://localhost:3000/health)