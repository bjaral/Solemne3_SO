package main

import (
    "github.com/bjaral/solemne3_SO/config"
    "github.com/bjaral/solemne3_SO/node"
    "github.com/bjaral/solemne3_SO/sync"
    "time"
)

func main() {
    // Crear nodo
    peers := config.NodeAddresses
    myNode := node.NewNode("Nodo1", "localhost:8000", peers)

    // Iniciar listener en segundo plano
    go myNode.StartListener()

    // Esperar que los nodos estén listos
    time.Sleep(2 * time.Second)

    // Sincronizar con un servidor
    sync.CristianSync(myNode, "localhost:8001")

	// sync.BerkeleySync(myNode)
	relojLogico := sync.NewRelojLogico()
	//sync.EnviarMensajeLogico(myNode, "localhost:8001", relojLogico, "Hola desde Nodo1")


    // Mantener activo
    select {}
}
