# Virgo4 Collections Context Service

This is a web service provides collections context servives to the Virgo4 Client

Requires Go 1.16.0+

### Current API

* GET /version - get version information
* GET /lookup?q=[target name] - lookup collection context for a named collection. Ex: Daily Progress Digitized Microfilm
* GET /collections/:id/dates?year=YYYY - get publication dates for a collection year
* GET /collections/:id/items/:date/next - get the next published item; date format=yyyy-mm-dd
* GET /collections/:id/items/:date/previous - get the previous published item; date format=yyyy-mm-dd

### Database Notes

This service uses a Postgres DB to track collection data. The schema is managed by
https://github.com/golang-migrate/migrate and the scripts are in ./db/migrations.

Install the migrate binary on your host system. For OSX, the easiest method is brew. Execute:

`brew install golang-migrate`.

Define your PSQL connection params in an environment variable, like this:

`export COLL_DB=postgres://v4_collections:pass@localhost:5432/v4_collections?sslmode=disable`

Run migrations like this:

`migrate -database ${COLL_DB} -path db/migrations up`

Example migrate commads to create a migration and run one:

* `migrate create -ext sql -dir db/migrations -seq update_user_auth`
* `migrate -database ${COLL_DB} -path db/migrations/ up`

Sample setup of a user for local testing:

`postgres=# create database v4_collections;`
`postgres=# create user v4_collections with encrypted password 'pass';`
`postgres=# grant all privileges on database v4_collections to v4_collections;`
