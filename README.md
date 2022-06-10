# Go Clean Architecture Service

## API
Please make sure you create ```.env``` file first, take a look at ```.env.example``` for the layout of environment variables.
#### Run 
```sh
make run
```
#### Build
```sh
make build
```

## Migration
This service uses [database migration](https://en.wikipedia.org/wiki/Schema_migration) to manage the changes of the 
database schema over the whole project development phase. The following commands are commonly used with regard to database schema changes:
#### Up
```sh
make migrate-up
```
#### Down
```sh
make migrate-down
```
#### Fresh
Drop All Tables & Migrate
or 
```sh
make migrate-fresh
```
#### Create a New Migration File
```sh
make migrate-new
```

## Scheduler
There are 1 scheduler for this service:
- cleanup

To run a scheduler, use the command below:
```sh
./application cron <scheduler_type>
```

## API Docs
Please go to ```<BASE_URL>/swagger/index.html``` for the API docs

## Common Development Tasks

This section describes some common development tasks using this service.

#### Implementing a New Feature

Implementing a new feature typically involves the following steps:

1. Develop the service that implements the business logic supporting the feature. Please refer to `internal/<module>/service.go` as an example.
2. Develop the RESTful API exposing the service about the feature. Please refer to `internal/<module>/api.go` as an example.
3. Develop the repository that persists the data entities needed by the service. Please refer to `internal/repository/<module>.go` as an example.
4. Wire up the above components together by injecting their dependencies in the main function. Please refer to the `<module>.RegisterAPI()` call in `app/api/api.go`.

## Deployment
The application can be run as a docker container. You can use ```make build-docker``` to build the application into a docker image. The docker container starts with the ```./application```. Later you can pass the docker args to run the spesific command.
