# Users gRPC

Simple user management microservice.

## Goals of this project:
- Learn about microservices
- Learn about gRPC
- Be able to use this service along with other microservice (have an API Gateway, a frontend in React, and interact with a To-Dos service.

## To-do:

- [x] Creates the migrations/seeds for the database.
- [x] Request for login, returning JWT token and Refresh token, that the API Gateway will use.
- [x] Refresh the access tokens, if current refresh token is valid.
- [x] Register a new user with a certain role (needs to have a valid JWT).
- [ ] CRUD For the users and roles (incomplete right now).
- [x] Add proper gRPC error handling.
- [x] Unit Test everything.
- [ ] Add integration tests.

## Installation
### Requirements
1. Have docker and docker-compose installed.

### Process to install
```
docker network create -d bridge todos-infrastructure
cd ~/workspace
git clone git@github.com:plagioriginal/user-api-todos.git
cd user-api-todos
cp .env.example .env
docker-compose up
```
You will have available a gRPC service running at http://localhost:8080 (port is editable in the `.env` file).

## Overall project architecture

### API Gateway
- REST API (returns in JSON)
- A client for all the gRPC microservices (reads all the gRPC services).
- Responsible for the interactions between all the microservices.

### Users gRPC
- Manages the users in the architecture

### To-dos gRPC
- Manages all the to-dos.
- Maybe will be used for note-taking as well.

### Frontend
- In React, will only call and use data from the API Gateway.
