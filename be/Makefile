run:
	go run cmd/main.go
run-air:
	air -c .air.toml
sqlc:
	sqlc generate -f infrastructure/db/sqlc/sqlc.yaml
newmg:
	migrate create -ext sql -dir infrastructure/db/migrations $(name)
upmg:
	@bash scripts/migrate-up.sh
downmg:
	@bash scripts/migrate-down.sh

.PHONY: run sqlc upmg downmg