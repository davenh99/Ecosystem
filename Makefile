build:
	@go build -o bin/david-erp-backend cmd/main.go

test:
	@go test -v ./...

run: build
	@./bin/david-erp-backend
	# @npm run dev --prefix ui/

migration:
	@migrate create -ext sql -dir cmd/migrate/migrations $(filter-out $@, $(MAKECMDGOALS))

# change the below two to use the ./backend migrate or something?
migrate-up:
	@go run cmd/migrate/main.go up

migrate-down:
	@go run cmd/migrate/main.go down

# DROP DATABASE dave_erp;CREATE DATABASE dave_erp;GRANT ALL PRIVILEGES ON dave_erp.* TO 'david'@'localhost';FLUSH PRIVILEGES;USE dave_erp;
