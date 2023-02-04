# 도커 관련된 변수와 세팅입니다.
SERVER_NAME = jgc-api
SERVER_VERSION = 0.1
CONTAINER_NAME = JGC-API
PORT = 8806

# 콘솔 색 관련 세팅입니다.
# 아래와 같이 사용하면 됩니다.
# @echo "$(GRENN)Hello World$(RESET)"
GREEN = \033[92m
RESET = \033[0m

# 명령어에 붙는 prefix입니다.
PREFIX = $(GREEN)[JGC]$(RESET)

server:
	@echo "$(PREFIX) Building api server image..."
	@docker build \
		--platform linux/x86_64 \
		-t $(SERVER_NAME):$(SERVER_VERSION) .
	@echo "$(PREFIX) Done building api server image."
.PHONY: server

run-product:
	@echo "$(PREFIX) Running api server image..."
	@docker run \
		--platform linux/x86_64 \
		--name $(CONTAINER_NAME) \
		--network $(DOCKER_NETWORK) \
		--restart always \
		-it -d -p $(PORT):$(PORT) \
		$(SERVER_NAME):$(SERVER_VERSION) \
			-env product
	@echo "$(PREFIX) Success running api server image."
.PHONY: run-product

run-test:
	@echo "$(PREFIX) Running api server image..."
	@docker run \
		--platform linux/x86_64 \
		--name $(CONTAINER_NAME) \
		--network $(DOCKER_NETWORK) \
		--restart always \
		-it -d -p $(PORT):$(PORT) \
		$(SERVER_NAME):$(SERVER_VERSION) \
			-env test
	@echo "$(PREFIX) Success running api server image."
.PHONY: run-test-product

migrate-product:
	@echo "$(PREFIX) Migrate Product DB..."
	@docker run \
		--platform linux/x86_64 \
		--network $(DOCKER_NETWORK) \
		--rm -it -d \
		$(SERVER_NAME):$(SERVER_VERSION) \
			-env product \
			-migrate
.PHONY: migrate-product

migrate-test:
	@echo "$(PREFIX) Migrate Test DB..."
	@docker run \
		--platform linux/x86_64 \
		--network $(DOCKER_NETWORK) \
		--rm -it -d \
		$(SERVER_NAME):$(SERVER_VERSION) \
			-env test \
			-migrate
.PHONY: migrate-test

migrate-local:
	@echo "$(PREFIX) Migrate Native DB..."
	@go run main.go \
		-env native \
		-migrate
.PHONY: migrate-local

serve:
	@echo "$(PREFIX) Running api server..."
	@go run main.go \
		-env native
	@echo "$(PREFIX) Success running api server."
.PHONY: serve

stop:
	@echo "$(PREFIX) Stopping api server..."
	@docker stop \
		$(CONTAINER_NAME)
	@echo "$(PREFIX) Success stopping api server."
.PHONY: stop

DANGLING_IMAGE = $(shell docker images -f dangling=true -q)
API_IMAGE = $(shell docker images --filter=reference="jgc-api" -q)
clean:
	@echo "$(PREFIX) Removing dangling images..."
ifneq ($(DANGLING_IMAGE),)
	@docker rmi $(DANGLING_IMAGE)
endif
	@echo "$(PREFIX) Done removing dangling images."

	@echo "$(PREFIX) Removing all jgc-api images..."
ifneq ($(API_IMAGE),)
	@docker rmi -f $(API_IMAGE)
endif
	@echo "$(PREFIX) Done removing all jgc-api images."
.PHONY: clean
