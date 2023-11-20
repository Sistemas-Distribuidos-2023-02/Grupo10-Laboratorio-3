package main

import (
	"bufio"
	pb "central/github.com/Sistemas-Distribuidos-2023-02/Grupo10-Laboratorio-3"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
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

func MonotonicWrites(sector, reloj, puerto string) string {
	contenidos, err := ioutil.ReadFile("registro.txt")
	if err != nil {
		return "problemas al abrir el archivo "
	}
	// Separar el contenido por líneas
	lineas := strings.Split(string(contenidos), "\n")
	for i := range lineas {
		// Verificar si la línea está vacía
		if lineas[i] == "" {
			continue
		}

		// Separar elementos del archivo
		partes := strings.Split(lineas[i], ",")
		if len(partes) < 6 {
			return "Error de formato en el archivo"
		}

		sectorA := partes[1]
		relojA := partes[4]
		puertoA := partes[5]

		// Separar reloj archivo
		separar := strings.Split(strings.Trim(relojA, "[]"), ",")
		numero1, err := strconv.Atoi(strings.TrimSpace(separar[0]))
		if err != nil {
			return "problemas al abrir el archivo "
		}
		numero2, err := strconv.Atoi(strings.TrimSpace(separar[1]))
		if err != nil {
			return "problemas al abrir el archivo "
		}
		numero3, err := strconv.Atoi(strings.TrimSpace(separar[2]))
		if err != nil {
			return "problemas al abrir el archivo "
		}

		// Separar reloj a agregar
		separar = strings.Split(strings.Trim(reloj, "[]"), ",")
		numero4, err := strconv.Atoi(strings.TrimSpace(separar[0]))
		if err != nil {
			return "problemas al abrir el archivo "
		}
		numero5, err := strconv.Atoi(strings.TrimSpace(separar[1]))
		if err != nil {
			return "problemas al abrir el archivo "
		}
		numero6, err := strconv.Atoi(strings.TrimSpace(separar[2]))
		if err != nil {
			return "problemas al abrir el archivo "
		}

		if sectorA == sector && puertoA == puerto {
			if numero1 > numero4 || numero2 > numero5 || numero3 > numero6 {
				return "Error de consistencia"
			}
		}
	}

	return "No hay problemas"
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
		puerto = "dist039:50053"
	case " fulcrum3":
		puerto = "dist040:50054"
	}

	fmt.Printf("Respuesta del servidor: %s soldados\n", soldados)

	monotonic := MonotonicWrites(nombreSector, reloj, puerto)
	if monotonic == "No hay problemas" {
		// Escribir en registro.txt
		logfile, err := os.OpenFile("registro.txt", os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Printf("Log de registro no pudo abrirse exitosamente")
		}
		defer logfile.Close()

		_, err = fmt.Fprintf(logfile, "GetSoldados,%s,%s,%s,%s,%s\n", nombreSector, nombreBase, soldados, reloj, puerto)
		if err != nil {
			fmt.Printf("No pudo escribirse correctamente en archivo log")
		}
	} else {
		fmt.Printf("No se pudo escribir en el archivo registro.txt debido a que hubo un error de consistencia detectado por Monotonic Writes \n")
	}
}

func main() {

	conn, err := grpc.Dial("dist037:50051", grpc.WithInsecure())
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
		fmt.Println("Ingrese otro comando:")
	}

}
