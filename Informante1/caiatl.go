package main

import (
	"bufio"
	pb "central/github.com/Sistemas-Distribuidos-2023-02/Grupo10-Laboratorio-3" // Asegúrate de ajustar la importación correctamente
	"context"
	"fmt"
	"log"
	"os"

	"google.golang.org/grpc"
)

func enviarComandoAgregarBase(client pb.MiServicioClient, nombreSector, nombreBase string, valor float32) {
	req := &pb.AgregarBaseRequest{
		NombreSector: nombreSector,
		NombreBase:   nombreBase,
		Valor:        valor,
	}

	resp, err := client.AgregarBase(context.Background(), req)
	if err != nil {
		log.Fatalf("Error al enviar comando AgregarBase: %v", err)
	}

	fmt.Printf("Respuesta del servidor: %s\n", resp.Mensaje)
}

func enviarComandoRenombrarBase(client pb.MiServicioClient, nombreSector, nombreBase string, valor interface{}) {
	req := &pb.RenombrarBaseRequest{
		NombreSector: nombreSector,
		NombreBase:   nombreBase,
		NuevoNombre:  fmt.Sprintf("%v", valor),
	}

	resp, err := client.RenombrarBase(context.Background(), req)
	if err != nil {
		log.Fatalf("Error al enviar comando RenombrarBase: %v", err)
	}

	fmt.Printf("Respuesta del servidor: %s\n", resp.Mensaje)
}
func enviarComandoActualizarValor(client pb.MiServicioClient, nombreSector, nombreBase string, nuevoValor float32) {
	req := &pb.ActualizarValorRequest{
		NombreSector: nombreSector,
		NombreBase:   nombreBase,
		NuevoValor:   nuevoValor,
	}

	resp, err := client.ActualizarValor(context.Background(), req)
	if err != nil {
		log.Fatalf("Error al enviar comando ActualizarBase: %v", err)
	}

	fmt.Printf("Respuesta del servidor: %s\n", resp.Mensaje)
}
func enviarComandoBorrarBase(client pb.MiServicioClient, nombreSector, nombreBase string) {
	req := &pb.BorrarBaseRequest{
		NombreSector: nombreSector,
		NombreBase:   nombreBase,
	}

	resp, err := client.BorrarBase(context.Background(), req)
	if err != nil {
		log.Fatalf("Error al enviar comando BorrarBase: %v", err)
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
	fmt.Println("Ejemplos:")
	fmt.Println(" - AgregarBase: AgregarBase nombre_sector nombre_base 42.0")
	fmt.Println(" - RenombrarBase: RenombrarBase nombre_sector nombre_base nuevo_nombre")
	fmt.Println(" - ActualizarValor: ActualizarValor nombre_sector nombre_base 43.0")
	fmt.Println(" - BorrarBase: BorrarBase nombre_sector nombre_base")

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		entrada := scanner.Text()
		// Parsea la entrada para obtener los parámetros del comando
		// Aquí deberías implementar la lógica para parsear la entrada y llamar a la función correspondiente
		// Puedes usar strings.Split o alguna otra técnica de parsing

		// Ejemplo de parsing (puede necesitar ser modificado según tus necesidades):
		var comando, nombreSector, nombreBase string
		var valor interface{}
		n, _ := fmt.Sscanf(entrada, "%s %s %s %v", &comando, &nombreSector, &nombreBase, &valor)
		if n >= 3 {
			switch comando {
			case "AgregarBase":
				enviarComandoAgregarBase(client, nombreSector, nombreBase, valor)
			case "RenombrarBase":
				enviarComandoRenombrarBase(client, nombreSector, nombreBase, valor)
			case "ActualizarValor":
				enviarComandoActualizarValor(client, nombreSector, nombreBase, valor)
			case "BorrarBase":
				enviarComandoBorrarBase(client, nombreSector, nombreBase)
			default:
				fmt.Println("Comando no reconocido")
			}
		} else {
			fmt.Println("Entrada no válida")
		}

		fmt.Println("Ingrese otro comando:")
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error al leer la entrada estándar: %v", err)
	}
}
