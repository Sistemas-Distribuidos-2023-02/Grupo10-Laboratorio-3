package main

import (
	pb "central/github.com/Sistemas-Distribuidos-2023-02/Grupo10-Laboratorio-3" // Asegúrate de ajustar la importación correctamente
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"

	"google.golang.org/grpc"
)

type baseServiceServer struct {
	pb.UnimplementedMiServicioServer
}

func (s *baseServiceServer) AgregarBase(ctx context.Context, req *pb.AgregarBaseRequest) (*pb.Respuesta, error) {
	puerto := s.RandomFulcrum()

	// Crea una conexión gRPC al servidor Fulcrum
	conn, err := grpc.Dial(fmt.Sprintf(puerto), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	clienteFulcrum := pb.NewMiServicioClient(conn)

	// Envia la solicitud AgregarBase a Fulcrum
	respuestaFulcrum, err := clienteFulcrum.AgregarBase(ctx, req)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Respuesta del fulcrum:%s %s\n", puerto, respuestaFulcrum.Mensaje)
	return &pb.Respuesta{Mensaje: "Comando AgregarBase ejecutado", Exitoso: true}, nil
}

func (s *baseServiceServer) RenombrarBase(ctx context.Context, req *pb.RenombrarBaseRequest) (*pb.Respuesta, error) {
	puerto := s.RandomFulcrum()

	// Crea una conexión gRPC al servidor Fulcrum
	conn, err := grpc.Dial(fmt.Sprintf(puerto), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	clienteFulcrum := pb.NewMiServicioClient(conn)

	// Envia la solicitud RenombrarBase a Fulcrum
	respuestaFulcrum, err := clienteFulcrum.RenombrarBase(ctx, req)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Respuesta del fulcrum:%s %s\n", puerto, respuestaFulcrum.Mensaje)
	return &pb.Respuesta{Mensaje: "Comando RenombrarBase ejecutado", Exitoso: true}, nil
}

func (s *baseServiceServer) ActualizarValor(ctx context.Context, req *pb.ActualizarValorRequest) (*pb.Respuesta, error) {
	puerto := s.RandomFulcrum()

	// Crea una conexión gRPC al servidor Fulcrum
	conn, err := grpc.Dial(fmt.Sprintf(puerto), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	clienteFulcrum := pb.NewMiServicioClient(conn)

	// Envia la solicitud ActualizarValor a Fulcrum
	respuestaFulcrum, err := clienteFulcrum.ActualizarValor(ctx, req)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Respuesta del fulcrum:%s %s\n", puerto, respuestaFulcrum.Mensaje)
	return &pb.Respuesta{Mensaje: "Comando ActualizarValor ejecutado", Exitoso: true}, nil
}

func (s *baseServiceServer) BorrarBase(ctx context.Context, req *pb.BorrarBaseRequest) (*pb.Respuesta, error) {
	puerto := s.RandomFulcrum()

	// Crea una conexión gRPC al servidor Fulcrum
	conn, err := grpc.Dial(fmt.Sprintf(puerto), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	clienteFulcrum := pb.NewMiServicioClient(conn)

	// Envia la solicitud BorrarBase a Fulcrum
	respuestaFulcrum, err := clienteFulcrum.BorrarBase(ctx, req)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Respuesta del fulcrum:%s %s\n", puerto, respuestaFulcrum.Mensaje)
	return &pb.Respuesta{Mensaje: "Comando BorrarBase ejecutado", Exitoso: true}, nil
}

func (s *baseServiceServer) RandomFulcrum() string {
	fulcrum := rand.Intn(3) + 1
	switch fulcrum {
	case 1:
		return "localhost:50052"
	case 2:
		return "localhost:50053"
	case 3:
		return "localhost:50054"
	default:
		return "localhost:50053"
	}
}

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
