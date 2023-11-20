package main

import (
	"bufio"
	pb "central/github.com/Sistemas-Distribuidos-2023-02/Grupo10-Laboratorio-3" // Asegúrate de ajustar la importación correctamente
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"google.golang.org/grpc"
)

type baseServiceServer struct {
	pb.UnimplementedMiServicioServer
}

func retornarReloj(nombreArchivo string) (string, error) {
	archivo, err := os.Open(nombreArchivo)
	if err != nil {
		return "", err
	}
	defer archivo.Close()

	scanner := bufio.NewScanner(archivo)

	if scanner.Scan() {
		return scanner.Text(), nil
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	// Si el archivo está vacío
	return "", fmt.Errorf("El archivo está vacío")
}

func editarReloj(nombreArchivo string) error {
	// Leer el contenido actual del archivo
	contenidos, err := ioutil.ReadFile(nombreArchivo)
	if err != nil {
		return err
	}

	// Separar el contenido por líneas
	lineas := strings.Split(string(contenidos), "\n")
	// Separar el elemento a sumar
	nuevaL := lineas[0]
	separar := strings.Split(strings.Trim(nuevaL, "[]"), ",")
	numero1, err := strconv.Atoi(strings.TrimSpace(separar[0]))
	numero2, err := strconv.Atoi(strings.TrimSpace(separar[1]))
	numero3, err := strconv.Atoi(strings.TrimSpace(separar[2]))
	numero3 = numero3 + 1
	// Crear una nueva linea
	nuevaLinea := fmt.Sprintf("[%d,%d,%d]", numero1, numero2, numero3)

	// Abrir el archivo en modo escritura, truncándolo para eliminar el contenido original
	file, err := os.OpenFile(nombreArchivo, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Escribir el nuevo contenido en el archivo (primera línea)
	_, err = fmt.Fprintln(file, nuevaLinea)
	if err != nil {
		return err
	}

	// Escribir el resto del contenido original después de la nueva línea
	for i := 1; i < len(lineas); i++ {
		_, err := fmt.Fprintln(file, lineas[i])
		if err != nil {
			return err
		}
	}
	return nil
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

func (s *baseServiceServer) AgregarLOG(info, sector, base, nuevonombre string, valor, nuevovalor float32) (*pb.Respuesta, error) {
	// Creación del log en caso de no existir
	if _, err := os.Stat("Registro.txt"); os.IsNotExist(err) {
		if err := s.CrearRegistro("Registro.txt"); err != nil {
			return &pb.Respuesta{Mensaje: "Log de registro no pudo ser creado", Exitoso: false}, err
		}
	}
	logfile, err := os.OpenFile("Registro.txt", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return &pb.Respuesta{Mensaje: "Log de registro no pudo abrirse exitosamente", Exitoso: false}, err
	}
	defer logfile.Close()

	// Agregar Hora y formatear
	horaActual := time.Now()
	horaFormateada := horaActual.Format("15:04:05")

	switch info { // Escribir la información de la base en el log
	case "agregar":
		_, err = fmt.Fprintf(logfile, "AgregarBase %s %s %.0f [%s]\n", sector, base, valor, horaFormateada)
		if err != nil {
			return &pb.Respuesta{Mensaje: "No pudo escribirse correctamente en archivo log", Exitoso: false}, err
		}
	case "renombrar":
		_, err = fmt.Fprintf(logfile, "RenombrarBase %s %s %s [%s]\n", sector, base, nuevonombre, horaFormateada)
		if err != nil {
			return &pb.Respuesta{Mensaje: "No pudo escribirse correctamente en archivo log", Exitoso: false}, err
		}
	case "actualizar":
		_, err = fmt.Fprintf(logfile, "ActualizarBase %s %s %.0f [%s]\n", sector, base, nuevovalor, horaFormateada)
		if err != nil {
			return &pb.Respuesta{Mensaje: "No pudo escribirse correctamente en archivo log", Exitoso: false}, err
		}
	case "borrar":
		_, err = fmt.Fprintf(logfile, "BorrarBase %s %s [%s] \n", sector, base, horaFormateada)
		if err != nil {
			return &pb.Respuesta{Mensaje: "No pudo escribirse correctamente en archivo log", Exitoso: false}, err
		}
	}
	return &pb.Respuesta{Mensaje: "No pudo escribirse correctamente en archivo log", Exitoso: false}, err
}

func (s *baseServiceServer) AgregarBase(ctx context.Context, req *pb.AgregarBaseRequest) (*pb.Respuesta, error) {
	nombreArchivo := fmt.Sprintf("%s.txt", req.NombreSector)
	if _, err := os.Stat(nombreArchivo); os.IsNotExist(err) {
		// El archivo no existe, entonces se crea uno nuevo
		if err := s.CrearRegistro(nombreArchivo); err != nil {
			return &pb.Respuesta{Mensaje: "Archivo de Sector no pudo ser creado", Exitoso: false}, err
		}
		file, err := os.OpenFile(nombreArchivo, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return &pb.Respuesta{Mensaje: "Archivo de Sector no pudo abrirse exitosamente", Exitoso: false}, err
		}
		defer file.Close()
		// Definir reloj
		_, err = fmt.Fprintf(file, "[0,0,1]\n")
		if err != nil {
			return &pb.Respuesta{Mensaje: "No pudo escribirse correctamente en archivo de sector", Exitoso: false}, err
		}
		// Escribir la información de la base en el archivo
		_, err = fmt.Fprintf(file, "%s %s %.0f", req.NombreSector, req.NombreBase, req.Valor)
		if err != nil {
			return &pb.Respuesta{Mensaje: "No pudo escribirse correctamente en archivo de sector", Exitoso: false}, err
		}
		// función para manipular el log
		_, err = s.AgregarLOG("agregar", req.NombreSector, req.NombreBase, "", req.Valor, 0)
		if err != nil {
			return &pb.Respuesta{Mensaje: "No pudo escribirse correctamente en archivo log", Exitoso: false}, err
		}
		// Obtener reloj para retornarlo al informante
		reloj, err := retornarReloj(nombreArchivo)
		if err != nil {
			return &pb.Respuesta{Mensaje: "Error al extraer el reloj", Exitoso: false}, err
		}
		return &pb.Respuesta{Mensaje: reloj, Exitoso: true}, nil
	}
	// Abrir archivo de sector
	file, err := os.OpenFile(nombreArchivo, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return &pb.Respuesta{Mensaje: "Archivo de Sector no pudo abrirse exitosamente", Exitoso: false}, err
	}
	defer file.Close()

	err = editarReloj(nombreArchivo)
	if err != nil {
		fmt.Println("Error al editar reloj", err)
	}

	// Escribir la información de la base en el archivo
	_, err = fmt.Fprintf(file, "%s %s %.0f", req.NombreSector, req.NombreBase, req.Valor)
	if err != nil {
		return &pb.Respuesta{Mensaje: "No pudo escribirse correctamente en archivo de sector", Exitoso: false}, err
	}

	logfile, err := os.OpenFile("Registro.txt", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return &pb.Respuesta{Mensaje: "Log de registro no pudo abrirse exitosamente", Exitoso: false}, err
	}
	defer logfile.Close()

	// Escribir la información de la base en el log
	_, err = fmt.Fprintf(logfile, "AgregarBase %s %s %.0f\n", req.NombreSector, req.NombreBase, req.Valor)
	if err != nil {
		return &pb.Respuesta{Mensaje: "No pudo escribirse correctamente en archivo log", Exitoso: false}, err
	}
	// Obtener reloj para retornarlo al informante
	reloj, err := retornarReloj(nombreArchivo)
	if err != nil {
		return &pb.Respuesta{Mensaje: "Error al extraer el reloj", Exitoso: false}, err
	}
	return &pb.Respuesta{Mensaje: reloj, Exitoso: true}, nil
}

func (s *baseServiceServer) RenombrarBase(ctx context.Context, req *pb.RenombrarBaseRequest) (*pb.Respuesta, error) {
	nombreArchivo := fmt.Sprintf("%s.txt", req.NombreSector)
	if _, err := os.Stat(nombreArchivo); os.IsNotExist(err) {
		//Hay que agregar lo de los logs
		// El archivo no existe, entonces se crea uno nuevo
		if err := s.CrearRegistro(nombreArchivo); err != nil {
			return &pb.Respuesta{Mensaje: "Archivo de Sector no pudo ser creado", Exitoso: false}, err
		}
		file, err := os.OpenFile(nombreArchivo, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return &pb.Respuesta{Mensaje: "Archivo de Sector no pudo abrirse exitosamente", Exitoso: false}, err
		}
		defer file.Close()

		// Definir reloj
		_, err = fmt.Fprintf(file, "[0,0,1]\n")
		if err != nil {
			return &pb.Respuesta{Mensaje: "No pudo escribirse correctamente en archivo de sector", Exitoso: false}, err
		}
		// Escribir la nueva línea en el archivo
		nuevaLinea := fmt.Sprintf("%s %s 0", req.NombreSector, req.NuevoNombre)
		_, err = file.WriteString(nuevaLinea)
		if err != nil {
			log.Fatalf("Error al escribir en el archivo %s: %v", nombreArchivo, err)
			return &pb.Respuesta{Mensaje: "Comando RenombrarBase no pudo ser ejecutado", Exitoso: false}, err
		}
	} else {
		contenido, err := ioutil.ReadFile(nombreArchivo)
		if err != nil {
			return &pb.Respuesta{Mensaje: "Comando RenombrarBase no pudo ser ejecutado", Exitoso: false}, err
		}
		lineas := strings.Split(string(contenido), "\n")
		encontrar := false
		for i, linea := range lineas {
			if strings.Contains(linea, req.NombreBase) {
				lineas[i] = strings.Replace(linea, req.NombreBase, req.NuevoNombre, -1)
				encontrar = true
				break
			}
		}
		if encontrar == false { //No se encontró la base en el sector
			nuevalinea := fmt.Sprintf("%s %s 0", req.NombreSector, req.NuevoNombre)
			lineas = append(lineas, nuevalinea)
		}
		nuevoContenido := strings.Join(lineas, "\n")
		err = ioutil.WriteFile(nombreArchivo, []byte(nuevoContenido), 0644)
		if err != nil {
			return &pb.Respuesta{Mensaje: "Comando RenombrarBase no pudo ser ejecutado", Exitoso: false}, err
		}
		// Editar reloj
		err = editarReloj(nombreArchivo)
		if err != nil {
			fmt.Println("Error al editar reloj", err)
		}
	}
	// función para manipular el registro
	_, err := s.AgregarLOG("renombrar", req.NombreSector, req.NombreBase, req.NuevoNombre, 0, 0)
	if err != nil {
		return &pb.Respuesta{Mensaje: "No pudo escribirse correctamente en archivo log", Exitoso: false}, err
	}
	// Obtener reloj para retornarlo al informante
	reloj, err := retornarReloj(nombreArchivo)
	if err != nil {
		return &pb.Respuesta{Mensaje: "Error al extraer el reloj", Exitoso: false}, err
	}
	return &pb.Respuesta{Mensaje: reloj, Exitoso: true}, nil
}

func (s *baseServiceServer) ActualizarValor(ctx context.Context, req *pb.ActualizarValorRequest) (*pb.Respuesta, error) {
	nombreArchivo := fmt.Sprintf("%s.txt", req.NombreSector)
	if _, err := os.Stat(nombreArchivo); os.IsNotExist(err) {
		//Hay que agregar lo de los logs
		if err := s.CrearRegistro(nombreArchivo); err != nil {
			return &pb.Respuesta{Mensaje: "Archivo de Sector no pudo ser creado", Exitoso: false}, err
		}
		file, err := os.OpenFile(nombreArchivo, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return &pb.Respuesta{Mensaje: "Archivo de Sector no pudo abrirse exitosamente", Exitoso: false}, err
		}
		defer file.Close()

		// Definir reloj
		_, err = fmt.Fprintf(file, "[0,0,1]\n")
		if err != nil {
			return &pb.Respuesta{Mensaje: "No pudo escribirse correctamente en archivo de sector", Exitoso: false}, err
		}
		// Escribir la nueva línea en el archivo
		nuevaLinea := fmt.Sprintf("%s %s %.0f", req.NombreSector, req.NombreBase, req.NuevoValor)
		_, err = file.WriteString(nuevaLinea)
		if err != nil {
			log.Fatalf("Error al escribir en el archivo %s: %v", nombreArchivo, err)
			return &pb.Respuesta{Mensaje: "Comando ActualizaValor no pudo ser ejecutado", Exitoso: false}, err
		}
	} else {
		contenido, err := ioutil.ReadFile(nombreArchivo)
		if err != nil {
			return &pb.Respuesta{Mensaje: "Comando ActualizaValor no pudo ser ejecutado", Exitoso: false}, err
		}
		encontrar := false
		lineas := strings.Split(string(contenido), "\n")
		for i, linea := range lineas {
			if strings.Contains(linea, req.NombreBase) {
				nuevaLinea := fmt.Sprintf("%s %s %.0f", req.NombreSector, req.NombreBase, req.NuevoValor)
				lineas[i] = nuevaLinea
				encontrar = true
				break
			}
		}
		if encontrar == false {
			nuevalinea := fmt.Sprintf("%s %s %.0f", req.NombreSector, req.NombreBase, req.NuevoValor)
			lineas = append(lineas, nuevalinea)
		}
		nuevoContenido := strings.Join(lineas, "\n")

		err = ioutil.WriteFile(nombreArchivo, []byte(nuevoContenido), 0644)
		if err != nil {
			return &pb.Respuesta{Mensaje: "Comando ActualizaValor no pudo ser ejecutado", Exitoso: false}, err
		}
		// Editar Reloj
		err = editarReloj(nombreArchivo)
		if err != nil {
			fmt.Println("Error al editar reloj", err)
		}
	}
	// función para manipular el log
	_, err := s.AgregarLOG("actualizar", req.NombreSector, req.NombreBase, "", 0, req.NuevoValor)
	if err != nil {
		return &pb.Respuesta{Mensaje: "No pudo escribirse correctamente en archivo log", Exitoso: false}, err
	}
	// Obtener reloj para retornarlo al informante
	reloj, err := retornarReloj(nombreArchivo)
	if err != nil {
		return &pb.Respuesta{Mensaje: "Error al extraer el reloj", Exitoso: false}, err
	}
	return &pb.Respuesta{Mensaje: reloj, Exitoso: true}, nil
}

func (s *baseServiceServer) BorrarBase(ctx context.Context, req *pb.BorrarBaseRequest) (*pb.Respuesta, error) {
	nombreArchivo := fmt.Sprintf("%s.txt", req.NombreSector)
	if _, err := os.Stat(nombreArchivo); os.IsNotExist(err) {
		// El archivo no existe, entonces se crea uno nuevo
		if err := s.CrearRegistro(nombreArchivo); err != nil {
			return &pb.Respuesta{Mensaje: "Archivo de Sector no pudo ser creado", Exitoso: false}, err
		}
		file, err := os.OpenFile(nombreArchivo, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return &pb.Respuesta{Mensaje: "Archivo de Sector no pudo abrirse exitosamente", Exitoso: false}, err
		}
		defer file.Close()

		// Definir reloj
		_, err = fmt.Fprintf(file, "[0,0,1]\n")
		if err != nil {
			return &pb.Respuesta{Mensaje: "No pudo escribirse correctamente en archivo de sector", Exitoso: false}, err
		}
	} else {
		contenido, err := ioutil.ReadFile(nombreArchivo)
		if err != nil {
			return &pb.Respuesta{Mensaje: "Comando BorrarBase no pudo ser ejecutado", Exitoso: false}, err
		}
		lineas := strings.Split(string(contenido), "\n")
		var nuevasLineas []string
		encontrar := false
		for _, linea := range lineas {
			if !strings.Contains(linea, req.NombreBase) {
				nuevasLineas = append(nuevasLineas, linea)
				encontrar = true
			}
		}
		if encontrar == false {
			//implementar código de registro
		}
		nuevoContenido := strings.Join(nuevasLineas, "\n")

		err = ioutil.WriteFile(nombreArchivo, []byte(nuevoContenido), 0644)
		if err != nil {
			return &pb.Respuesta{Mensaje: "Comando BorrarBase no pudo ser ejecutado", Exitoso: false}, err
		}
		// Editar reloj
		err = editarReloj(nombreArchivo)
		if err != nil {
			fmt.Println("Error al editar reloj", err)
		}
	}
	// función para manipular el log
	_, err := s.AgregarLOG("borrar", req.NombreSector, req.NombreBase, "", 0, 0)
	if err != nil {
		return &pb.Respuesta{Mensaje: "No pudo escribirse correctamente en archivo log", Exitoso: false}, err
	}
	// Obtener reloj para retornarlo al informante
	reloj, err := retornarReloj(nombreArchivo)
	if err != nil {
		return &pb.Respuesta{Mensaje: "Error al extraer el reloj", Exitoso: false}, err
	}
	return &pb.Respuesta{Mensaje: reloj, Exitoso: true}, nil
}

func (s *baseServiceServer) GetSoldados(ctx context.Context, req *pb.GetSoldadosRequest) (*pb.Respuesta, error) {
	// Nombre del archivo del sector
	nombreArchivo := fmt.Sprintf("%s.txt", req.NombreSector)

	// Leer el contenido del archivo
	data, err := ioutil.ReadFile(nombreArchivo)
	if err != nil {
		return &pb.Respuesta{
			Mensaje: "Sector no encontrado en Fulcrum 1", Exitoso: false,
		}, nil
	}

	// Convertir el contenido del archivo a líneas
	lineas := strings.Split(string(data), "\n")
	reloj := lineas[0]

	// Buscar la base en las líneas del archivo
	for _, linea := range lineas {
		elementos := strings.Fields(linea)
		if len(elementos) >= 3 && elementos[0] == req.NombreSector && elementos[1] == req.NombreBase {
			// Encontramos la base, devolver la cantidad de soldados
			contenido := fmt.Sprintf("%s - %s - fulcrum3", elementos[2], reloj)
			return &pb.Respuesta{
				Mensaje: contenido, Exitoso: true,
			}, nil
		}
	}

	// La base no fue encontrada
	return &pb.Respuesta{
		Mensaje: "Base no encontrada en comando GetSoldados", Exitoso: true,
	}, nil
}

func main() {
	listener, err := net.Listen("tcp", ":50054")
	if err != nil {
		log.Fatalf("Error al escuchar: %v", err)
	}

	server := grpc.NewServer()
	pb.RegisterMiServicioServer(server, &baseServiceServer{})

	fmt.Println("Servidor gRPC iniciado en el puerto 50054")

	if err := server.Serve(listener); err != nil {
		log.Fatalf("Error al servir: %v", err)
	}
}
