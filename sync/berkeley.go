package sync

import (
    "bufio"
    "fmt"
    "net"
    "strconv"
    "strings"
    "time"

    "github.com/bjaral/solemne3_SO/node" // Cambia por tu nombre real de módulo
)

// BerkeleySync inicia una sincronización desde un nodo coordinador hacia todos los nodos
func BerkeleySync(coordinator *node.Node) {
    fmt.Println("[" + coordinator.Name + "] Iniciando sincronización Berkeley...")

    var totalDiff time.Duration
    var responses int
    timeDiffs := make(map[string]time.Duration)

    // Enviar solicitud de hora a cada nodo
    for _, peer := range coordinator.Peers {
        conn, err := net.Dial("tcp", peer)
        if err != nil {
            fmt.Println("[" + coordinator.Name + "] No se pudo conectar a", peer)
            continue
        }

        // Solicitar hora
        _, err = fmt.Fprintf(conn, "GET_TIME\n")
        if err != nil {
            fmt.Println("[" + coordinator.Name + "] Error enviando solicitud:", err)
            conn.Close()
            continue
        }

        // Recibir respuesta
        message, err := bufio.NewReader(conn).ReadString('\n')
        conn.Close()
        if err != nil {
            fmt.Println("[" + coordinator.Name + "] Error leyendo respuesta:", err)
            continue
        }

        // Parsear hora
        remoteTime, err := time.Parse("2006-01-02 15:04:05", strings.TrimSpace(message))
        if err != nil {
            fmt.Println("[" + coordinator.Name + "] Error parseando hora de", peer)
            continue
        }

        // Calcular diferencia
        diff := remoteTime.Sub(coordinator.GetClock())
        timeDiffs[peer] = diff
        totalDiff += diff
        responses++
    }

    // Agregar la propia hora del coordinador
    timeDiffs[coordinator.Address] = 0
    totalDiff += 0
    responses++

    // Calcular promedio de diferencias
    avgDiff := time.Duration(int64(totalDiff) / int64(responses))
    fmt.Println("[" + coordinator.Name + "] Diferencia promedio:", avgDiff)

    // Enviar ajuste a cada nodo
    for peer, diff := range timeDiffs {
        adjustment := avgDiff - diff
        if peer == coordinator.Address {
            // Ajustar su propio reloj
            newTime := coordinator.GetClock().Add(adjustment)
            coordinator.SetClock(newTime)
            fmt.Println("[" + coordinator.Name + "] Ajustó su reloj a", newTime)
            continue
        }

        conn, err := net.Dial("tcp", peer)
        if err != nil {
            fmt.Println("[" + coordinator.Name + "] No se pudo conectar a", peer)
            continue
        }

        message := "ADJUST_TIME:" + strconv.FormatInt(int64(adjustment.Seconds()), 10) + "\n"
        _, err = fmt.Fprintf(conn, message)
        if err != nil {
            fmt.Println("[" + coordinator.Name + "] Error enviando ajuste:", err)
        }
        conn.Close()
    }
}

// HandleBerkeleyMessage interpreta los mensajes relacionados a Berkeley
func HandleBerkeleyMessage(n *node.Node, message string, conn net.Conn) {
    msg := strings.TrimSpace(message)

    switch {
    case msg == "GET_TIME":
        currentTime := n.GetClock().Format("2006-01-02 15:04:05")
        conn.Write([]byte(currentTime + "\n"))
        fmt.Println("[" + n.Name + "] Enviando hora:", currentTime)

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
        fmt.Println("[" + n.Name + "] Reloj ajustado a", newTime)
    }
}
