VERSION=-ldflags="-X main.Version=$(shell git describe --tags)"

build-docker: ## Builds a docker image with the binary
	docker build -t miniapp/backend -f ./Dockerfile .

docker-run:
	docker run -d --name backend -v ./config.ini:/app/config.ini    miniapp/backend
