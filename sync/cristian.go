package sync

import (
	"fmt"
	"net"
	"strings"
	"time"

	"solemne3_SO/node" // Reemplaza con el nombre real de tu módulo
)

// CristianSync permite sincronizar el reloj de un cliente con un servidor
func CristianSync(client *node.Node, serverAddress string) {
	fmt.Printf("[%s] Cristian: Iniciando sincronización con servidor %s\n", client.Name, serverAddress)

	// Obtener hora actual del cliente antes de la sincronización
	initialTime := client.GetClock()
	fmt.Printf("[%s] Cristian: Hora inicial del cliente: %s\n", client.Name, initialTime.Format("15:04:05"))

	// Marca de tiempo antes de enviar la solicitud
	T0 := time.Now()
	fmt.Printf("[%s] Cristian: Enviando solicitud de tiempo al servidor\n", client.Name)

	// Conectar al servidor
	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		fmt.Printf("[%s] Cristian: ERROR - No se pudo conectar al servidor %s: %v\n", client.Name, serverAddress, err)
		return
	}
	defer conn.Close()

	fmt.Printf("[%s] Cristian: Conexión establecida con el servidor\n", client.Name)

	// Enviar solicitud
	_, err = fmt.Fprintf(conn, "TIME_REQUEST\n")
	if err != nil {
		fmt.Printf("[%s] Cristian: ERROR - No se pudo enviar solicitud al servidor: %v\n", client.Name, err)
		return
	}

	// Recibir respuesta
	buffer := make([]byte, 128)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Printf("[%s] Cristian: ERROR - No se pudo leer respuesta del servidor: %v\n", client.Name, err)
		return
	}

	// Marca de tiempo al recibir respuesta
	T1 := time.Now()

	// Procesar respuesta
	reply := strings.TrimSpace(string(buffer[:n]))
	fmt.Printf("[%s] Cristian: Respuesta recibida del servidor: %s\n", client.Name, reply)

	serverTime, err := time.Parse("2006-01-02 15:04:05", reply)
	if err != nil {
		fmt.Printf("[%s] Cristian: ERROR - Formato de hora inválido del servidor: %s\n", client.Name, reply)
		return
	}

	// Calcular retardo estimado
	roundTrip := T1.Sub(T0)
	estimatedLatency := roundTrip / 2

	fmt.Printf("[%s] Cristian: Tiempo de ida y vuelta (RTT): %v\n", client.Name, roundTrip)
	fmt.Printf("[%s] Cristian: Latencia estimada: %v\n", client.Name, estimatedLatency)

	// Calcular tiempo estimado del servidor al momento de recibir la respuesta
	estimatedTime := serverTime.Add(estimatedLatency)

	fmt.Printf("[%s] Cristian: Hora del servidor: %s\n", client.Name, serverTime.Format("15:04:05"))
	fmt.Printf("[%s] Cristian: Hora estimada ajustada por latencia: %s\n", client.Name, estimatedTime.Format("15:04:05"))

	// Calcular diferencia entre relojes
	timeDifference := estimatedTime.Sub(initialTime)
	fmt.Printf("[%s] Cristian: Diferencia entre relojes: %v\n", client.Name, timeDifference)

	// Ajustar reloj del cliente
	client.SetClock(estimatedTime)

	finalTime := client.GetClock()
	fmt.Printf("[%s] Cristian: Sincronización completada exitosamente\n", client.Name)
	fmt.Printf("[%s] Cristian: Hora anterior: %s\n", client.Name, initialTime.Format("15:04:05"))
	fmt.Printf("[%s] Cristian: Hora nueva: %s\n", client.Name, finalTime.Format("15:04:05"))
	fmt.Printf("[%s] Cristian: Ajuste aplicado: %v\n", client.Name, timeDifference)

	// Mostrar resumen de la sincronización
	if timeDifference > 0 {
		fmt.Printf("[%s] Cristian: Reloj adelantado en %v\n", client.Name, timeDifference)
	} else if timeDifference < 0 {
		fmt.Printf("[%s] Cristian: Reloj atrasado en %v\n", client.Name, -timeDifference)
	} else {
		fmt.Printf("[%s] Cristian: Reloj ya estaba sincronizado\n", client.Name)
	}
}

// HandleTimeRequest procesa solicitudes de hora de otros nodos
func HandleTimeRequest(n *node.Node, message string, conn net.Conn) {
	if strings.TrimSpace(message) != "TIME_REQUEST" {
		fmt.Printf("[%s] Cristian: Mensaje no reconocido: %s\n", n.Name, strings.TrimSpace(message))
		return
	}

	fmt.Printf("[%s] Cristian: Solicitud de tiempo recibida de un cliente\n", n.Name)

	currentTime := n.GetClock()
	timeString := currentTime.Format("2006-01-02 15:04:05")

	_, err := conn.Write([]byte(timeString + "\n"))
	if err != nil {
		fmt.Printf("[%s] Cristian: ERROR - No se pudo enviar respuesta al cliente: %v\n", n.Name, err)
		return
	}

	fmt.Printf("[%s] Cristian: Tiempo enviado al cliente: %s\n", n.Name, currentTime.Format("15:04:05"))
	fmt.Printf("[%s] Cristian: Respuesta enviada exitosamente\n", n.Name)
}
