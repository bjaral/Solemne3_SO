# Sistemas Distribuidos — Sincronización de Relojes y Comunicación entre Nodos en Go

## Descripción General

Este proyecto es una simulación de un sistema distribuido desarrollado en **Go**, donde múltiples nodos se comunican entre sí, sincronizan sus relojes mediante distintos algoritmos y manejan aspectos de seguridad y tolerancia a fallos.

El objetivo es demostrar de forma práctica cómo se gestionan los procesos de:
- Comunicación entre nodos
- Sincronización de relojes
- Gestión de fallos
- Seguridad en la transmisión de datos

Además, este proyecto forma parte de una presentación académica sobre **Sistemas Distribuidos**.

---

## Contexto Teórico

Los **Sistemas Distribuidos** son un conjunto de computadoras independientes que se presentan ante los usuarios como un sistema único y coherente. En estos entornos, la **sincronización de relojes** y la **comunicación entre procesos** son esenciales para garantizar consistencia, confiabilidad y seguridad.

Este proyecto implementa:
- **4 algoritmos de sincronización de relojes**:
  - Cristian
  - Berkeley
  - Reloj Lógico (Lamport)
  - Reloj Vectorial
- **Tolerancia a fallos simulada** mediante desconexiones y reintentos
- **Seguridad básica** con cifrado de mensajes y autenticación
- **Comunicación entre nodos** vía sockets TCP

---

## Estructura del Proyecto

```bash
.
├── config/
│ ├── config.go # Configuración de direcciones IP, puertos y parámetros
│ └── README.md
├── go.mod # Archivo de dependencias del proyecto Go
├── main.go # Punto de entrada de la aplicación
├── node/
│ ├── node.go # Definición de nodos, sus relojes y métodos de comunicación
│ └── README.md
├── sync/
│ ├── berkeley.go # Algoritmo Berkeley de sincronización
│ ├── cristian.go # Algoritmo Cristian de sincronización
│ ├── logical.go # Reloj lógico de Lamport
│ ├── vector.go # Reloj vectorial
│ └── README.md
├── utils/
│ ├── security.go # Funciones de cifrado, firmas y validación
│ └── README.md
└── README.md # Este archivo
```

---

## Tecnologías Utilizadas

- **Go / Golang** (lenguaje principal)
- **Sockets TCP** (`net` package de Go)
- **Criptografía simétrica / hashing**
- **Algoritmos distribuidos de sincronización**
- **Simulación de fallos controlados**

---

## Funcionamiento General

1. **Configuración de nodos**  
   Se definen direcciones IP y puertos en `config/config.go`.

2. **Inicio de nodos**  
   Cada nodo se lanza como una instancia que escucha conexiones TCP y mantiene un reloj local.

3. **Comunicación**  
   Los nodos pueden enviarse mensajes cifrados, solicitudes de hora y coordinar sincronizaciones.

4. **Sincronización de relojes**  
   Los algoritmos implementados permiten al sistema ajustar sus relojes y mantener coherencia temporal.

5. **Tolerancia a fallos**  
   Simulación de caídas de nodos, reconexión automática y respaldo de estados.

6. **Seguridad**  
   Los mensajes transmitidos entre nodos se cifran y validan para evitar suplantación o manipulación.

---

## Ejecución del Proyecto

### Requisitos
- Go 1.22 o superior
- Puertos TCP libres (por ejemplo: `8000`, `8001`, `8002`)

### Clonar el proyecto

```bash
git clone https://github.com/tu_usuario/tu_repositorio.git
cd tu_repositorio
```

### Compilar y ejecutar el proyecto

```bash
go mod tidy
go run main.go 
```

flags de ejecución:  

* `--port`:  
Puerto TCP donde el nodo escuchará.  
puertos disponibles:  
`8000`, `8001`, `8002`  

* `--algo`:  
Algoritmo de sincronización a usar. Opciones:  
   * `cristian` (por defecto)  
   * `berkeley`  
   * `logical`  
   * `vector`  

ejemplo de uso:

```bash
   # iniciar un nodo en el puerto 8001
   go run main.go --port=8001 --algo=logical


```
