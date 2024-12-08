sqlc:
	-rm schema/*.sql
	-rm sqlcgen/*.sql.go
	cp ../sql/1-schema.sql schema/schema.sql
	sqlc generate

install: ## Install sqlc
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.27.0

.PHONY: sqlc install
