package sync

import (
    "fmt"
    "net"
    "strings"
    "time"

    "solemne3_SO/node"  // Reemplaza con el nombre real de tu módulo
)

// CristianSync permite sincronizar el reloj de un cliente con un servidor
func CristianSync(client *node.Node, serverAddress string) {
    fmt.Println("[" + client.Name + "] Iniciando sincronización con", serverAddress)

    // Marca de tiempo antes de enviar la solicitud
    T0 := time.Now()

    // Conectar al servidor
    conn, err := net.Dial("tcp", serverAddress)
    if err != nil {
        fmt.Println("[" + client.Name + "] Error conectando a servidor:", err)
        return
    }
    defer conn.Close()

    // Enviar solicitud
    _, err = fmt.Fprintf(conn, "TIME_REQUEST\n")
    if err != nil {
        fmt.Println("[" + client.Name + "] Error enviando solicitud:", err)
        return
    }

    // Recibir respuesta
    buffer := make([]byte, 128)
    n, err := conn.Read(buffer)
    if err != nil {
        fmt.Println("[" + client.Name + "] Error leyendo respuesta:", err)
        return
    }

    // Marca de tiempo al recibir respuesta
    T1 := time.Now()

    // Procesar respuesta
    reply := strings.TrimSpace(string(buffer[:n]))
    serverTime, err := time.Parse("2006-01-02 15:04:05", reply)
    if err != nil {
        fmt.Println("[" + client.Name + "] Error parseando hora:", err)
        return
    }

    // Calcular retardo estimado
    roundTrip := T1.Sub(T0)
    estimatedTime := serverTime.Add(roundTrip / 2)

    // Ajustar reloj del cliente
    client.SetClock(estimatedTime)

    fmt.Println("[" + client.Name + "] Sincronización completada. Nuevo reloj:", estimatedTime)
}

// HandleTimeRequest procesa solicitudes de hora de otros nodos
func HandleTimeRequest(n *node.Node, message string, conn net.Conn) {
    if strings.TrimSpace(message) == "TIME_REQUEST" {
        currentTime := n.GetClock().Format("2006-01-02 15:04:05")
        conn.Write([]byte(currentTime + "\n"))
        fmt.Println("[" + n.Name + "] Hora enviada a cliente:", currentTime)
    }
}
