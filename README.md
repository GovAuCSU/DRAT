### Dependency Risk Analysis Tool - DRAT 

DRAT aims to provide risk indicator for libraries used by the developer within an organisation.

This tool's intial code and design is inspired by the certificate transparency work done by Adam Eijdenberg from Cloud.gov.au team. You can find the repository for that [here](https://github.com/govau/certwatch/tree/master/jobs)

## Running locally

```bash
docker run -p 5432:5432 --name dratpg -e POSTGRES_USER=dratpg -e POSTGRES_PASSWORD=dratpg -d postgres
export VCAP_SERVICES='{"postgres": [{"credentials": {"username": "certpg", "host": "localhost", "password": "certpg", "name": "certpg", "port": 5434}, "tags": ["postgres"]}]}'
go run *.go
```

## Running in docker with latest code

Note: We are running dratpg seperately from docker-compose.yml so the db docker persist while we update the other containers.

```bash
docker run -p 5432:5432 --name dratpg -e POSTGRES_USER=dratpg -e POSTGRES_PASSWORD=dratpg -d 
GOOS=linux GOARCH=amd64 go build -o bin/drat cmd/drat/main.go
docker-compose up
```

To checkout database:

```bash
psql "dbname=dratpg host=localhost user=dratpg password=dratpg port=5432"
```


### Postgresql and queueing tricks

To drop all tables and start fresh, we could run the following instead of re-creating the docker container
```
DROP SCHEMA public CASCADE;CREATE SCHEMA public;
```

To re-run organisation crawl job (which run daily), we can delete the metadata for cronjob and set `run_at` to `now()` for that job.

To start all jobs again:
```
update que_jobs set run_at = now();
```
