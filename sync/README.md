# Algoritmos de Sincronización

Esta carpeta contiene las implementaciones de los cuatro algoritmos de sincronización de relojes que se usan en el proyecto.

## Archivos y su función

- `cristian.go`: Implementa el algoritmo Cristian, donde el cliente solicita la hora a un servidor y ajusta su reloj compensando la latencia.
- `berkeley.go`: Implementa el algoritmo Berkeley, donde un nodo maestro calcula el promedio de las horas de los nodos y envía ajustes.
- `logical.go`: Implementa el reloj lógico (Lamport) para mantener el orden de eventos en sistemas distribuidos.
- `vector.go`: Implementa el reloj vectorial para mantener el orden parcial y la causalidad entre eventos.