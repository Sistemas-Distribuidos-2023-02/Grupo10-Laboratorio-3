package main

import (
	pb "central/github.com/Sistemas-Distribuidos-2023-02/Grupo10-Laboratorio-3" // Asegúrate de ajustar la importación correctamente
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"

	"google.golang.org/grpc"
)

type baseServiceServer struct {
	pb.UnimplementedMiServicioServer
}

func editarReloj(nombreArchivo, nuevoContenido string) error {
	// Leer el contenido actual del archivo
	contenidos, err := ioutil.ReadFile(nombreArchivo)
	if err != nil {
		return err
	}

	// Separar el contenido por líneas
	lineas := strings.Split(string(contenidos), "\n")

	// Abrir el archivo en modo escritura, truncándolo para eliminar el contenido original
	file, err := os.OpenFile(nombreArchivo, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Escribir el nuevo contenido en el archivo (primera línea)
	_, err = fmt.Fprintln(file, nuevoContenido)
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

var n int

func (s *baseServiceServer) AgregarBase(ctx context.Context, req *pb.AgregarBaseRequest) (*pb.Respuesta, error) {
	nombreArchivo := fmt.Sprintf("Sector%s.txt", req.NombreSector)
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
		n = 1
		// Definir reloj
		_, err = fmt.Fprintf(file, "[%d,0,0]\n", n)
		if err != nil {
			return &pb.Respuesta{Mensaje: "No pudo escribirse correctamente en archivo de sector", Exitoso: false}, err
		}
		// Escribir la información de la base en el archivo
		_, err = fmt.Fprintf(file, "Sector %s %s %.0f", req.NombreSector, req.NombreBase, req.Valor)
		if err != nil {
			return &pb.Respuesta{Mensaje: "No pudo escribirse correctamente en archivo de sector", Exitoso: false}, err
		}
		// Creación de log en caso de no existir
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

		// Escribir la información de la base en el log
		_, err = fmt.Fprintf(logfile, "AgregarBase Sector %s %s %.0f\n", req.NombreSector, req.NombreBase, req.Valor)
		if err != nil {
			return &pb.Respuesta{Mensaje: "No pudo escribirse correctamente en archivo log", Exitoso: false}, err
		}
		return &pb.Respuesta{Mensaje: "Comando AgregarBase ejecutado", Exitoso: true}, nil
	}
	// Abrir archivo de sector
	file, err := os.OpenFile(nombreArchivo, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return &pb.Respuesta{Mensaje: "Archivo de Sector no pudo abrirse exitosamente", Exitoso: false}, err
	}
	defer file.Close()

	n += 1
	nuevoContenido := fmt.Sprintf("[%d,0,0]", n)
	err = editarReloj(nombreArchivo, nuevoContenido)
	if err != nil {
		fmt.Println("Error al editar reloj", err)
	} else {
	}

	// Escribir la información de la base en el archivo
	_, err = fmt.Fprintf(file, "Sector %s %s %.0f", req.NombreSector, req.NombreBase, req.Valor)
	if err != nil {
		return &pb.Respuesta{Mensaje: "No pudo escribirse correctamente en archivo de sector", Exitoso: false}, err
	}

	logfile, err := os.OpenFile("Registro.txt", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return &pb.Respuesta{Mensaje: "Log de registro no pudo abrirse exitosamente", Exitoso: false}, err
	}
	defer logfile.Close()

	// Escribir la información de la base en el log
	_, err = fmt.Fprintf(logfile, "AgregarBase Sector %s %s %.0f\n", req.NombreSector, req.NombreBase, req.Valor)
	if err != nil {
		return &pb.Respuesta{Mensaje: "No pudo escribirse correctamente en archivo log", Exitoso: false}, err
	}
	return &pb.Respuesta{Mensaje: "Comando AgregarBase ejecutado", Exitoso: true}, nil
}

func (s *baseServiceServer) RenombrarBase(ctx context.Context, req *pb.RenombrarBaseRequest) (*pb.Respuesta, error) {
	nombreArchivo := fmt.Sprintf("Sector%s.txt", req.NombreSector)
	if _, err := os.Stat(nombreArchivo); os.IsNotExist(err) {
		//Hay que agregar lo de los logs
		// Crear una nueva línea con el formato deseado
		nuevaLinea := fmt.Sprintf("Sector %s %s 0", req.NombreSector, req.NuevoNombre)

		// Escribir la nueva línea en el archivo
		err = ioutil.WriteFile(nombreArchivo, []byte(nuevaLinea), 0644)
		if err != nil {
			fmt.Printf("Error al crear el archivo %s: %v\n", nombreArchivo, err)
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
			nuevalinea := fmt.Sprintf("Sector %s %s 0", req.NombreSector, req.NuevoNombre)
			lineas = append(lineas, nuevalinea)
		}
		nuevoContenido := strings.Join(lineas, "\n")
		err = ioutil.WriteFile(nombreArchivo, []byte(nuevoContenido), 0644)
		if err != nil {
			return &pb.Respuesta{Mensaje: "Comando RenombrarBase no pudo ser ejecutado", Exitoso: false}, err
		}
	}
	// Escribir la información de la base en el log
	if _, err := os.Stat("Registro.txt"); os.IsNotExist(err) {
		// El log no existe, entonces se crea uno
		if err := s.CrearRegistro("Registro.txt"); err != nil {
			return &pb.Respuesta{Mensaje: "Log de registro no pudo ser creado", Exitoso: false}, err
		}
	}
	logfile, err := os.OpenFile("Registro.txt", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return &pb.Respuesta{Mensaje: "Log de registro no pudo abrirse exitosamente", Exitoso: false}, err
	}
	defer logfile.Close()
	_, err = fmt.Fprintf(logfile, "RenombrarBase Sector %s %s %s\n", req.NombreSector, req.NombreBase, req.NuevoNombre)
	if err != nil {
		return &pb.Respuesta{Mensaje: "No pudo escribirse correctamente en archivo log", Exitoso: false}, err
	}
	return &pb.Respuesta{Mensaje: "Comando RenombrarBase ejecutado", Exitoso: true}, nil
}

func (s *baseServiceServer) ActualizarValor(ctx context.Context, req *pb.ActualizarValorRequest) (*pb.Respuesta, error) {
	nombreArchivo := fmt.Sprintf("Sector%s.txt", req.NombreSector)
	if _, err := os.Stat(nombreArchivo); os.IsNotExist(err) {
		//Hay que agregar lo de los logs
		// Crear una nueva línea con el formato deseado
		nuevaLinea := fmt.Sprintf("Sector %s %s %.0f", req.NombreSector, req.NombreBase, req.NuevoValor)

		// Escribir la nueva línea en el archivo
		err = ioutil.WriteFile(nombreArchivo, []byte(nuevaLinea), 0644)
		if err != nil {
			fmt.Printf("Error al crear el archivo %s: %v", nombreArchivo, err)
			return &pb.Respuesta{Mensaje: "Comando RenombrarBase no pudo ser ejecutado", Exitoso: false}, err
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
				nuevaLinea := fmt.Sprintf("Sector %s %s %.0f", req.NombreSector, req.NombreBase, req.NuevoValor)
				lineas[i] = nuevaLinea
				encontrar = true
				break
			}
		}
		if encontrar == false {
			nuevalinea := fmt.Sprintf("Sector %s %s %.0f", req.NombreSector, req.NombreBase, req.NuevoValor)
			lineas = append(lineas, nuevalinea)
		}
		nuevoContenido := strings.Join(lineas, "\n")

		err = ioutil.WriteFile(nombreArchivo, []byte(nuevoContenido), 0644)
		if err != nil {
			return &pb.Respuesta{Mensaje: "Comando ActualizaValor no pudo ser ejecutado", Exitoso: false}, err
		}
	}
	// Escribir la información de la base en el log
	if _, err := os.Stat("Registro.txt"); os.IsNotExist(err) {
		// El log no existe, entonces se crea uno
		if err := s.CrearRegistro("Registro.txt"); err != nil {
			return &pb.Respuesta{Mensaje: "Log de registro no pudo ser creado", Exitoso: false}, err
		}
	}
	logfile, err := os.OpenFile("Registro.txt", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return &pb.Respuesta{Mensaje: "Log de registro no pudo abrirse exitosamente", Exitoso: false}, err
	}
	defer logfile.Close()
	_, err = fmt.Fprintf(logfile, "ActualizaValor Sector %s %s %.0f\n", req.NombreSector, req.NombreBase, req.NuevoValor)
	if err != nil {
		return &pb.Respuesta{Mensaje: "No pudo escribirse correctamente en archivo log", Exitoso: false}, err
	}
	return &pb.Respuesta{Mensaje: "Comando ActualizaValor ejecutado", Exitoso: true}, nil
}

