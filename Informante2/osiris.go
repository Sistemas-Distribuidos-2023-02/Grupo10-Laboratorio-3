package main

import (
	"bufio"
	pb "central/github.com/Sistemas-Distribuidos-2023-02/Grupo10-Laboratorio-3" // Asegúrate de ajustar la importación correctamente
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"google.golang.org/grpc"
)

var PrimeraEscritura = true

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
		partes := strings.Split(lineas[i], " ")
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

	puerto := resp.Mensaje
	conn, err := grpc.Dial(puerto, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Error al conectar con el fulcrum: %v", err)
	}
	defer conn.Close()
	clienteFulcrum := pb.NewMiServicioClient(conn)

	respFulcrum, err := clienteFulcrum.AgregarBase(context.Background(), req)
	if err != nil {
		log.Fatalf("Error al enviar comando AgregarBase al fulcrum: %v", err)
	}
	fulcrum := asignarNombreFulcrum(puerto)
	if PrimeraEscritura == true {
		err := inicializarArchivo()
		if err != nil {
			log.Fatal("Problemas con el archivo registro.txt")
		}
		PrimeraEscritura = false
	}
	reloj := respFulcrum.Mensaje
	//Metodo Monotonic Writes
	monotonic := MonotonicWrites(req.NombreSector, reloj, puerto)
	if monotonic == "No hay problemas" {
		//Escritura en el registro.txt
		logfile, err := os.OpenFile("registro.txt", os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Printf("Log de registro no pudo abrirse exitosamente")
		}
		defer logfile.Close()

		_, err = fmt.Fprintf(logfile, "AgregarBase %s %s %.0f %s %s\n", req.NombreSector, req.NombreBase, req.Valor, reloj, puerto)
		if err != nil {
			fmt.Printf("No pudo escribirse correctamente en archivo log")
		}
		if resp.Exitoso {
			fmt.Printf("Respuesta del %s: Comando AgregarBase Ejecutado con éxito\n", fulcrum)
		} else {
			fmt.Printf("Respuesta del %s: Comando AgregarBase no pudo ser ejecutado\n", fulcrum)
		}
	} else {
		fmt.Printf("No se pudo escribir en el archivo registro.txt debido a que hubo un error de consistencia detectado por Monotonic Writes \n")
	}
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

	puerto := resp.Mensaje
	conn, err := grpc.Dial(puerto, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Error al conectar con el fulcrum: %v", err)
	}
	defer conn.Close()
	clienteFulcrum := pb.NewMiServicioClient(conn)

	respFulcrum, err := clienteFulcrum.RenombrarBase(context.Background(), req)
	if err != nil {
		log.Fatalf("Error al enviar comando RenombrarBase al fulcrum: %v", err)
	}
	fulcrum := asignarNombreFulcrum(puerto)
	if PrimeraEscritura {
		err := inicializarArchivo()
		if err != nil {
			log.Fatal("Problemas con el archivo registro.txt")
		}
		PrimeraEscritura = false
	}
	reloj := respFulcrum.Mensaje
	//Metodo Monotonic Writes
	monotonic := MonotonicWrites(req.NombreSector, reloj, puerto)
	if monotonic == "No hay problemas" {
		//Escritura en el registro.txt
		logfile, err := os.OpenFile("registro.txt", os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Printf("Log de registro no pudo abrirse exitosamente")
		}
		defer logfile.Close()

		_, err = fmt.Fprintf(logfile, "RenombrarBase %s %s %s %s %s\n", req.NombreSector, req.NombreBase, req.NuevoNombre, reloj, puerto)
		if err != nil {
			fmt.Printf("No pudo escribirse correctamente en archivo log")
		}
		if resp.Exitoso {
			fmt.Printf("Respuesta del %s: Comando RenombrarBase Ejecutado con éxito\n", fulcrum)
		} else {
			fmt.Printf("Respuesta del %s: Comando RenombrarBase no pudo ser ejecutado\n", fulcrum)
		}
	} else {
		fmt.Printf("No se pudo escribir en el archivo registro.txt debido a que hubo un error de consistencia detectado por Monotonic Writes \n")
	}
}

