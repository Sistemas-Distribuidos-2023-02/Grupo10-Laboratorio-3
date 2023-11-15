package main

import (
	pb "central/github.com/Sistemas-Distribuidos-2023-02/Grupo10-Laboratorio-3" // Asegúrate de ajustar la importación correctamente
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
)

type baseServiceServer struct {
	pb.UnimplementedMiServicioServer
}

// CrearRegistro crea un nuevo archivo de registro para el sector
func (s *baseServiceServer) CrearRegistro(sectorFileName string) error {
	// Lógica para crear un nuevo archivo de registro para el sector
	file, err := os.Create(sectorFileName)
	if err != nil {
		return err
	}
	file.Close()
	return nil
}

func (s *baseServiceServer) AgregarBase(ctx context.Context, req *pb.AgregarBaseRequest) (*pb.Respuesta, error) {
	nombresector := fmt.Sprintf("Sector%s.txt", req.NombreSector)
	if _, err := os.Stat(nombresector); os.IsNotExist(err) {
		// El archivo no existe, entonces se crea uno nuevo
		if err := s.CrearRegistro(nombresector); err != nil {
			return &pb.Respuesta{Mensaje: "Comando AgregarBase no pudo ser ejecutado", Exitoso: false}, err
		}
	}
	file, err := os.OpenFile(nombresector, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return &pb.Respuesta{Mensaje: "Comando AgregarBase no pudo ser ejecutado", Exitoso: false}, err
	}
	defer file.Close()
	// Escribir la información de la base en el archivo
	_, err = fmt.Fprintf(file, "Sector %s %s %.0f\n", req.NombreSector, req.NombreBase, req.Valor)
	if err != nil {
		return &pb.Respuesta{Mensaje: "Comando AgregarBase no pudo ser ejecutado", Exitoso: false}, err
	}
	return &pb.Respuesta{Mensaje: "Comando AgregarBase ejecutado", Exitoso: true}, nil
}

func (s *baseServiceServer) RenombrarBase(ctx context.Context, req *pb.RenombrarBaseRequest) (*pb.Respuesta, error) {
	// Lógica para el comando AgregarBase
	// ...

	return &pb.Respuesta{Mensaje: "Comando RenombrarBase ejecutado", Exitoso: true}, nil
}

func (s *baseServiceServer) ActualizarValor(ctx context.Context, req *pb.ActualizarValorRequest) (*pb.Respuesta, error) {
	// Lógica para el comando AgregarBase
	// ...

	return &pb.Respuesta{Mensaje: "Comando ActualizaValor ejecutado", Exitoso: true}, nil
}

func (s *baseServiceServer) BorrarBase(ctx context.Context, req *pb.BorrarBaseRequest) (*pb.Respuesta, error) {
	// Lógica para el comando AgregarBase
	// ...

	return &pb.Respuesta{Mensaje: "Comando BorrarBase ejecutado", Exitoso: true}, nil
}

func main() {
	listener, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatalf("Error al escuchar: %v", err)
	}

	server := grpc.NewServer()
	pb.RegisterMiServicioServer(server, &baseServiceServer{})

	fmt.Println("Servidor gRPC iniciado en el puerto 50053")

	if err := server.Serve(listener); err != nil {
		log.Fatalf("Error al servir: %v", err)
	}
}
