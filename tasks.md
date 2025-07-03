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
  
## immediate  Next steps
- Optimize Data base Schema, Table collumns and best way to run db connections.
- logging [Done]
- use [go migrate](https://github.com/golang-migrate/migrate) to run db migration
- 



## Next Steps

**Implement Authentication**: Replace the placeholder Auth middleware with actual authentication
**Implement Caching**: Add caching for frequently accessed data
