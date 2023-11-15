package main

import (
	pb "central/github.com/Sistemas-Distribuidos-2023-02/Grupo10-Laboratorio-3"
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

type vanguardiaServiceServer struct {
	pb.UnimplementedMiServicioServer
	brokerClient pb.MiServicioClient
}

func NewVanguardiaServiceServer(brokerClient pb.MiServicioClient) *vanguardiaServiceServer {
	return &vanguardiaServiceServer{
		brokerClient: brokerClient,
	}
}

func (s *vanguardiaServiceServer) GetSoldados(ctx context.Context, req *pb.GetSoldadosRequest) (*pb.GetSoldadosResponse, error) {
	//Implementar petición a broker y que este le pida a fulcrums los soldados del sector
	return nil, nil
}

func main() {
	// Configurar la conexión con el servidor Broker
	connBroker, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Error al conectar al servidor Broker: %v", err)
	}
	defer connBroker.Close()

	// Crear el cliente para comunicarse con el servidor Broker
	brokerClient := pb.NewMiServicioClient(connBroker)

	// Configurar el servidor de Vanguardia
	listener, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Error al escuchar: %v", err)
	}

	server := grpc.NewServer()
	vanguardiaServer := NewVanguardiaServiceServer(brokerClient)
	pb.RegisterMiServicioServer(server, vanguardiaServer)

	fmt.Println("Servidor de Vanguardia gRPC iniciado en el puerto 50052")

	if err := server.Serve(listener); err != nil {
		log.Fatalf("Error al servir: %v", err)
	}
}
