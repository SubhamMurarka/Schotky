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

![Sharded Redis for Rate limiting (2)](https://github.com/user-attachments/assets/08f31c6d-d402-49e3-89b2-bb88a1fe34b3)

### Optimizations

- **Reduced Network Hops**: By directly connecting Redis to Nginx for rate limiting, the need for an additional service is eliminated. This setup reduces the number of network hops by 50%, leading to a significant reduction in latency and improving overall system efficiency.

## Shortening and Redirecting

### Overview

This service handles the shortening of long URLs and redirects users to the corresponding long URLs using a unique short URL.

### How Shortening Works

- **Request Reception**: A user sends a request to shorten a long URL.
  
- **Unique ID Generation**: If the system has exhausted its available unique IDs, it requests a range of IDs from Zookeeper to ensure sequential ID allocation.

- **Base62 Encoding**: The obtained unique ID is encoded using Base62, converting the numeric ID into a compact alphanumeric string, which will serve as the short URL.
  
- **Storing URL**: The short URL along with the corresponding long URL is saved in DynamoDB for persistence.

### How Redirection Works

- **Short URL Request**: A user requests the short URL for redirection.
  
- **URL Lookup**: The service looks up the corresponding long URL for the short URL in DynamoDB Accelerator (DAX).DAX handles cache misses efficiently, ensuring fast retrieval of the long URL.
 
- **Redirection**: Once the long URL is fetched, the user is redirected to the long URL.

![httpswww xyz com (1)](https://github.com/user-attachments/assets/fc05f22c-23b8-4293-8edc-e0565e09de53)

## Why Zookeeper?

// TODO

## Analytics

### Overview

The system collects and processes click event data, stores it for analysis, and visualizes the data in an interactive and insightful manner.

### How Analytics Work

- **Click Event Capture**: Every click event is captured and passed to Kafka for streaming.
  
- **Kafka Consumption**: Kafka consumers read the click events from the stream.
  
- **Data Enrichment**: The consumers batch the data and fetch additional details such as: User IP address, Geolocation, Operating System, Other relevant information from header.
  
- **Data Ingestion**:After enrichment, the batched data is inserted into Elasticsearch for indexing and efficient querying.
  
- **Data Visualization**:Grafana fetches the data from Elasticsearch and presents it in a visual and interactive dashboard, allowing real-time analytics and reporting.
  
### Optimizations:
- **Batching of Data**: The consumers batch the click events before inserting them into Elasticsearch, reducing the number of individual requests and optimizing throughput.Batching improves Elasticsearch indexing performance and reduces network I/O, as Elasticsearch is optimized for bulk operations.

- **ElasticSearch**: Elasticsearch’s indexing and querying capabilities enable efficient retrieval of large volumes of analytics data, ensuring minimal response times for complex queries making it ideal for analytics purpose.

![Short Url](https://github.com/user-attachments/assets/36fd51db-77f4-451d-9cb3-02f28ee5c9ef)

### Grafana Dashboard

![Screenshot 2025-02-07 014028](https://github.com/user-attachments/assets/9a258908-6616-4997-955b-e04ab8955487)

## Calculations

lets consider 2000 req/s

Requests per month = 2000req/s × 60s/min × 60min/h × 24h/day × 30days/month = 5,184,000,000requests/month

lets consider the minimum length of short url be 3 and maximum 7, with base62 encoding.

Total unique IDs = 62^3 + 62^4 + 62^5 + 62^6 + 62^7 = 3583328087528 unique ids can be generated.

Then, total years of serving = Total Unique IDs / Requests per month ≈ 691.07years

## 🛠️ Tech Stack

![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![Fiber](https://img.shields.io/badge/Fiber-FF2D20?style=for-the-badge&logo=fiber&logoColor=white)
![DynamoDB](https://img.shields.io/badge/DynamoDB-4053D6?style=for-the-badge&logo=amazonaws&logoColor=white)
![DAX](https://img.shields.io/badge/DAX-4053D6?style=for-the-badge&logo=amazonaws&logoColor=white)
![Kafka](https://img.shields.io/badge/Apache%20Kafka-231F20?style=for-the-badge&logo=apache-kafka&logoColor=white)
![Nginx](https://img.shields.io/badge/Nginx-009639?style=for-the-badge&logo=nginx&logoColor=white)
![Docker](https://img.shields.io/badge/Docker-2496ED?style=for-the-badge&logo=docker&logoColor=white)
![Docker Compose](https://img.shields.io/badge/Docker_Compose-2496ED?style=for-the-badge&logo=docker&logoColor=white)
![Elasticsearch](https://img.shields.io/badge/Elasticsearch-005571?style=for-the-badge&logo=elasticsearch&logoColor=white)
![Grafana](https://img.shields.io/badge/Grafana-F46800?style=for-the-badge&logo=grafana&logoColor=white)
![Zookeeper](https://img.shields.io/badge/Zookeeper-000000?style=for-the-badge&logo=zookeeper&logoColor=white)

## 🔧 Setup Instructions

To run the Schotky service on your machine:

1. Clone the repository:
   ```bash
   git clone https://github.com/SubhamMurarka/Schotky.git

2. Run with Docker
```bash
docker-compose up -d --build
