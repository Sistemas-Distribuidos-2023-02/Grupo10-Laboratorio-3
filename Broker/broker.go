package main

import (
	pb "central/github.com/Sistemas-Distribuidos-2023-02/Grupo10-Laboratorio-3" // Asegúrate de ajustar la importación correctamente
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strings"

	"google.golang.org/grpc"
)

type baseServiceServer struct {
	pb.UnimplementedMiServicioServer
}

func (s *baseServiceServer) AgregarBase(ctx context.Context, req *pb.AgregarBaseRequest) (*pb.Respuesta, error) {
	puerto := s.RandomFulcrum()
	respuesta := &pb.Respuesta{Mensaje: puerto, Exitoso: true}
	fmt.Printf("Solicitud de AgregarBase recibida, Se envían los datos del puerto %s", puerto)
	return respuesta, nil
}

func (s *baseServiceServer) RenombrarBase(ctx context.Context, req *pb.RenombrarBaseRequest) (*pb.Respuesta, error) {
	puerto := s.RandomFulcrum()
	respuesta := &pb.Respuesta{Mensaje: puerto, Exitoso: true}
	fmt.Printf("Solicitud de RenombrarBase recibida, Se envían los datos del puerto %s", puerto)
	return respuesta, nil
}

func (s *baseServiceServer) ActualizarValor(ctx context.Context, req *pb.ActualizarValorRequest) (*pb.Respuesta, error) {
	puerto := s.RandomFulcrum()
	respuesta := &pb.Respuesta{Mensaje: puerto, Exitoso: true}
	fmt.Printf("Solicitud de ActualizarValor recibida, Se envían los datos del puerto %s", puerto)
	return respuesta, nil
}

func (s *baseServiceServer) BorrarBase(ctx context.Context, req *pb.BorrarBaseRequest) (*pb.Respuesta, error) {
	puerto := s.RandomFulcrum()
	respuesta := &pb.Respuesta{Mensaje: puerto, Exitoso: true}
	fmt.Printf("Solicitud de BorrarBase recibida, Se envían los datos del puerto %s", puerto)
	return respuesta, nil
}

func (s *baseServiceServer) GetSoldados(ctx context.Context, req *pb.GetSoldadosRequest) (*pb.Respuesta, error) {
	puertofc1 := "localhost:50052"
	puertofc2 := "localhost:50053"
	puertofc3 := "localhost:50054"

	//Conectar a fc1
	conn1, err1 := grpc.Dial(fmt.Sprintf(puertofc1), grpc.WithInsecure())
	if err1 != nil {
		return nil, err1
	}
	defer conn1.Close()
	clienteFulcrum1 := pb.NewMiServicioClient(conn1)

	// Envia la solicitud GetSoldados a Fulcrum
	respuestaFulcrum1, err11 := clienteFulcrum1.GetSoldados(ctx, req)
	if err11 != nil {
		return nil, err11
	}
	fmt.Printf("Respuesta del fulcrum1:%s %s\n", puertofc1, respuestaFulcrum1.Mensaje)

	resultado1 := fmt.Sprintf("ubicación: Fulcrum1, cantidad: %s", respuestaFulcrum1.Mensaje)

	//Conectar a fc2
	conn2, err2 := grpc.Dial(fmt.Sprintf(puertofc2), grpc.WithInsecure())
	if err2 != nil {
		return nil, err2
	}
	defer conn2.Close()
	clienteFulcrum2 := pb.NewMiServicioClient(conn2)

	// Envia la solicitud GetSoldados a Fulcrum
	respuestaFulcrum2, err22 := clienteFulcrum2.GetSoldados(ctx, req)
	if err22 != nil {
		return nil, err22
	}
	fmt.Printf("Respuesta del fulcrum2:%s %s\n", puertofc2, respuestaFulcrum2.Mensaje)

	resultado2 := fmt.Sprintf("ubicación: Fulcrum2, cantidad: %s", respuestaFulcrum2.Mensaje)

	//Conectar a fc3
	conn3, err3 := grpc.Dial(fmt.Sprintf(puertofc3), grpc.WithInsecure())
	if err3 != nil {
		return nil, err3
	}
	defer conn3.Close()
	clienteFulcrum3 := pb.NewMiServicioClient(conn3)

	// Envia la solicitud GetSoldados a Fulcrum
	respuestaFulcrum3, err33 := clienteFulcrum3.GetSoldados(ctx, req)
	if err33 != nil {
		return nil, err33
	}
	fmt.Printf("Respuesta del fulcrum3:%s %s\n", puertofc3, respuestaFulcrum3.Mensaje)

	resultado3 := fmt.Sprintf("ubicación: Fulcrum3, cantidad: %s", respuestaFulcrum3.Mensaje)

	var resultados []string

	// Verificar si la base fue encontrada en cada Fulcrum
	if respuestaFulcrum1.Mensaje != "Base no encontrada en comando GetSoldados" {
		resultados = append(resultados, resultado1)
	}
	if respuestaFulcrum2.Mensaje != "Base no encontrada en comando GetSoldados" {
		resultados = append(resultados, resultado2)
	}
	if respuestaFulcrum3.Mensaje != "Base no encontrada en comando GetSoldados" {
		resultados = append(resultados, resultado3)
	}

	var respuestaFinal string

	// Construir la respuesta final
	if len(resultados) > 0 {
		respuestaFinal = fmt.Sprintf("Comando GetSoldados ejecutado, encontrado en: %s", strings.Join(resultados, ", "))
	} else {
		respuestaFinal = "Comando GetSoldados ejecutado, no encontrado en ninguno"
	}

	return &pb.Respuesta{Mensaje: respuestaFinal, Exitoso: true}, nil
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
