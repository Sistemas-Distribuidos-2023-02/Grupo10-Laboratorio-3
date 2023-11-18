package main

import (
	"bufio"
	pb "central/github.com/Sistemas-Distribuidos-2023-02/Grupo10-Laboratorio-3" // Asegúrate de ajustar la importación correctamente
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"google.golang.org/grpc"
)

var PrimeraEscritura = true

func inicializarArchivo() error {
	// Reinicia el contenido del archivo registro
	file, err := os.OpenFile("Registro.txt", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Error al abrir el archivo DATA.txt: %v", err)
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
	// ASEGURA EL MODELO READ YOUR WRITES
	if respFulcrum.Exitoso == false {
		fmt.Printf("Sucedió un error al momento de ejecutar el comando AgregarBase")
		return
	}
	fulcrum := asignarNombreFulcrum(puerto)
	if PrimeraEscritura == true {
		err := inicializarArchivo()
		if err != nil {
			log.Fatal("Problemas con el archivo registro.txt")
		}
		PrimeraEscritura = false
	}
	//Escritura en el registro.txt
	logfile, err := os.OpenFile("registro.txt", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Log de registro no pudo abrirse exitosamente")
	}
	defer logfile.Close()

	_, err = fmt.Fprintf(logfile, "AgregarBase %s %s %.0f\n", req.NombreSector, req.NombreBase, req.Valor)
	if err != nil {
		fmt.Printf("No pudo escribirse correctamente en archivo log")
	}
	fmt.Printf("Respuesta del %s: %s\n", fulcrum, respFulcrum.Mensaje)
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
	// ASEGURA EL MODELO READ YOUR WRITES
	if respFulcrum.Exitoso == false {
		fmt.Printf("Sucedió un error al momento de ejecutar el comando AgregarBase")
		return
	}
	fulcrum := asignarNombreFulcrum(puerto)
	if PrimeraEscritura {
		err := inicializarArchivo()
		if err != nil {
			log.Fatal("Problemas con el archivo registro.txt")
		}
		PrimeraEscritura = false
	}
	//Escritura en el registro.txt
	logfile, err := os.OpenFile("registro.txt", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Log de registro no pudo abrirse exitosamente")
	}
	defer logfile.Close()

	_, err = fmt.Fprintf(logfile, "RenombrarBase %s %s %s\n", req.NombreSector, req.NombreBase, req.NuevoNombre)
	if err != nil {
		fmt.Printf("No pudo escribirse correctamente en archivo log")
	}
	fmt.Printf("Respuesta del %s: %s\n", fulcrum, respFulcrum.Mensaje)
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
	// ASEGURA EL MODELO READ YOUR WRITES
	if respFulcrum.Exitoso == false {
		fmt.Printf("Sucedió un error al momento de ejecutar el comando AgregarBase")
		return
	}
	fulcrum := asignarNombreFulcrum(puerto)
	if PrimeraEscritura {
		err := inicializarArchivo()
		if err != nil {
			log.Fatal("Problemas con el archivo registro.txt")
		}
		PrimeraEscritura = false
	}
	//Escritura en el registro.txt
	logfile, err := os.OpenFile("registro.txt", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Log de registro no pudo abrirse exitosamente")
	}
	defer logfile.Close()

	_, err = fmt.Fprintf(logfile, "ActualizarValor %s %s %.0f\n", req.NombreSector, req.NombreBase, req.NuevoValor)
	if err != nil {
		fmt.Printf("No pudo escribirse correctamente en archivo log")
	}
	fmt.Printf("Respuesta del %s: %s\n", fulcrum, respFulcrum.Mensaje)
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
	// ASEGURA EL MODELO READ YOUR WRITES
	if respFulcrum.Exitoso == false {
		fmt.Printf("Sucedió un error al momento de ejecutar el comando AgregarBase")
		return
	}
	fulcrum := asignarNombreFulcrum(puerto)
	if PrimeraEscritura {
		err := inicializarArchivo()
		if err != nil {
			log.Fatal("Problemas con el archivo registro.txt")
		}
		PrimeraEscritura = false
	}
	//Escritura en el registro.txt
	logfile, err := os.OpenFile("registro.txt", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Log de registro no pudo abrirse exitosamente")
	}
	defer logfile.Close()

	_, err = fmt.Fprintf(logfile, "BorrarBase %s %s\n", req.NombreSector, req.NombreBase)
	if err != nil {
		fmt.Printf("No pudo escribirse correctamente en archivo log")
	}
	fmt.Printf("Respuesta del %s: %s\n", fulcrum, respFulcrum.Mensaje)
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

	fmt.Println("Ingrese un comando:")
	fmt.Println("Ejemplos:")
	fmt.Println(" - AgregarBase: AgregarBase nombre_sector nombre_base 42")
	fmt.Println(" - RenombrarBase: RenombrarBase nombre_sector nombre_base nuevo_nombre")
	fmt.Println(" - ActualizarValor: ActualizarValor nombre_sector nombre_base 43")
	fmt.Println(" - BorrarBase: BorrarBase nombre_sector nombre_base")

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		entrada := scanner.Text()
		// Parsea la entrada para obtener los parámetros del comando
		// Aquí deberías implementar la lógica para parsear la entrada y llamar a la función correspondiente
		// Puedes usar strings.Split o alguna otra técnica de parsing

		// Ejemplo de parsing (puede necesitar ser modificado según tus necesidades):
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
