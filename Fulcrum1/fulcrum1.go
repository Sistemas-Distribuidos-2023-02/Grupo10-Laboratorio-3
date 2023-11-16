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
	nombreArchivo := fmt.Sprintf("Sector%s.txt", req.NombreSector)
	if _, err := os.Stat(nombreArchivo); os.IsNotExist(err) {
		//Hay que agregar lo de los logs
		// Crear una nueva línea con el formato deseado
		nuevaLinea := fmt.Sprintf("%s %s 0", req.NombreSector, req.NombreBase)

		// Escribir la nueva línea en el archivo
		err = ioutil.WriteFile(nombreArchivo, []byte(nuevaLinea), 0644)
		if err != nil {
			fmt.Printf("Error al crear el archivo %s: %v\n", nombreArchivo, err)
			return &pb.Respuesta{Mensaje: "Comando RenombrarBase no pudo ser ejecutado", Exitoso: false}, err
		}
		return &pb.Respuesta{Mensaje: "Comando RenombrarBase ejecutado", Exitoso: true}, nil
	} else {
		contenido, err := ioutil.ReadFile(nombreArchivo)
		if err != nil {
			return &pb.Respuesta{Mensaje: "Comando RenombrarBase no pudo ser ejecutado", Exitoso: false}, err
		}
		lineas := strings.Split(string(contenido), "\n")
		for i, linea := range lineas {
			if strings.Contains(linea, req.NombreBase) {
				lineas[i] = strings.Replace(linea, req.NombreBase, req.NuevoNombre, -1)
				break
			}
		}
		nuevoContenido := strings.Join(lineas, "\n")
		err = ioutil.WriteFile(nombreArchivo, []byte(nuevoContenido), 0644)
		if err != nil {
			return &pb.Respuesta{Mensaje: "Comando RenombrarBase no pudo ser ejecutado", Exitoso: false}, err
		}
	}
	return &pb.Respuesta{Mensaje: "Comando RenombrarBase ejecutado", Exitoso: true}, nil
}

func (s *baseServiceServer) ActualizarValor(ctx context.Context, req *pb.ActualizarValorRequest) (*pb.Respuesta, error) {
	nombresector := fmt.Sprintf("Sector%s.txt", req.NombreSector)
	if _, err := os.Stat(nombresector); os.IsNotExist(err) {
		//Hay que agregar lo de los logs
		return &pb.Respuesta{Mensaje: "Comando ActualizaValor ejecutado", Exitoso: true}, nil
	} else {
		contenido, err := ioutil.ReadFile(nombresector)
		if err != nil {
			return &pb.Respuesta{Mensaje: "Comando ActualizaValor no pudo ser ejecutado", Exitoso: false}, err
		}
		lineas := strings.Split(string(contenido), "\n")
		for i, linea := range lineas {
			if strings.Contains(linea, req.NombreBase) {
				nuevaLinea := fmt.Sprintf("Sector %s %s %.0f\n", req.NombreSector, req.NombreBase, req.NuevoValor)
				lineas[i] = nuevaLinea
				break
			}
		}
		nuevoContenido := strings.Join(lineas, "\n")

		err = ioutil.WriteFile(nombresector, []byte(nuevoContenido), 0644)
		if err != nil {
			return &pb.Respuesta{Mensaje: "Comando ActualizaValor no pudo ser ejecutado", Exitoso: false}, err
		}
	}
	return &pb.Respuesta{Mensaje: "Comando ActualizaValor ejecutado", Exitoso: true}, nil
}

func (s *baseServiceServer) BorrarBase(ctx context.Context, req *pb.BorrarBaseRequest) (*pb.Respuesta, error) {
	nombresector := fmt.Sprintf("Sector%s.txt", req.NombreSector)
	if _, err := os.Stat(nombresector); os.IsNotExist(err) {
		//Hay que agregar lo de los logs
		return &pb.Respuesta{Mensaje: "Comando BorrarBase ejecutado", Exitoso: true}, nil
	} else {
		contenido, err := ioutil.ReadFile(nombresector)
		if err != nil {
			return &pb.Respuesta{Mensaje: "Comando BorrarBase no pudo ser ejecutado", Exitoso: false}, err
		}
		lineas := strings.Split(string(contenido), "\n")
		var nuevasLineas []string
		for _, linea := range lineas {
			if !strings.Contains(linea, req.NombreBase) {
				nuevasLineas = append(nuevasLineas, linea)
			}
		}
		nuevoContenido := strings.Join(nuevasLineas, "\n")

		err = ioutil.WriteFile(nombresector, []byte(nuevoContenido), 0644)
		if err != nil {
			return &pb.Respuesta{Mensaje: "Comando BorrarBase no pudo ser ejecutado", Exitoso: false}, err
		}
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
