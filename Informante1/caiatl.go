package main

import (
	pb "central/github.com/Sistemas-Distribuidos-2023-02/Grupo10-Laboratorio-3" // Asegúrate de ajustar la importación correctamente
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
)

func enviarComandoAgregarBase(client pb.MiServicioClient) {
	req := &pb.AgregarBaseRequest{
		// Establece los campos necesarios para el comando AgregarBase
	}

	resp, err := client.AgregarBase(context.Background(), req)
	if err != nil {
		log.Fatalf("Error al enviar comando AgregarBase: %v", err)
	}

	fmt.Printf("Respuesta del servidor: %s\n", resp.Mensaje)
}

// Implementa funciones similares para los demás comandos

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Error al conectar al servidor gRPC: %v", err)
	}
	defer conn.Close()

	client := pb.NewMiServicioClient(conn)

	// Llama a las funciones para enviar comandos según sea necesario
	enviarComandoAgregarBase(client)
}
