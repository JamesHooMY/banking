run:
	go run main.go apiserver

lint:
	golangci-lint run --timeout 10m ./... --fix

tidy:
	go mod tidy && go mod vendor

test:
	go clean -testcache && go test ./...

cover:
	go clean -testcache && go test ./... -coverprofile=coverage.out && go tool cover -html=coverage.out

swagger:
	swag init

docker_up:
	docker-compose -f ./build/docker-compose.yml up -d

docker_down:
	docker-compose -f ./build/docker-compose.yml down
