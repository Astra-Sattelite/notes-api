.PHONY: build run test check fmt tidy update update-patch swagger
docker-build docker-up docker-down docker-logs docker-clean

Go
build: go build ./...

run: go run .

test: go test ./...

check: go vet ./... go test ./... go build ./...

fmt: go fmt ./...

tidy: go mod tidy

swagger: swag init

update: go get -u -t ./... go mod tidy go test ./... go build ./...

update-patch: go get -u=patch -t ./... go mod tidy go test ./... go build ./...

Docker
docker-build: docker compose build

docker-up: docker compose up -d

docker-down: docker compose down

docker-logs: docker compose logs -f

WARNING: Removes Docker volumes, including persistent database data.
docker-clean: docker compose down -v