func (s *baseServiceServer) BorrarBase(ctx context.Context, req *pb.BorrarBaseRequest) (*pb.Respuesta, error) {
	nombreArchivo := fmt.Sprintf("Sector%s.txt", req.NombreSector)
	if _, err := os.Stat(nombreArchivo); os.IsNotExist(err) {
		// El archivo no existe, entonces se crea uno nuevo
		if err := s.CrearRegistro(nombreArchivo); err != nil {
			return &pb.Respuesta{Mensaje: "Comando AgregarBase no pudo ser ejecutado", Exitoso: false}, err
		}
		// Hay que implementar logs aquí
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
	}
	// Escribir la información de la base en el log
	if _, err := os.Stat("Registro.txt"); os.IsNotExist(err) {
		// El log no existe, entonces se crea uno
		if err := s.CrearRegistro("Registro.txt"); err != nil {
			return &pb.Respuesta{Mensaje: "Log de registro no pudo ser creado", Exitoso: false}, err
		}
	}
	logfile, err := os.OpenFile("Registro.txt", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return &pb.Respuesta{Mensaje: "Log de registro no pudo abrirse exitosamente", Exitoso: false}, err
	}
	defer logfile.Close()
	_, err = fmt.Fprintf(logfile, "BorrarBase Sector %s %s\n", req.NombreSector, req.NombreBase)
	if err != nil {
		return &pb.Respuesta{Mensaje: "No pudo escribirse correctamente en archivo log", Exitoso: false}, err
	}
	return &pb.Respuesta{Mensaje: "Comando BorrarBase ejecutado", Exitoso: true}, nil
}

