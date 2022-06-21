deps:
	go get ./...

up:
	docker-compose up -d
	docker-compose logs -f

down:
	docker-compose down

test:
	go test -v ./...

test-ci: down
	docker-compose -f ./docker-compose.test.yaml build
	docker-compose -f ./docker-compose.test.yaml up -d && go test -v --tags=integration ./...
	docker-compose -f ./docker-compose.test.yaml down

build:
	docker build --build-arg ENV=prod -t picture-room-api .

deploy:
	gcloud app deploy

migrate-redo:
	goose -dir internal/db/migrations postgres "user=postgres password=postgres dbname=picture-room sslmode=disable" redo

migrate-up:
	goose -dir internal/db/migrations postgres "user=postgres password=postgres dbname=picture-room sslmode=disable" up

migrate-down:
	goose -dir internal/db/migrations postgres "user=postgres password=postgres dbname=picture-room sslmode=disable" down

sqlc:
	sqlc generate

.PHONY: deps up test test-ci build deploy down migrate-up migrate-down migrate-redo sqlc
