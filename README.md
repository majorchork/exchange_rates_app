# rates_app

## About the Project.

#### This project is to show historical exchange rates using Clean Architecture and use only the standard library other than database connection and mocks for testing.
#### The data was gotten from [European Central Bank](https://www.ecb.europa.eu/stats/eurofxref/eurofxref-hist-90d.xml)

### this app requires

- Golang: go version go1.18 
- Go SQLite3
- Optional: Docker

### usage
#### sync all dependencies with the command 
 ```
 $ go mod tidy
 ```
#### run the backend with the command 
```
$ make run
```
#### Or, run using docker
```make docker
docker run --name ratesapp -p 8085:8085 --rm air 
```
#### test
```
$make mock

$make test
````
### List of Endpoints
- Get forex latest rates: http://localhost:8085/rates/latest

- Get forex rates at certain dates (format: yyyy-mm-dd): http://localhost:8085/rates/2022-06-24

- Analyze forex rates: http://localhost:8085/rates/analyze