func (s *baseServiceServer) GetSoldados(ctx context.Context, req *pb.GetSoldadosRequest) (*pb.Respuesta, error) {
	// Nombre del archivo del sector
	nombreArchivo := fmt.Sprintf("Sector%s.txt", req.NombreSector)

	// Leer el contenido del archivo
	data, err := ioutil.ReadFile(nombreArchivo)
	if err != nil {
		return &pb.Respuesta{
			Mensaje: "Sector no encontrado en Fulcrum 1", Exitoso: false,
		}, nil
	}

	// Convertir el contenido del archivo a líneas
	lineas := strings.Split(string(data), "\n")

	// Buscar la base en las líneas del archivo
	for _, linea := range lineas {
		elementos := strings.Fields(linea)
		if len(elementos) >= 4 && elementos[0] == "Sector" && elementos[1] == req.NombreSector && elementos[2] == req.NombreBase {
			// Encontramos la base, devolver la cantidad de soldados
			return &pb.Respuesta{
				Mensaje: elementos[3], Exitoso: true,
			}, nil
		}
	}

	// La base no fue encontrada
	return &pb.Respuesta{
		Mensaje: "Base no encontrada en comando GetSoldados", Exitoso: true,
	}, nil
}

func main() {
	listener, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Error al escuchar: %v", err)
	}

	server := grpc.NewServer()
	pb.RegisterMiServicioServer(server, &baseServiceServer{})

	fmt.Println("Servidor gRPC iniciado en el puerto 50052")

	if err := server.Serve(listener); err != nil {
		log.Fatalf("Error al servir: %v", err)
	}
}
