# Define the binaries for the front end and broker
FRONT_END_BINARY=frontApp
BROKER_BINARY=brokerApp


## start-containers: starts all containers in the background without forcing build
start-containers:
	@echo "Starting Docker images..."
	docker-compose up -d
	@echo "Docker images started!"


## stop-containers: stop docker compose
stop-containers:
	@echo "Stopping docker compose..."
	docker-compose down
	@echo "Done!"


## rebuild-containers: stops docker-compose (if running), builds all projects and starts docker compose
rebuild-containers: build-broker
	@echo "Stopping docker images (if running...)"
	docker-compose down
	@echo "Building (when required) and starting docker images..."
	docker-compose up --build -d
	@echo "Docker images built and started!"




## build-broker: builds the broker binary as a linux executable
build-broker:
	@echo "Building broker binary..."
	cd ../broker-service && env GOOS=linux CGO_ENABLED=0 go build -o ${BROKER_BINARY} ./service
	@echo "Done!"




## build-front: builds the front end binary
build-front:
	@echo "Building front end binary..."
	cd ../front-end && env CGO_ENABLED=0 go build -o ${FRONT_END_BINARY} .
	@echo "Done!"




## start-front: starts the front end
start-front: build-front
	@echo "Starting front end"
	cd ../front-end && ./${FRONT_END_BINARY} 




## stop-front: stop the front end
stop-front:
	@echo "Stopping front end..."
	@-pkill -SIGTERM -f "./${FRONT_END_BINARY}"
	@echo "Stopped front end!"


.PHONY: clean-docker-images
clean-docker-images:
	@echo "Removing all Docker images..."
	docker rmi $(shell docker images -q)
	@echo "All Docker images removed!"