func enviarComandoActualizarValor(client pb.MiServicioClient, nombreSector, nombreBase string, nuevoValor float32) {
	req := &pb.ActualizarValorRequest{
		NombreSector: nombreSector,
		NombreBase:   nombreBase,
		NuevoValor:   nuevoValor,
	}

	resp, err := client.ActualizarValor(context.Background(), req)
	if err != nil {
		log.Fatalf("Error al enviar comando ActualizarValor: %v", err)
	}

	puerto := resp.Mensaje
	conn, err := grpc.Dial(puerto, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Error al conectar con el fulcrum: %v", err)
	}
	defer conn.Close()
	clienteFulcrum := pb.NewMiServicioClient(conn)

	respFulcrum, err := clienteFulcrum.ActualizarValor(context.Background(), req)
	if err != nil {
		log.Fatalf("Error al enviar comando ActualizarValor al fulcrum: %v", err)
	}
	fulcrum := asignarNombreFulcrum(puerto)
	if PrimeraEscritura {
		err := inicializarArchivo()
		if err != nil {
			log.Fatal("Problemas con el archivo registro.txt")
		}
		PrimeraEscritura = false
	}
	reloj := respFulcrum.Mensaje
	//Metodo Monotonic Writes
	monotonic := MonotonicWrites(req.NombreSector, reloj, puerto)
	if monotonic == "No hay problemas" {
		//Escritura en el registro.txt
		logfile, err := os.OpenFile("registro.txt", os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Printf("Log de registro no pudo abrirse exitosamente")
		}
		defer logfile.Close()

		_, err = fmt.Fprintf(logfile, "ActualizarValor %s %s %.0f %s %s\n", req.NombreSector, req.NombreBase, req.NuevoValor, reloj, puerto)
		if err != nil {
			fmt.Printf("No pudo escribirse correctamente en archivo log")
		}
		if resp.Exitoso {
			fmt.Printf("Respuesta del %s: Comando ActualizarValor Ejecutado con éxito\n", fulcrum)
		} else {
			fmt.Printf("Respuesta del %s: Comando ActualizarValor no pudo ser ejecutado\n", fulcrum)
		}
	} else {
		fmt.Printf("No se pudo escribir en el archivo registro.txt debido a que hubo un error de consistencia detectado por Monotonic Writes \n")
	}

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

	puerto := resp.Mensaje
	conn, err := grpc.Dial(puerto, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Error al conectar con el fulcrum: %v", err)
	}
	defer conn.Close()
	clienteFulcrum := pb.NewMiServicioClient(conn)

	respFulcrum, err := clienteFulcrum.BorrarBase(context.Background(), req)
	if err != nil {
		log.Fatalf("Error al enviar comando BorrarBase al fulcrum: %v", err)
	}
	fulcrum := asignarNombreFulcrum(puerto)
	if PrimeraEscritura {
		err := inicializarArchivo()
		if err != nil {
			log.Fatal("Problemas con el archivo registro.txt")
		}
		PrimeraEscritura = false
	}
	reloj := respFulcrum.Mensaje
	//Metodo Monotonic Writes
	monotonic := MonotonicWrites(req.NombreSector, reloj, puerto)
	if monotonic == "No hay problemas" {
		//Escritura en el registro.txt
		logfile, err := os.OpenFile("registro.txt", os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Printf("Log de registro no pudo abrirse exitosamente")
		}
		defer logfile.Close()

		_, err = fmt.Fprintf(logfile, "BorrarBase %s %s %s %s\n", req.NombreSector, req.NombreBase, reloj, puerto)
		if err != nil {
			fmt.Printf("No pudo escribirse correctamente en archivo log")
		}
		if resp.Exitoso {
			fmt.Printf("Respuesta del %s: Comando BorrarBase Ejecutado con éxito\n", fulcrum)
		} else {
			fmt.Printf("Respuesta del %s: Comando BorrarBase no pudo ser ejecutado\n", fulcrum)
		}
	} else {
		fmt.Printf("No se pudo escribir en el archivo registro.txt debido a que hubo un error de consistencia detectado por Monotonic Writes \n")
	}
}

func asignarNombreFulcrum(puerto string) string {
	var fulcrum string

	switch puerto {
	case "localhost:50052":
		fulcrum = "fulcrum1"
	case "localhost:50053":
		fulcrum = "fulcrum2"
	case "localhost:50054":
		fulcrum = "fulcrum3"
	}

	return fulcrum
}

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Error al conectar al servidor gRPC: %v", err)
	}
	defer conn.Close()

	client := pb.NewMiServicioClient(conn)

	fmt.Println("Ingrese algun comando:")
	fmt.Println("-AgregarBase nombre_sector nombre_base valor")
	fmt.Println("-RenombrarBase nombre_sector nombre_base nuevo_nombre")
	fmt.Println("-ActualizarValor nombre_sector nombre_base nuevo_valor")
	fmt.Println("-BorrarBase nombre_sector nombre_base")

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		entrada := scanner.Text()
		// Parsea la entrada para obtener los parámetros del comando
		var comando, nombreSector, nombreBase, nuevaBase string
		n, _ := fmt.Sscanf(entrada, "%s %s %s %s", &comando, &nombreSector, &nombreBase, &nuevaBase)
		if n >= 3 {
			if floatValue, err := strconv.ParseFloat(nuevaBase, 32); err == nil {
				valor := float32(floatValue)
				switch comando {
				case "AgregarBase":
					enviarComandoAgregarBase(client, nombreSector, nombreBase, valor)
				case "ActualizarValor":
					enviarComandoActualizarValor(client, nombreSector, nombreBase, valor)
				case "BorrarBase":
					enviarComandoBorrarBase(client, nombreSector, nombreBase)
				default:
					fmt.Println("Comando no reconocido")
				}
			} else {
				switch comando {
				case "RenombrarBase":
					enviarComandoRenombrarBase(client, nombreSector, nombreBase, nuevaBase)
				case "BorrarBase":
					enviarComandoBorrarBase(client, nombreSector, nombreBase)
				default:
					fmt.Println("Comando no reconocido")
				}
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
