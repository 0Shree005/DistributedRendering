# 🚀 Distributed Rendering (Go + Blender)

A custom, high-performance remote rendering pipeline built in **Golang**. This project allows a "Client" machine (e.g., a workstation) to automatically offload Blender rendering tasks to a "Server" machine (e.g., a high-spec laptop) the moment a file is saved.

## 🧠 The "Why"

I prefer working on my PC, but my laptop has the **RTX 3050** and **Ryzen 7 5800H**. Instead of manually moving files via USB or cloud storage every time I want to see a render, I built this to bridge the gap.

It turns any networked machine into a dedicated render node.

## 🛠️ Key Features

* **Automated File Watcher:** Uses system-level events to detect when you hit `Ctrl + S` in Blender.
* **Stack-Based Middleware:** A custom `stackMiddle.go` logic that "de-bounces" noisy file-system events. It ensures the transfer only triggers once the file is fully written to disk.
* **TCP/IP Networking:** Direct socket communication for low-latency file transfer and command signaling.
* **Concurrent Design:** Utilizes Go’s goroutines and channels to handle multiple events without blocking the main execution thread.

## 🏗️ Architecture

1. **Client (Workstation):** Monitors the `.blend` file. When a save is detected, it gates the event through a stack-based validator and streams the file over TCP.
2. **Middleware (`stackMiddle.go`):** Acts as a locking mechanism to prevent multiple redundant transfers during a single save operation.
3. **Server (Render Node):** Listens for incoming data, reconstructs the `.blend` file, and triggers the Blender CLI (`blender -b -a`) to render the frames.

## 📂 Branch Guide

* **`main`:** The core TCP server/client implementation. Basic string-based communication to establish the handshake.
* **`fileChange`:** The "engine room." Contains the file monitoring logic, the stack-based event gate, and the file-handling structures.

## 🚀 Getting Started

### Prerequisites

* [Go](https://go.dev/) (1.21+)
* [Blender](https://www.blender.org/) (Added to your System PATH)

### Installation

```bash
git clone https://github.com/0Shree005/DistributedRendering.git
cd DistributedRendering

```

### Running the Server (The Laptop)

```bash
go run server.go

```

### Running the Client (The PC)

```bash
go run main.go

```

---

## 📄 License

Distributed under the **MIT License**. See `LICENSE` for more information.
