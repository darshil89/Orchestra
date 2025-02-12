# 🎼 Orchestra - Distributed Job Orchestration

## 🚀 Overview

**Orchestra** is a **distributed job orchestration system** designed for **CI/CD pipeline execution** and **data processing** tasks. It efficiently manages task execution across multiple workers using a **message queue**, enabling parallel execution and real-time status updates.

## 🎯 Applications

- **Task Execution & Distribution**: Dynamically schedules and distributes tasks across worker nodes.
- **CI/CD Automation**: Supports running tests, building images, and deploying code.
- **Data Processing**: Handles parallel data processing jobs like ML training, image processing, etc.
- **Message Queues (RabbitMQ)**: Ensures tasks are processed asynchronously.
- **Real-time Status Updates (WebSockets)**: Users get instant updates when tasks are completed.
- **Scalable & Flexible Architecture**: Easily integrates with Docker, Kubernetes, and cloud platforms.

## 🏗️ Architecture

---
![Mail](https://drive.google.com/file/d/1SFGdy9O0Kg2oIHzalqJxqaDt6DfGzKh6/view)
---

1. **Task Producer (Client)**
   - A frontend app (React) or a service sends a task request to the backend.
2. **Backend (Task Manager)**
   - Accepts the task, stores initial data in **Redis**, and pushes the job to **RabbitMQ**.
3. **Message Broker (RabbitMQ)**
   - Queues tasks and ensures workers process them asynchronously.
4. **Worker Nodes**
   - Runs in **Docker/Kubernetes**, pulls tasks from RabbitMQ, and executes them.
5. **Result Storage (Redis)**
   - Stores task progress and final results.
6. **WebSockets (Socket.IO)**
   - Pushes task updates to the frontend in real-time.

## 🛠️ Tech Stack

| Component       | Technology Used  |
|---------------|----------------|
| **Frontend**   | React, Next.js  |
| **Backend**    | Go (Gin + Chi Router) , AMQP |
| **Queue**      | RabbitMQ (pub/sub)  |
| **Database**   | Redis  |
| **Worker Nodes** | Docker, Kubernetes |
| **Real-time Updates** | Gorilla WebSockets (Socket.IO) |

## 🏗️ Installation & Setup

### 1️⃣ Clone the Repository
```sh
git clone https://github.com/yourusername/orchestra.git
cd orchestra
```

### 2️⃣ Start RabbitMQ & Redis (Docker)
```sh
docker run -d --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:management
docker run -d --name redis -p 6379:6379 redis
```

### 3️⃣ Start Backend Server
```sh
go run main.go
```

### 4️⃣ Start Frontend
```sh
cd frontend
npm install
npm run dev
```

### 5️⃣ Run Worker Nodes
```sh
docker-compose up worker
docker run -d --name redis -p 6379:6379 redis
docker run -d --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:management
```

## 📡 WebSocket Integration

To receive real-time updates, the frontend connects to the WebSocket server:
```js
const ws = new WebSocket("ws://localhost:8080/ws");
ws.onmessage = (event) => {
    console.log("Task update:", event.data);
};
```

## 🎬 Usage
1. Open the frontend UI.
2. Click **Start Task** to send a job request.
3. Backend adds it to **RabbitMQ**, and a worker processes it.
4. Get **real-time task updates** via WebSockets.

## 🚀 Future Enhancements
- ✅ Add support for **multiple queues** for different task priorities.
- ✅ Implement **Kubernetes Auto-scaling** for dynamic worker allocation.
- ✅ Provide a **UI Dashboard** to monitor job progress visually.

## 🎉 Contributing
Pull requests are welcome! Please fork the repo and submit a PR. 😊

## 📜 License
MIT License © 2025 Orchestra Project
