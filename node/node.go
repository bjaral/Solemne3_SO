package node

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// // cristian /sync
// if strings.TrimSpace(message) == "TIME_REQUEST" {
//     // Llama a la función del paquete sync
//     sync.HandleTimeRequest(n, message, conn)
//     return
// }

// // berkeley /sync
// if strings.TrimSpace(message) == "GET_TIME" ||
//    strings.HasPrefix(message, "ADJUST_TIME:") {
//     sync.HandleBerkeleyMessage(n, message, conn)
//     return
// }

// // logical /sync
// if strings.HasPrefix(message, "LAMPORT:") {
//     sync.HandleLamportMessage(n, relojLogico, message)
//     return
// }

// Node representa un nodo dentro del sistema distribuido
type Node struct {
	Name      string     // Nombre del nodo
	Address   string     // Dirección IP:Puerto
	Clock     time.Time  // Reloj local del nodo
	Peers     []string   // Lista de direcciones de otros nodos
	Mutex     sync.Mutex // Para acceso concurrente seguro al reloj
	IsRunning bool       // Estado del nodo
}

// NewNode crea una nueva instancia de nodo
func NewNode(name, address string, peers []string) *Node {
	return &Node{
		Name:      name,
		Address:   address,
		Clock:     time.Now().UTC(),
		Peers:     peers,
		IsRunning: true,
	}
}

// StartListener inicia un servidor TCP que escucha mensajes entrantes
func (n *Node) StartListener() {
	ln, err := net.Listen("tcp", n.Address)
	if err != nil {
		fmt.Println("Error iniciando listener:", err)
		os.Exit(1)
	}
	defer ln.Close()

	fmt.Println("["+n.Name+"] escuchando en", n.Address)

	for n.IsRunning {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error aceptando conexión:", err)
			continue
		}
		go n.handleConnection(conn)
	}
}

// handleConnection procesa cada conexión entrante
func (n *Node) handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	message, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("["+n.Name+"] Error leyendo mensaje:", err)
		return
	}

	message = strings.TrimSpace(message)
	n.HandleMessage(message, conn)
}

// HandleMessage interpreta y responde a un mensaje recibido
func (n *Node) HandleMessage(message string, conn net.Conn) {
	fmt.Println("["+n.Name+"] Mensaje recibido:", message)

	if strings.HasPrefix(message, "SETCLOCK:") {
		newTimeStr := strings.TrimPrefix(message, "SETCLOCK:")
		newTime, err := time.Parse("2006-01-02 15:04:05", newTimeStr)
		if err == nil {
			n.Mutex.Lock()
			n.Clock = newTime
			n.Mutex.Unlock()
			fmt.Println("["+n.Name+"] Reloj ajustado a", newTime)
		}
	}

	if message == "TIME_REQUEST" {
		n.HandleTimeRequest(conn)
		return
	}

	if message == "GET_TIME" || strings.HasPrefix(message, "ADJUST_TIME:") {
		n.HandleBerkeleyMessage(message, conn)
		return
	}
}

// SendMessage envía un mensaje a un nodo remoto
func (n *Node) SendMessage(toAddress, message string) {
	conn, err := net.Dial("tcp", toAddress)
	if err != nil {
		fmt.Println("["+n.Name+"] No se pudo conectar a", toAddress, "-", err)
		return
	}
	defer conn.Close()

	_, err = fmt.Fprintf(conn, message+"\n")
	if err != nil {
		fmt.Println("["+n.Name+"] Error enviando mensaje a", toAddress, "-", err)
	}
}

// BroadcastMessage envía un mensaje a todos los nodos conectados
func (n *Node) BroadcastMessage(message string) {
	for _, peer := range n.Peers {
		go n.SendMessage(peer, message)
	}
}

// Stop detiene el nodo (cierra el servidor)
func (n *Node) Stop() {
	n.IsRunning = false
}

// GetClock obtiene el reloj actual del nodo (con protección de concurrencia)
func (n *Node) GetClock() time.Time {
	n.Mutex.Lock()
	defer n.Mutex.Unlock()
	return n.Clock
}

// SetClock ajusta el reloj del nodo (con protección de concurrencia)
func (n *Node) SetClock(t time.Time) {
	n.Mutex.Lock()
	n.Clock = t
	n.Mutex.Unlock()
}

func (n *Node) HandleTimeRequest(conn net.Conn) {
	currentTime := n.GetClock().Format("2006-01-02 15:04:05")
	conn.Write([]byte(currentTime + "\n"))
	fmt.Println("["+n.Name+"] Hora enviada a cliente:", currentTime)
}

func (n *Node) HandleBerkeleyMessage(message string, conn net.Conn) {
	msg := strings.TrimSpace(message)

	switch {
	case msg == "GET_TIME":
		currentTime := n.GetClock().Format("2006-01-02 15:04:05")
		conn.Write([]byte(currentTime + "\n"))
		fmt.Println("["+n.Name+"] Enviando hora:", currentTime)

	case strings.HasPrefix(msg, "ADJUST_TIME:"):
		parts := strings.Split(msg, ":")
		if len(parts) != 2 {
			return
		}
		adjustmentSec, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			return
		}
		newTime := n.GetClock().Add(time.Duration(adjustmentSec) * time.Second)
		n.SetClock(newTime)
		fmt.Println("["+n.Name+"] Reloj ajustado a", newTime)
	}
}
