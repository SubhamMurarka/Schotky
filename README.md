# Schotky: Scalable URL Shortener Service

**Schotky** is a high-performance and scalable URL shortener designed to transform long URLs into compact, shareable links. With its robust infrastructure and efficient design, Schotky ensures quick URL resolution, high reliability, and optimal user experience.

---

## 🚀 Key Features

- **High Scalability**: Supports up to **3.5 trillion unique URLs** and processes **thousands of requests per second**.
- **Efficient URL Resolution**: Optimized for low latency and high performance.
- **Distributed Counter Management**: Implements **ZooKeeper** for atomic counter updates across distributed systems.
- **Base62 Encoding**: Generates collision-free, compact short URLs.
- **Modular Design**: Follows the **repository design pattern** for clean, maintainable code.

---

## 🛠️ Tech Stack

- **Programming Language**: Golang
- **Framework**: Fiber
- **Database**: AWS DynamoDB with DAX for faster read operations
- **Distributed System Management**: ZooKeeper
- **Load Balancing**: NGINX
- **Containerization**: Docker

---

## 📈 System Design

Below is the high-level architecture of Schotky:

![System Design Diagram](https://github.com/user-attachments/assets/7085d73c-6ccc-4f86-ba91-bf32645a1864)

---

## 📹 Demo Video

Watch the [demo video][Screencast from 21-11-24 02:15:48 PM IST.webm](https://github.com/user-attachments/assets/d4a7e8e0-4877-49e1-9773-2fee788310b9)
to see Schotky in action.

---

## 🔧 Setup Instructions

To run the Schotky service on your machine:

1. Clone the repository:
   ```bash
   git clone https://github.com/SubhamMurarka/Schotky.git

2. Move to Schotky/api
   ```bash
   cd Schotky/api

3. Run with Docker
```bash
docker-compose up -d --build
