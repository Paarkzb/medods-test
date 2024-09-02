migrate-up:
	goose -dir .\migrations postgres "host=localhost port=5443 user=postgres password=postgres dbname=postgres sslmode=disable" up
migrate-down:
	goose -dir .\migrations postgres "host=localhost port=5443 user=postgres password=postgres dbname=postgres sslmode=disable" down
build:
	docker-compose build
start:
	docker-compose up -d --build
stop:
	docker-compose stop
down:
	docker-compose down
test:
	go test -v ./...