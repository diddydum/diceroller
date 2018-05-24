MIGRATE := /usr/local/bin/migrate

$(MIGRATE):
	go get -u -d github.com/golang-migrate/migrate/cli github.com/lib/pq
	go build -tags 'postgres' -o /usr/local/bin/migrate github.com/golang-migrate/migrate/cli

.PHONY: create-docker-db
create-docker-db:
	docker run --name diceroller-db -p 5432:5432 --mount type=bind,source="$(CURDIR)"/docker/db/docker-entrypoint-initdb.d,target=/docker-entrypoint-initdb.d -d postgres:10

.PHONY: clean-docker-containers
clean-docker-containers:
	docker stop diceroller-db && docker rm diceroller-db

migrate: $(MIGRATE)
	$(MIGRATE) -path $(CURDIR)/migrations -database "postgres://diceroller:diceroller@localhost/diceroller?sslmode=disable" up

.PHONY: clean
clean:
	rm diceroller

.PHONY: build
build:
	go build -v -o diceroller
	go test -v
	go vet
