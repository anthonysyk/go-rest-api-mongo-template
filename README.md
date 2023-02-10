# REST API with MongoDB template

## Purpose

REST API implementing a CRUD on MongoDB and a login.

## Usage

Setup local environment with docker-compose

```make compose-up```

Tear down local environment with docker-compose

```make compose-down```

Run server locally

```make run```

Run tests

```make test```

## Deploy

### Local

We use docker-compose to setup local environment (mongodb and mongo-express)

### Env variables

Use `.env.example` to generate your `.env` file

## Authentication

User must be authorized with a JWT token to access to his data.

```Authorization: Bearer jwt-token```

### Routes

These routes require authentication

```
GET /users/:id
PUT /users/:id
DELETE /delete/user/:id
```

## Structure

```
├── api
│   └── user
│     ├── handler.go
│     ├── model.go
│     ├── repository.go
│     └── service.go
├── config
├── internal
│   ├── auth
│   ├── db
│   ├── fs
│   ├── jsonflex
│   ├── middleware
│   └── password
└── test
    └── integration
```

- api: contains REST API routes
    - user: contains handler, service, repository, model
- config: contains config for the API
- internal: contains packages not to be shared
- test: contains tests
    - integration: contains integration tests 

