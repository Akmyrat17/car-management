# 📝 Task Manager

This project is a simple **Car Management System** built with Go and Docker, using PostgreSQL for persistence and integrated with observability tools like **Jaeger**, **Prometheus**, and **Grafana**.

---

## ⚙️ Tech Stack

- **Go** — main programming language
- **PostgreSQL** — database
- **Docker** — containerization
- **Jaeger** — distributed tracing
- **Prometheus + Grafana** — monitoring and metrics

---

## 🚀 Getting Started

### 1. 📦 Build & Run

```bash
docker-compose up --build
```

---

## 🏗 Build Info

- Main application Dockerfile: `./Dockerfile`
- PostgreSQL image is built from: `db/Dockerfile`
- Prometheus configuration: `./prometheus.yml`

---

## 🌐 Useful Ports

| Service    | Description            | URL/Port                 |
| ---------- | ---------------------- | ------------------------ |
| App (Go)   | Main HTTP server       | `http://localhost:8080`  |
| PostgreSQL | Database               | `localhost:5424`         |
| Jaeger UI  | Distributed tracing UI | `http://localhost:16686` |
| Prometheus | Metrics collection     | `http://localhost:9090`  |
| Grafana    | Dashboards & metrics   | `http://localhost:3000`  |

> **Grafana Default Credentials:**  
> Username: `admin`  
> Password: `admin`

---

## 📂 Useful Commands

### 📍 Stop Containers

```bash
docker-compose down
```

### 🧹 Clean with Volumes

```bash
docker-compose down -v
```

### ♻️ Rebuild Everything

```bash
docker-compose up --build --force-recreate
```

---

## 📈 Observability

- **Jaeger** – Visualize request traces across the service
- **Prometheus** – Store and aggregate metrics
- **Grafana** – Beautiful dashboards powered by Prometheus

---
