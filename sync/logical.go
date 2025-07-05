package sync

import (
    "fmt"
    "strconv"
    "strings"

    "solemne3_SO/node"
)

// RelojLógico almacena el valor entero del reloj Lamport
type RelojLógico struct {
    Valor int
}

// NewRelojLogico inicializa el reloj lógico
func NewRelojLogico() *RelojLógico {
    return &RelojLógico{Valor: 0}
}

// Incrementa aumenta el contador local antes de un evento
func (r *RelojLógico) Incrementa() {
    r.Valor++
}

// Sincroniza actualiza el reloj con otro valor recibido
func (r *RelojLógico) Sincroniza(valorRemoto int) {
    if valorRemoto > r.Valor {
        r.Valor = valorRemoto
    }
    r.Valor++
}

// Get retorna el valor actual del reloj
func (r *RelojLógico) Get() int {
    return r.Valor
}

// EnviarMensajeLogico envía un mensaje con el reloj lógico actual
func EnviarMensajeLogico(from *node.Node, to string, reloj *RelojLógico, contenido string) {
    reloj.Incrementa()

    message := fmt.Sprintf("LAMPORT:%d:%s", reloj.Get(), contenido)
    from.SendMessage(to, message)

    fmt.Printf("[%s] Envió mensaje a %s con reloj lógico %d\n", from.Name, to, reloj.Get())
}

// HandleLamportMessage procesa un mensaje con reloj lógico Lamport
func HandleLamportMessage(from *node.Node, reloj *RelojLógico, message string) {
    parts := strings.SplitN(message, ":", 3)
    if len(parts) != 3 {
        return
    }

    valorRemoto, err := strconv.Atoi(parts[1])
    if err != nil {
        return
    }

    reloj.Sincroniza(valorRemoto)

    fmt.Printf("[%s] Recibió mensaje: '%s' con reloj lógico remoto %d, nuevo reloj local: %d\n",
        from.Name, parts[2], valorRemoto, reloj.Get())
}
