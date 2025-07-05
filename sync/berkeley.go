package sync

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"solemne3_SO/node" // Cambia por tu nombre real de módulo
)

// BerkeleySync inicia una sincronización desde un nodo coordinador hacia todos los nodos
func BerkeleySync(coordinator *node.Node) {
	fmt.Printf("[%s] Berkeley: Iniciando proceso de sincronizacion como coordinador\n", coordinator.Name)

	var totalDiff time.Duration
	var responses int
	timeDiffs := make(map[string]time.Duration)

	fmt.Printf("[%s] Berkeley: Solicitando hora actual a todos los nodos\n", coordinator.Name)

	// Enviar solicitud de hora a cada nodo
	for _, peer := range coordinator.Peers {
		if peer == coordinator.Address {
			continue // Saltar a sí mismo
		}

		fmt.Printf("[%s] Berkeley: Conectando con nodo %s\n", coordinator.Name, peer)

		conn, err := net.Dial("tcp", peer)
		if err != nil {
			fmt.Printf("[%s] Berkeley: ERROR - No se pudo conectar con %s: %v\n", coordinator.Name, peer, err)
			continue
		}

		// Solicitar hora
		_, err = fmt.Fprintf(conn, "GET_TIME\n")
		if err != nil {
			fmt.Printf("[%s] Berkeley: ERROR - No se pudo enviar solicitud a %s: %v\n", coordinator.Name, peer, err)
			conn.Close()
			continue
		}

		// Recibir respuesta
		message, err := bufio.NewReader(conn).ReadString('\n')
		conn.Close()
		if err != nil {
			fmt.Printf("[%s] Berkeley: ERROR - No se pudo leer respuesta de %s: %v\n", coordinator.Name, peer, err)
			continue
		}

		// Parsear hora
		remoteTime, err := time.Parse("2006-01-02 15:04:05", strings.TrimSpace(message))
		if err != nil {
			fmt.Printf("[%s] Berkeley: ERROR - Formato de hora inválido de %s: %s\n", coordinator.Name, peer, strings.TrimSpace(message))
			continue
		}

		// Calcular diferencia
		diff := remoteTime.Sub(coordinator.GetClock())
		timeDiffs[peer] = diff
		totalDiff += diff
		responses++

		fmt.Printf("[%s] Berkeley: Recibido de %s - Hora: %s, Diferencia: %v\n",
			coordinator.Name, peer, remoteTime.Format("15:04:05"), diff)
	}

	// Agregar la propia hora del coordinador
	timeDiffs[coordinator.Address] = 0
	totalDiff += 0
	responses++

	fmt.Printf("[%s] Berkeley: Hora propia del coordinador: %s\n",
		coordinator.Name, coordinator.GetClock().Format("15:04:05"))

	if responses == 0 {
		fmt.Printf("[%s] Berkeley: ERROR - No se pudo obtener respuesta de ningún nodo\n", coordinator.Name)
		return
	}

	// Calcular promedio de diferencias
	avgDiff := time.Duration(int64(totalDiff) / int64(responses))
	fmt.Printf("[%s] Berkeley: Calculando ajuste promedio basado en %d respuestas\n", coordinator.Name, responses)
	fmt.Printf("[%s] Berkeley: Diferencia promedio calculada: %v\n", coordinator.Name, avgDiff)

	// Enviar ajuste a cada nodo
	fmt.Printf("[%s] Berkeley: Enviando ajustes a todos los nodos\n", coordinator.Name)

	for peer, diff := range timeDiffs {
		adjustment := avgDiff - diff

		if peer == coordinator.Address {
			// Ajustar su propio reloj
			oldTime := coordinator.GetClock()
			newTime := oldTime.Add(adjustment)
			coordinator.SetClock(newTime)
			fmt.Printf("[%s] Berkeley: Ajuste propio - Hora anterior: %s, Hora nueva: %s, Ajuste: %v\n",
				coordinator.Name, oldTime.Format("15:04:05"), newTime.Format("15:04:05"), adjustment)
			continue
		}

		fmt.Printf("[%s] Berkeley: Enviando ajuste a %s: %v\n", coordinator.Name, peer, adjustment)

		conn, err := net.Dial("tcp", peer)
		if err != nil {
			fmt.Printf("[%s] Berkeley: ERROR - No se pudo conectar para enviar ajuste a %s: %v\n", coordinator.Name, peer, err)
			continue
		}

		message := "ADJUST_TIME:" + strconv.FormatInt(int64(adjustment.Seconds()), 10) + "\n"
		_, err = fmt.Fprintf(conn, message)
		if err != nil {
			fmt.Printf("[%s] Berkeley: ERROR - No se pudo enviar ajuste a %s: %v\n", coordinator.Name, peer, err)
		} else {
			fmt.Printf("[%s] Berkeley: Ajuste enviado exitosamente a %s\n", coordinator.Name, peer)
		}
		conn.Close()
	}

	fmt.Printf("[%s] Berkeley: Proceso de sincronización completado\n", coordinator.Name)
}

// HandleBerkeleyMessage interpreta los mensajes relacionados a Berkeley
func HandleBerkeleyMessage(n *node.Node, message string, conn net.Conn) {
	msg := strings.TrimSpace(message)

	switch {
	case msg == "GET_TIME":
		currentTime := n.GetClock().Format("2006-01-02 15:04:05")
		conn.Write([]byte(currentTime + "\n"))
		fmt.Printf("[%s] Berkeley: Solicitud de hora recibida - Enviando: %s\n",
			n.Name, currentTime)

	case strings.HasPrefix(msg, "ADJUST_TIME:"):
		parts := strings.Split(msg, ":")
		if len(parts) != 2 {
			fmt.Printf("[%s] Berkeley: ERROR - Formato de ajuste inválido: %s\n", n.Name, msg)
			return
		}

		adjustmentSec, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			fmt.Printf("[%s] Berkeley: ERROR - Valor de ajuste inválido: %s\n", n.Name, parts[1])
			return
		}

		oldTime := n.GetClock()
		newTime := oldTime.Add(time.Duration(adjustmentSec) * time.Second)
		n.SetClock(newTime)

		fmt.Printf("[%s] Berkeley: Ajuste recibido del coordinador\n", n.Name)
		fmt.Printf("[%s] Berkeley: Hora anterior: %s\n", n.Name, oldTime.Format("15:04:05"))
		fmt.Printf("[%s] Berkeley: Hora nueva: %s\n", n.Name, newTime.Format("15:04:05"))
		fmt.Printf("[%s] Berkeley: Ajuste aplicado: %d segundos\n", n.Name, adjustmentSec)

	default:
		fmt.Printf("[%s] Berkeley: Mensaje no reconocido: %s\n", n.Name, msg)
	}
}
