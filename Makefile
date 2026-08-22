include .env
export

migrate-up:
	migrate -path migrations -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path migrations -database "$(DATABASE_URL)" down 1

migrate-create:
	migrate create -ext sql -dir migrations -seq $(name)

swagger:m 
	swag init -g cmd/api/main.go -o docs

migrate-test-up:
	migrate -path migrations -database "$(TEST_DATABASE_URL)" up

migrate-test-down:
	migrate -path migrations -database "$(TEST_DATABASE_URL)" down 1
	
test:
	go test ./...

test-integration:
	go test -tags=integration ./internal/repository/postgres/... -v

test-cover:
	go test ./... -coverprofile=coverage.out
	go tool cover -func=coverage.out

test-cover-html:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

docker-up:
	docker compose up --build

docker-down:
	docker compose down

docker-down-clean:
	docker compose down -v