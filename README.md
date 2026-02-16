# Task Manager

## Description

A simple task manager with a pre-defined set of features.

### Features

- todo
  - Add new task
  - Update task
    - Description changes
    - Status changes
      - Created
      - Started
      - Canceled
      - Done
  - Show progress details
  - Delete task
  - Archive task
- user profile
- login/logout

### Requirements

REST-Full api for CRUD

>note: must use gin

Use of PostgreSQL for data persistency

At least 70% test coverage

>note: usage of mocking

Multi-stage docker and docker compose build files

Usage of OpenAPI swagger for documenting endpoints

Write a readme explaining the project structure and possible trade-offs.

Basic observability

- Prometheus Metrics
- tracing

#### Optional Requirements

Cache layer using redis for /GET endpoints

>note: Add cache invalidation after update or delete occurred

Pagination and filtering  (status, assignee)

Usage of grafana k6 to do benchmark and load test

### Project Structure

```
    RestApi Port
        |
    Routes
        |
    Handlers
        |
    Controller
        |
    Repository
        |
    Data Models
```

#### Setup

You can project with one one of following ways.

#### Air ( auto reload )

`
    #if air is not in your PATH.
    go install github.com/air-verse/air@latest
    air
`

#### Docker Compose Setup ( Server )

`
    docker compose up -d --build
    docker compose logs -f app postgres redis
`


#### Testing Project

Inside project there are test files name in idiomatic golang style ` <pkg_name>_test.go`.
in order to run tests and check coverage run following command :


`
    go test ./internal/... -coverprofile=coverage.out -covermode=atomic
`

`
#Sample response
        graph-interview/internal/api            coverage: 0.0% of statements
ok      graph-interview/internal/api/handlers   1.237s  coverage: 84.8% of statements
ok      graph-interview/internal/api/handlers/dto       0.485s  coverage: 100.0% of statements
        graph-interview/internal/api/handlers/errors            coverage: 0.0% of statements
ok      graph-interview/internal/api/middlewares        0.553s  coverage: 38.3% of statements
        graph-interview/internal/cfg            coverage: 0.0% of statements
?       graph-interview/internal/domain [no test files]
?       graph-interview/internal/repository     [no test files]
ok      graph-interview/internal/repository/cache       2.301s  coverage: 16.7% of statements
ok      graph-interview/internal/repository/enum        2.134s  coverage: 100.0% of statements
        graph-interview/internal/repository/mock                coverage: 0.0% of statements
        graph-interview/internal/repository/storage             coverage: 0.0% of statements
        graph-interview/internal/repository/storage/postgres            coverage: 0.0% of statements
ok      graph-interview/internal/services       1.107s  coverage: 73.7% of statements
`