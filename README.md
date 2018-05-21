# Diceroller

## Usage
Diceroller requires a postgres db running. Easiest way to launch one is via docker:

```
$ docker run --name diceroller-db -p 5432:5432 --mount type=bind,source="$(pwd)"/docker/db/docker-entrypoint-initdb.d,target=/docker-entrypoint-initdb.d -d postgres:10
```

This launches a db instance bound to port 5432. You can start and stop it via `docker start diceroller-db` and `docker stop diceroller-db` respectively.
