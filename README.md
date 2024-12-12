# Schotky: Scalable URL Shortener Service

Schotky is a high-performance and scalable URL shortener designed to transform long URLs into compact, shareable links. With its robust infrastructure, efficient design, and built-in analytics service, Schotky ensures quick URL resolution, high reliability, and detailed insights into user interactions for an optimal user experience.

---

## 🚀 **Key Features**

- **Scalability & Performance**: Handles **3.5 trillion URLs** and processes **thousands of requests per second** with low latency.  
- **Counter Management**: Uses **ZooKeeper** for atomic updates and **Base62 Encoding** for compact URLs.  
- **Analytics**: Tracks user data (OS, browser, device, location) with **Elasticsearch** and **Grafana** for real-time insights.  
- **Rate Limiting**: Implements **IP-based rate limiting (Fixed Window)** with **NGINX(Lua script)** and **Redis(sharded)** reducing network hops and efficient checks. 
- **Modular Design**: Follows the **repository design pattern** for maintainability.

---

## 🛠️ **Tech Stack**

- **Programming Language**: Golang
- **Framework**: Fiber
- **Database**: AWS DynamoDB with DAX for faster read operations
- **Distributed Counter Management**: ZooKeeper
- **Message Queue**: Kafka for data streaming to Elasticsearch
- **Load Balancing**: NGINX with Lua scripting
- **Rate Limiting**: Redis (sharded) for efficient IP-based rate limiting
- **Analytics**: Elasticsearch and Grafana
- **Containerization**: Docker

---

## 📹 Demo Video

Watch the [demo video][Screencast from 12-12-24 03:48:44 PM IST.webm](https://github.com/user-attachments/assets/89171014-3bf1-4017-9229-e1f007fc3267)

to see Schotky in action.

---

## 📈 System Design

Below is the high-level architecture of Schotky:

![System Design Diagram]![Screenshot 2024-12-03 194100](https://github.com/user-attachments/assets/f2974b96-bbd8-4281-8c0d-bb90da870bc7)

---

## 🔧 Setup Instructions

To run the Schotky service on your machine:

1. Clone the repository:
   ```bash
   git clone https://github.com/SubhamMurarka/Schotky.git

2. Run with Docker
```bash
docker-compose up -d --build
