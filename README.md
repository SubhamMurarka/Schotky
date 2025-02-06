# Schotky: Scalable URL Shortener Service

## Rate Limiting Service

### Overview

This rate limiting service is built using Nginx OpenResty with Lua scripts and is designed to efficiently rate limit user requests based on their IP addresses using the Fixed Window approach. The rate limiter utilizes Redis for fast and added benefits of TTL (Time-To-Live) and sharding based on the user’s IP address.

## 📹 Demo Video

Watch the [demo video][Screencast from 12-12-24 03:48:44 PM IST.webm](https://github.com/user-attachments/assets/89171014-3bf1-4017-9229-e1f007fc3267)
to see Schotky in action.

### How It Works

- **Redis Shard Lookup**: Based on the IP address of requested user, Nginx determines the Redis shard where the rate limit data should be stored.The hashing mechanism ensures that requests from the same IP always hit the same Redis shard, distributing the load across multiple Redis instances and minimizing the risk of contention.

- **Rate Limit Check & Count Update**: The rate limiting logic implemented in Lua checks if the IP is present in the Redis shard. If the IP is not found, it registers the IP address and 
initializes the request count for that IP. If the IP is found, Nginx forwards the request to the server, where the request count for that IP is incremented.

- **Rate Limit Enforcement**: Once the server processes the request and updates the count, the Lua script checks if the user has exceeded the allowed number of requests within a fixed time window.
If the rate limit is exceeded, Nginx blocks the request and sends a 429 (Too Many Requests) response to the user.If the rate limit is not exceeded, the request proceeds to the destination server as normal.

- **TTL Management in Redis**: Redis ensures the request count is stored with a Time-To-Live (TTL), automatically resetting the count after the defined time window has passed (e.g., 1 minute).
This ensures that once the window expires, the count is reset, and the user can start making requests again.

### Optimizations

- **Reduced Network Hops**: Since the rate limiting decision (IP check and request count increment) is handled in-memory in Redis, the system minimizes network latency. This is a significant optimization because, in traditional systems, each request might require querying a central server for rate limiting, resulting in additional network hops.
  
By having only Redis as an external dependency for rate limiting, the system’s efficiency is significantly improved, reducing the number of network hops by 50% compared to a setup that involves multiple services.

## 📈 System Design

Below is the high-level architecture of Schotky:

![System Design Diagram]![Screenshot 2024-12-03 194100](https://github.com/user-attachments/assets/f2974b96-bbd8-4281-8c0d-bb90da870bc7)

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

## 🔧 Setup Instructions

To run the Schotky service on your machine:

1. Clone the repository:
   ```bash
   git clone https://github.com/SubhamMurarka/Schotky.git

2. Run with Docker
```bash
docker-compose up -d --build
