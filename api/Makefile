include .env

NAME = init

create_migration:
	migrate create -ext=sql -dir=migrations -seq ${NAME}

migrate_up:
	migrate -path=migrations -database \
	    "postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable" \
		-verbose up

migrate_down:
	migrate -path=migrations -database \
	    "postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable" \
		-verbose down

migrate_force:
	@[ -n "$(LEVEL)" ] || (echo "ERROR: LEVEL is required. Usage: make migrate_force LEVEL=<value>"; exit 1)
	migrate -path=migrations -database \
	    "postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable" \
		force ${LEVEL}

.PHONY: create_migration migrate_up migrate_down
