run:
	gofmt -w .
	go run main.go
up:
	docker-compose up
down:
	docker-compose down
build:
	docker-compose up --build
docker:
	docker build -t bambamload .
run_docker:
	docker run bambamload
lint:
	cd $(shell dirname $(realpath go.mod)) && golangci-lint run ./... --timeout=2m -D staticcheck,govet

clean:
	docker-compose down --volumes --remove-orphans
	docker container prune -f
	docker volume prune -f
	docker network prune -f
	docker system prune -a --volumes -f