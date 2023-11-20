# Grupo10-Laboratorio-3

* Brandon Monsalve 202073545-k
* Vicente Neira 202004663-8
* Alvaro Soto 202004608-5

## Información Máquinas Virtuales
* dist000.inf.santiago.usm.cl
* MV 129(uUMVDw8hBMnb): Broker Luna
* MV 130(L9uVYbMEdVhj): Fulcrum 1 y La Vanguardia
* MV 131(CjLdwuSkb73v): Fulcrum 2 y Caiatl
* MV 132(64Rb8buuXPeX): Fulcrum 3 y Osiris

## Indicaciones
Al momento de ingresar en la máquina virtual, nosotros operamos de la siguiente forma:
* cd Grupo10-Laboratorio-3(tab) para entrar a la carpeta.
* sudo docker ps -a -> para ver las imagenes y poder empezarlas por ID.
* sudo docker start ID
* sudo docker logs -f ID -> Para ver la consola y los prints.

Comandos importantes para la creación de las imagenes:
* sudo docker stop ID
* sudo docker rm $(sudo docker ps -a -q) -> Borra todas las imagenes
* sudo make docker-broker
* sudo make docker-vanguardia
* sudo make docker-f1
* sudo make docker-f2
* sudo make docker-f3
* sudo make docker-i1
* sudo make docker-i2
