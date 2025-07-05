# Nodo (node.go)

Este componente define la estructura principal del nodo dentro del sistema distribuido.

## Qué incluye

- Definición del struct `Node`, que contiene información esencial como:
  - Identificador único del nodo (`ID`)
  - Puerto en el que escucha
  - Reloj local (hora)
  - Lista de pares (otros nodos con los que puede comunicarse)

- Métodos asociados al nodo para manejar su lógica básica, como:
  - Inicialización del nodo
  - Funciones para actualizar el reloj local
  - Manejo básico de comunicación o sincronización (en algunos casos)