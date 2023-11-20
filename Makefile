# Define variables
BROKER_DOCKER_IMAGE = broker-server
FULCRUM1_DOCKER_IMAGE = fulcrum1-server
FULCRUM2_DOCKER_IMAGE = fulcrum2-server
FULCRUM3_DOCKER_IMAGE = fulcrum3-server
CAIATL_DOCKER_IMAGE = caiatl-server
OSIRIS_DOCKER_IMAGE = osiris-server
VANGUARDIA_DOCKER_IMAGE = vanguardia-server

BROKER_PORT= 50051
FULCRUM1_PORT= 50052
FULCRUM2_PORT= 50053
FULCRUM3_PORT= 50054

# Define the default target
all: help

# Build the ONU server Docker image
docker-vanguardia:
	docker build -t $(VANGUARDIA_DOCKER_IMAGE) --build-arg SERVER_TYPE=vanguardia .
	docker run -it --name $(VANGUARDIA_DOCKER_IMAGE) -e SERVER_TYPE=vanguardia $(VANGUARDIA_DOCKER_IMAGE)

# Build the OMS server Docker image
docker-broker:
	docker build -t $(BROKER_DOCKER_IMAGE) --build-arg SERVER_TYPE=broker .
	docker run -d --name $(BROKER_DOCKER_IMAGE) -e SERVER_TYPE=broker -p $(BROKER_PORT):$(BROKER_PORT) $(BROKER_DOCKER_IMAGE)

# Build the FULCRUMS server Docker images
docker-fulcrums:
	@echo "SERVER_TYPE is set to: $(SERVER_TYPE)"
	@if [ "$(SERVER_TYPE)" = "fulcrum1" ]; then \
		docker build -t $(FULCRUM1_DOCKER_IMAGE) --build-arg SERVER_TYPE=fulcrum1 .; \
		docker run -d --name $(FULCRUM1_DOCKER_IMAGE) -e SERVER_TYPE=$(SERVER_TYPE) -p $(FULCRUM1_PORT):$(FULCRUM1_PORT) $(FULCRUM1_DOCKER_IMAGE); \
	elif [ "$(SERVER_TYPE)" = "fulcrum2" ]; then \
		docker build -t $(FULCRUM2_DOCKER_IMAGE) --build-arg SERVER_TYPE=fulcrum2 .; \
		docker run -d --name $(FULCRUM2_DOCKER_IMAGE) -e SERVER_TYPE=$(SERVER_TYPE) -p $(FULCRUM2_PORT):$(FULCRUM2_PORT) $(FULCRUM2_DOCKER_IMAGE); \
	elif [ "$(SERVER_TYPE)" = "fulcrum3" ]; then \
		docker build -t $(FULCRUM3_DOCKER_IMAGE) --build-arg SERVER_TYPE=fulcrum3 .; \
		docker run -d --name $(FULCRUM3_DOCKER_IMAGE) -e SERVER_TYPE=$(SERVER_TYPE) -p $(FULCRUM3_PORT):$(FULCRUM3_PORT) $(FULCRUM3_DOCKER_IMAGE); \
	else \
		echo "Invalid SERVER_TYPE argument. Use 'fulcrum1', 'fulcrum2', or 'fulcrum3'."; \
		exit 1; \
	fi

# Build the CONTINENTES server Docker images
docker-informantes:
	@echo "SERVER_TYPE is set to: $(SERVER_TYPE)"
	@if [ "$(SERVER_TYPE)" = "caiatl" ]; then \
		docker build -t $(CAIATL_DOCKER_IMAGE) --build-arg SERVER_TYPE=caiatl .; \
		docker run -it --name $(CAIATL_DOCKER_IMAGE) -e SERVER_TYPE=$(SERVER_TYPE) $(CAIATL_DOCKER_IMAGE); \
	elif [ "$(SERVER_TYPE)" = "osiris" ]; then \
		docker build -t $(OSIRIS_DOCKER_IMAGE) --build-arg SERVER_TYPE=osiris .; \
		docker run -it --name $(OSIRIS_DOCKER_IMAGE) -e SERVER_TYPE=$(SERVER_TYPE) $(OSIRIS_DOCKER_IMAGE); \
	else \
		echo "Invalid SERVER_TYPE argument. Use 'caiatl',or 'osiris'."; \
		exit 1; \
	fi


# Usage: make help
help:
	@echo "Available targets:"
	@echo "  docker-vanguardia   - Iniciar el codigo Docker para el servidor vanguardia"
	@echo "  docker-broker   - Iniciar el codigo Docker para el servidor broker"
	@echo "  docker-fulcrums  SERVER_TYPE={fulcrum1,fulcrum2,fulcrum3}  - Iniciar el codigo Docker para el servidor fulcrum especificado"
	@echo "  docker-informantes  SERVER_TYPE={caiatl,osiris}  - Iniciar el codigo Docker para el servidor informante especificado"
	@echo "  help             - Pide ayuda"