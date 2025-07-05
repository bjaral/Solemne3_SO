package main

import (
	"flag"
	"fmt"
	"solemne3_SO/config"
	"solemne3_SO/node"
	"solemne3_SO/sync"
	"time"
)

func main() {
	// ----- Port -----

	// Definir flags
	port := flag.String("port", "", "Puerto en el que se iniciará el nodo")

	// ----- Algorithm -----

	algo := flag.String("algo", "cristian", "Algoritmo de sincronización (cristian|berkeley|logical|vector)")

	flag.Parse()

	// ----- Port -----

	// Si no viene --port, revisar argumentos posicionales
	if *port == "" {
		args := flag.Args()
		if len(args) > 0 {
			*port = args[0]
		} else {
			*port = "8000" // valor por defecto
		}
	}

	address := "localhost:" + *port
	nombreNodo := "Nodo_" + *port

	fmt.Printf("[%s] Iniciando en %s usando algoritmo %s\n", nombreNodo, address, *algo)

	// Crear nodo
	peers := config.NodeAddresses
	myNode := node.NewNode(nombreNodo, address, peers)

	// Iniciar listener en segundo plano
	go myNode.StartListener()

	// Esperar que los nodos estén listos
	time.Sleep(2 * time.Second)

	// Sincronizar con todos los peers menos consigo mismo
	for _, peer := range peers {
		if peer != address {
			fmt.Println("[" + nombreNodo + "] Sincronizando con " + peer)

			// ----- Algorithm -----

			switch *algo {
			case "cristian":
				sync.CristianSync(myNode, peer)
			case "berkeley":
				sync.BerkeleySync(myNode)
			case "logical":
				reloj := sync.NewRelojLogico()
				sync.EnviarMensajeLogico(myNode, peer, reloj, "Hola desde "+nombreNodo)
			case "vector":
				fmt.Println("VectorClock pendiente.")
			default:
				fmt.Println("Algoritmo no reconocido: utilizando algoritmo cristian por defecto", *algo)
				sync.CristianSync(myNode, peer)
			}
		}
	}

	// Mantener activo
	select {}
}
