DB_URL=postgres://postgres:password@localhost:5432/go_k8s?sslmode=disable
TEST_DB_URL=postgres://postgres:password@localhost:5432/test_go_k8s?sslmode=disable

migrate_create:
	migrate create -ext sql -dir db/migrations -seq ${MIGRATE_NAME}
migrate_up:
	migrate -path db/migrations -database "${TEST_DB_URL}" -verbose up
migrate_down:
	migrate -path db/migrations -database "${TEST_DB_URL}" -verbose down
migrate_fix:
	migrate -path db/migrations -database "${TEST_DB_URL}" force ${CLEAN_VERSION}
sqlc:
	sqlc generate --file internal/db/sqlc.yaml
tests:
	go test -v -cover ./...
mocks:
	mockgen -package mockdb --destination pkg/mocks/store.go go-k8s/internal/sqlc TxStore
start:
	go run ./cmd/api/main.go

.PHONY: migrate_create migrate_up migrate_down sqlc tests mocks start