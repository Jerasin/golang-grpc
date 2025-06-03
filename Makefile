run-all:
	cd api-gateway && make server &
	cd auth-svc && make server 
	# cd order-svc && make server &
	# cd product-svc && make server &


dev: destroy build-package build up logs
prod: destroy build-prod up logs

run: up logs

build-package:
	cd api-gateway && go mod tidy &
	cd auth-svc && go mod tidy 
	
build:
	@COMPOSE_DOCKER_CLI_BUILD=1 DOCKER_BUILDKIT=1 docker-compose build $(c)
destroy:
	@docker-compose down -v $(c)
up:
	@docker-compose up -d $(c)
up-db:
	@docker-compose -f docker-compose.yml up -d db $(c)
logs:
	@docker-compose logs --tail=100 -f $(c)
down:
	@docker-compose down $(c)

build-prod:
	@COMPOSE_DOCKER_CLI_BUILD=1 DOCKER_BUILDKIT=1 docker-compose -f docker-compose.prod.yml build $(c)

set-env-test:
	export RUN_TEST_MODE=dev

go-test:
	RUN_TEST_MODE=prod go test -v -count=1 ./... 

go-test-debug:
	go test -v -count=1 ./... 

test:
	@make generate-mock
	@make go-test

generate-mock:
	mockery --name=BaseServiceInterface --dir=./app/service  --output=./app/mocks --filename=base_service_mock.go
	mockery --name=BaseRepositoryInterface --dir=./app/repository  --output=./app/mocks --filename=base_repository_mock.go

test-debug:
	@make go-test-debug