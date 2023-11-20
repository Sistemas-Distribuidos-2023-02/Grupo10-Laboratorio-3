package main

import (
	"bufio"
	pb "central/github.com/Sistemas-Distribuidos-2023-02/Grupo10-Laboratorio-3"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"google.golang.org/grpc"
)

type baseServiceServer struct {
	pb.UnimplementedMiServicioServer
}

var PrimeraEscritura = true

func inicializarArchivo() error {
	// Reinicia el contenido del archivo registro
	file, err := os.OpenFile("Registro.txt", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Error al abrir el archivo registro.txt: %v", err)
		return err
	}
	defer file.Close()
	return nil
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

	if PrimeraEscritura {
		err = inicializarArchivo()
		if err != nil {
			log.Fatal("Problemas con el archivo registro.txt")
		}
		PrimeraEscritura = false
	}

	// Separar elementos de la respuesta
	partes := strings.Split(resp.Mensaje, "-")
	soldados := partes[0]
	reloj := partes[1]
	fulcrum := partes[2]
	puerto := ""
	switch fulcrum {
	case " fulcrum1":
		puerto = "localhost:50052"
	case " fulcrum2":
		puerto = "localhost:50053"
	case " fulcrum3":
		puerto = "localhost:50054"
	}

	// Escribir en registro.txt
	logfile, err := os.OpenFile("registro.txt", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Log de registro no pudo abrirse exitosamente")
	}
	defer logfile.Close()

	_, err = fmt.Fprintf(logfile, "GetSoldados %s %s %s %s %s\n", nombreSector, nombreBase, soldados, reloj, puerto)
	if err != nil {
		fmt.Printf("No pudo escribirse correctamente en archivo log")
	}

	fmt.Printf("Respuesta del servidor: %s soldados\n", soldados)
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
