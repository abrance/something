PROJECT_NAME=something

.PHONY: docker-build
docker-build:
	@echo "Building Docker image..."
	@docker build -t $(PROJECT_NAME) .

.PHONY: docker-run
docker-run:
	@echo "Running Docker image..."
	@docker run -d --name $(PROJECT_NAME) $(PROJECT_NAME)

.PHONY: docker-rebuild
docker-rebuild:
	@echo "Rebuilding Docker image..."
	@docker stop $(PROJECT_NAME)
	@docker build -t $(PROJECT_NAME) .
