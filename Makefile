# Definir variables
BROKER_DOCKER_IMAGE = broker-server
FULCRUM1_DOCKER_IMAGE = f1-server
FULCRUM2_DOCKER_IMAGE = f2-server
FULCRUM3_DOCKER_IMAGE = f3-server
CAIATL_DOCKER_IMAGE = i1-server
OSIRIS_DOCKER_IMAGE = i2-server
VANGUARDIA_DOCKER_IMAGE = vanguardia-server

BROKER_PORT= 50051
FULCRUM1_PORT= 50052
FULCRUM2_PORT= 50053
FULCRUM3_PORT= 50054

# Define the default target
all: help

# Construir la imagen de la vanguardia
docker-vanguardia:
	docker build -t $(VANGUARDIA_DOCKER_IMAGE) --build-arg SERVER_TYPE=vanguardia .
	docker run -it --name $(VANGUARDIA_DOCKER_IMAGE) -e SERVER_TYPE=vanguardia $(VANGUARDIA_DOCKER_IMAGE)

# Construir la imagen del broker
docker-broker:
	docker build -t $(BROKER_DOCKER_IMAGE) --build-arg SERVER_TYPE=broker .
	docker run -d --name $(BROKER_DOCKER_IMAGE) -e SERVER_TYPE=broker -p $(BROKER_PORT):$(BROKER_PORT) $(BROKER_DOCKER_IMAGE)

# Construir la imagen de los fulcrum
docker-f1:
	docker build -t $(FULCRUM1_DOCKER_IMAGE) --build-arg SERVER_TYPE=f1 .
	docker run -d --name $(FULCRUM1_DOCKER_IMAGE) -e SERVER_TYPE=$(SERVER_TYPE) -p $(FULCRUM1_PORT):$(FULCRUM1_PORT) $(FULCRUM1_DOCKER_IMAGE)

docker-f2:
	docker build -t $(FULCRUM2_DOCKER_IMAGE) --build-arg SERVER_TYPE=f2 .
	docker run -d --name $(FULCRUM2_DOCKER_IMAGE) -e SERVER_TYPE=$(SERVER_TYPE) -p $(FULCRUM2_PORT):$(FULCRUM2_PORT) $(FULCRUM2_DOCKER_IMAGE)

docker-f3:
	docker build -t $(FULCRUM3_DOCKER_IMAGE) --build-arg SERVER_TYPE=f3 .
	docker run -d --name $(FULCRUM3_DOCKER_IMAGE) -e SERVER_TYPE=$(SERVER_TYPE) -p $(FULCRUM3_PORT):$(FULCRUM3_PORT) $(FULCRUM3_DOCKER_IMAGE)

# Construir la imagen de los informantes
docker-i1:
	docker build -t $(CAIATL_DOCKER_IMAGE) --build-arg SERVER_TYPE=i1 .
	docker run -it --name $(CAIATL_DOCKER_IMAGE) -e SERVER_TYPE=$(SERVER_TYPE) $(CAIATL_DOCKER_IMAGE)

docker-i2:
	docker build -t $(OSIRIS_DOCKER_IMAGE) --build-arg SERVER_TYPE=i2 .
	docker run -it --name $(OSIRIS_DOCKER_IMAGE) -e SERVER_TYPE=$(SERVER_TYPE) $(OSIRIS_DOCKER_IMAGE)


# Usage: make help
help:
	@echo "Available targets:"
	@echo "  docker-vanguardia   - Iniciar el codigo Docker para el servidor vanguardia"
	@echo "  docker-broker   - Iniciar el codigo Docker para el servidor broker"
	@echo "  docker-f1   - Iniciar el codigo Docker para el servidor fulcrum1"
	@echo "  docker-f2   - Iniciar el codigo Docker para el servidor fulcrum2"
	@echo "  docker-f3   - Iniciar el codigo Docker para el servidor fulcrum3"
	@echo "  docker-i1  - Iniciar el codigo Docker para el servidor caiatl"
	@echo "  docker-i2  - Iniciar el codigo Docker para el servidor osiris"
	@echo "  help             - Pide ayuda"