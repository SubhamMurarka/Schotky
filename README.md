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

![System Design Diagram]![Screenshot 2024-12-03 194100](https://github.com/user-attachments/assets/f2974b96-bbd8-4281-8c0d-bb90da870bc7)


---

## 📹 Demo Video

Watch the [demo video][Screencast from 12-12-24 03:48:44 PM IST.webm](https://github.com/user-attachments/assets/89171014-3bf1-4017-9229-e1f007fc3267)

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
