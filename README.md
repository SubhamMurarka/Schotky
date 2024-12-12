# Schotky: Scalable URL Shortener Service

Schotky is a high-performance and scalable URL shortener designed to transform long URLs into compact, shareable links. With its robust infrastructure, efficient design, and built-in analytics service, Schotky ensures quick URL resolution, high reliability, and detailed insights into user interactions for an optimal user experience.

---

## 🚀 **Key Features**

- **High Scalability**: Handles **up to 3.5 trillion unique URLs** and processes **thousands of requests per second**.
- **Low Latency**: Optimized for high performance and quick response times.
- **Distributed Counter Management**: Utilizes **ZooKeeper** for atomic counter updates across distributed systems.
- **Base62 Encoding**: Ensures collision-free, compact short URLs.
- **Click Analytics**: Tracks detailed metrics such as operating system, browser, device, and location.
- **Real-Time Analytics Dashboard**: Displays click analytics using **Elasticsearch** and **Grafana**.
- **Rate Limiting**: Implements **user IP-based rate limiting** with **Redis sharding**.
- **Network Optimization**: Reduces hops by directly connecting **NGINX** with **Redis** for rate-limiting checks.
- **TTL Management**: Automatically resets rate limit counters using Redis's Time to Live (TTL) feature.
- **Fixed Window Algorithm**: Uses the fixed window approach for rate limiting.
- **Modular Design**: Implements the **repository design pattern** for clean, maintainable code.

---

## 🛠️ **Tech Stack**

- **Programming Language**: Golang
- **Framework**: Fiber
- **Database**: AWS DynamoDB with DAX for faster read operations
- **Distributed System Management**: ZooKeeper
- **Message Queue**: Kafka for data streaming to Elasticsearch
- **Load Balancing**: NGINX with Lua scripting
- **Rate Limiting**: Redis (sharded) for efficient IP-based rate limiting
- **Analytics**: Elasticsearch and Grafana
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
