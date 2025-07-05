package main

import (
	"flag"
	"fmt"
	// "os"
	"solemne3_SO/config"
	"solemne3_SO/node"
	"solemne3_SO/sync"
	"time"
)

func main() {
	// Definir flag --port
	port := flag.String("port", "", "Puerto en el que se iniciará el nodo")
	flag.Parse()

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

	fmt.Printf("[%s] Iniciando en %s...\n", nombreNodo, address)

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
			sync.CristianSync(myNode, peer)
		}
	}

	// Mantener activo
	select {}
}
