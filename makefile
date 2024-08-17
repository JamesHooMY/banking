run:
	go run main.go apiserver

build:
	go build -o apiserver main.go

lint:
	golangci-lint run --timeout 10m ./... --fix

update:
	go get -u ./...

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

pprof_report:
	curl -o ./pprof/cpu.pprof http://localhost:6060/debug/pprof/profile?seconds=30

pprof_web:
	go tool pprof -http=:8081 ./pprof/cpu.pprof

pprof_mem:
	curl -o ./pprof/mem.pprof http://localhost:6060/debug/pprof/heap?seconds=30

pprof_flame:
	go tool pprof -svg ./pprof/cpu.pprof > ./pprof/cpu.svg
	go tool pprof -svg ./pprof/mem.pprof > ./pprof/mem.svg

pprof_top:
	go tool pprof -top ./pprof/cpu.pprof | head -n 20 > ./pprof/top_cpu.txt
	go tool pprof -top ./pprof/mem.pprof | head -n 20 > ./pprof/top_mem.txt

wire:
	wire ./...