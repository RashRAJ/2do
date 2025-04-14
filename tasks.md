# Tasks

## Essentials
- Docker-compose file to spin up all the services needed by this application
  - Database
  - Service
- Implement the database connection
- Binding repository layer to service layer using interface and dependency injection

## Optional
- Centralise config for ease of update and change
  - Read config in from yaml/json/toml file
  - Put it in configmap and mount volume path - k8s