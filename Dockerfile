# Use the specified base image
ARG BASE_IMAGE=golang:1.21.3
FROM ${BASE_IMAGE} AS builder

ARG BROKER_PORT= 50051
ARG FULCRUM1_PORT= 50052
ARG FULCRUM2_PORT= 50053
ARG FULCRUM3_PORT= 50054

ARG SERVER_TYPE

# Set the working directory inside the container
WORKDIR /app

# Copy the parent directory's go.mod and go.sum files to the container
COPY go.mod .
COPY go.sum .

# Download and install Go dependencies
RUN go mod download

# Copy the rest of your application code to the container
COPY . .

CMD if [ "$SERVER_TYPE" = "vanguardia" ]; then \
        cd /app/Vanguardia; \
        go build -o vanguardia-server; \
        ./vanguardia-server; \
    elif [ "$SERVER_TYPE" = "broker" ]; then \
        PORT=$BROKER_PORT; \
        cd /app/Broker; \
        go build -o broker-server; \
        ./broker-server; \
    elif [ "$SERVER_TYPE" = "f1" ]; then \
        PORT=$FULCRUM1_PORT; \
        cd /app/Fulcrum1; \
        go build -o f1-server; \
        ./f1-server; \
    elif [ "$SERVER_TYPE" = "f2" ]; then \
        PORT=$FULCRUM2_PORT; \
        cd /app/Fulcrum2; \
        go build -o f2-server; \
        ./f2-server; \
    elif [ "$SERVER_TYPE" = "f3" ]; then \
        PORT=$FULCRUM3_PORT; \
        cd /app/Fulcrum3; \
        go build -o f3-server; \
        ./f3-server; \
    elif [ "$SERVER_TYPE" = "i1" ]; then \
        cd /app/Informante1; \
        go build -o i1-server; \
        ./i1-server; \
    elif [ "$SERVER_TYPE" = "i2" ]; then \
        cd /app/Informante2; \
        go build -o i2-server; \
        ./i2-server; \
    else \
        echo "Invalid SERVER_TYPE argument."; \
    fi