package main

import (
	"bufio"
	pb "central/github.com/Sistemas-Distribuidos-2023-02/Grupo10-Laboratorio-3"
	"context"
	"fmt"
	"log"
	"os"

	"google.golang.org/grpc"
)

type baseServiceServer struct {
	pb.UnimplementedMiServicioServer
}

func enviarComandoGetSoldados(client pb.MiServicioClient, nombreSector, nombreBase string) {
	req := &pb.GetSoldadosRequest{
		NombreSector: nombreSector,
		NombreBase:   nombreBase,
	}

	resp, err := client.GetSoldados(context.Background(), req)
	if err != nil {
		log.Fatalf("Error al enviar comando GetSoldados: %v", err)
	}

	fmt.Printf("Respuesta del servidor: %s\n", resp.Mensaje)
}

func main() {

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Error al conectar al servidor gRPC: %v", err)
	}
	defer conn.Close()

	client := pb.NewMiServicioClient(conn)

	fmt.Println("Ingrese un comando:")
	fmt.Println(" - GetSoldados: GetSoldados <nombre_sector> <nombre_base>")

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		entrada := scanner.Text()
		var comando, nombreSector, nombreBase string
		n, _ := fmt.Sscanf(entrada, "%s %s %s %s", &comando, &nombreSector, &nombreBase)
		if n >= 3 {
			switch comando {
			case "GetSoldados":
				enviarComandoGetSoldados(client, nombreSector, nombreBase)
			default:
				fmt.Println("Comando no reconocido")
			}

		} else {
			fmt.Println("Entrada no válida")
		}
	}

}
