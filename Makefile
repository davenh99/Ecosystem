build:
	@go build -o bin/ecosystem cmd/main.go

test:
	@go test -v ./...

run: build
	@./bin/ecosystem
	# @npm run dev --prefix ui/

# change the below two to use the ./backend migrate or something?
migrate-up:
	@go run cmd/migrate/main.go up

migrate-down:
	@go run cmd/migrate/main.go down

# DROP DATABASE dave_erp;CREATE DATABASE dave_erp;GRANT ALL PRIVILEGES ON dave_erp.* TO 'david'@'localhost';FLUSH PRIVILEGES;USE dave_erp;
