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
