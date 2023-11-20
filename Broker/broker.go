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
	respuesta := &pb.Respuesta{Mensaje: puerto, Exitoso: true}
	fmt.Printf("Solicitud de AgregarBase recibida, Se envían los datos del puerto %s\n", puerto)
	return respuesta, nil
}

func (s *baseServiceServer) RenombrarBase(ctx context.Context, req *pb.RenombrarBaseRequest) (*pb.Respuesta, error) {
	puerto := s.RandomFulcrum()
	respuesta := &pb.Respuesta{Mensaje: puerto, Exitoso: true}
	fmt.Printf("Solicitud de RenombrarBase recibida, Se envían los datos del puerto %s\n", puerto)
	return respuesta, nil
}

func (s *baseServiceServer) ActualizarValor(ctx context.Context, req *pb.ActualizarValorRequest) (*pb.Respuesta, error) {
	puerto := s.RandomFulcrum()
	respuesta := &pb.Respuesta{Mensaje: puerto, Exitoso: true}
	fmt.Printf("Solicitud de ActualizarValor recibida, Se envían los datos del puerto %s\n", puerto)
	return respuesta, nil
}

func (s *baseServiceServer) BorrarBase(ctx context.Context, req *pb.BorrarBaseRequest) (*pb.Respuesta, error) {
	puerto := s.RandomFulcrum()
	respuesta := &pb.Respuesta{Mensaje: puerto, Exitoso: true}
	fmt.Printf("Solicitud de BorrarBase recibida, Se envían los datos del puerto %s\n", puerto)
	return respuesta, nil
}

func (s *baseServiceServer) GetSoldados(ctx context.Context, req *pb.GetSoldadosRequest) (*pb.Respuesta, error) {
	puerto := s.RandomFulcrum()

	//Conectar a fc1
	conn, err := grpc.Dial(fmt.Sprintf(puerto), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	clienteFulcrum := pb.NewMiServicioClient(conn)

	// Envia la solicitud GetSoldados a Fulcrum
	respuestaFulcrum, err := clienteFulcrum.GetSoldados(ctx, req)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Respuesta del fulcrum:%s %s\n", puerto, respuestaFulcrum.Mensaje)

	resultado := fmt.Sprintf("Cantidad: %s", respuestaFulcrum.Mensaje)

	return &pb.Respuesta{Mensaje: resultado, Exitoso: true}, nil
}

func (s *baseServiceServer) RandomFulcrum() string {
	fulcrum := rand.Intn(3) + 1
	switch fulcrum {
	case 1:
		return "dist130.inf.santiago.usm.cl:50052"
	case 2:
		return "dist131.inf.santiago.usm.cl:50053"
	case 3:
		return "dist132.inf.santiago.usm.cl:50054"
	default:
		return "dist131.inf.santiago.usm.cl:50053"
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
