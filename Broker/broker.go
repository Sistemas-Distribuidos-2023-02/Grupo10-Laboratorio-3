package main

import (
	pb "central/github.com/Sistemas-Distribuidos-2023-02/Grupo10-Laboratorio-3" // Asegúrate de ajustar la importación correctamente
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

type baseServiceServer struct {
	pb.UnimplementedMiServicioServer
}

func (s *baseServiceServer) AgregarBase(ctx context.Context, req *pb.AgregarBaseRequest) (*pb.Respuesta, error) {
	// Lógica para el comando AgregarBase
	// ...

	return &pb.Respuesta{Mensaje: "Comando AgregarBase ejecutado", Exitoso: true}, nil
}

// Implementa los demás métodos del servicio de manera similar

func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Error al escuchar: %v", err)
	}

	server := grpc.NewServer()
	pb.RegisterMiServicioServer(server, &baseServiceServer{})

	fmt.Println("Servidor gRPC iniciado en el puerto 50051")

	if err := server.Serve(listener); err != nil {
		log.Fatalf("Error al servir: %v", err)
	}
}