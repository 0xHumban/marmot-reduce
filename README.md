# marmot-reduce

🦫 - A little distributed system (Map Reduce simplified)

Marmot-Reduce is a simplified distributed system inspired by the MapReduce programming model. It allows you to perform distributed computations across multiple clients connected to a server. This project is designed to be educational and fun, providing a hands-on experience with distributed systems and concurrent programming in Go.


##  Features

### Server

- **Interactive menu** to manage operations.
- **Display connected clients** in real time.
- **Ping clients** to check their availability.
- **Proper connection management**, preventing resource leaks.
- **Distributed execution of calculations**:
  - Letter counting in a text.
  - Checking if a number is prime.
  - PI estimation using Monte Carlo algorithm

###  Client

- **Persistent connection**:
  - If the connection is lost, the client automatically attempts to reconnect.
  - It continues until the server sends an exit message.

---
## Installation and Execution

### Prerequisites

- **Go** (1.18 or higher) installed on your machine.

### Running the Project

1. Navigate to the `sources` directory:
   
   ```sh
   cd sources
   ```

2. Start the server:
   
   ```sh
   go run ./ server
   ```

3. Start one or more clients:
   
   ```sh
   go run ./
   ```

---

## How It Works

1. **Start the server**: it waits for client connections.
2. **Clients connect**: each client registers with the server.
3. **The server can send commands to clients**: ping, execute calculations, etc.
4. **Clients respond** and return results.
5. **If a client disconnects**, it automatically attempts to reconnect.
6. **The server can properly close all connections** before shutting down.

---

## Follow the Adventure

All project notes and updates are available in `notes.md`.


---
## License

This project is open-source and licensed under **MIT**.
