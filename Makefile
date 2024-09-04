VERSION=-ldflags="-X main.Version=$(shell git describe --tags)"

build-docker: ## Builds a docker image with the binary
	docker build -t miniapp/backend -f ./Dockerfile .

run-docker:
	docker run -d --name backend -p 10443:10443 -v ./config.ini:/app/config.ini    miniapp/backend
