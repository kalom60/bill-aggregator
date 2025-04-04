# Billing Aggregator Microservices Platform

A modular microservices architecture designed to support user registration, account linking, utility provider integration, bill aggregation, and inter-service messaging. This project is ideal for simulating a real-world utility bill tracking system.

## 🚀 Microservices Overview

### 1. **user-service**

Handles user registration, authentication, and validation logic.

- **Endpoints** for user register and login operations.
- Stores users in a relational database.
- Uses JWT for auth middleware.

### 2. **account-linking-service**

Allows users to link their accounts to utility providers.

- Stores linked accounts in DB.
- Communicates with `utility-provider-service` and `broker-service` via gRPC.
- Exposes endpoints to manage account links.

### 3. **utility-provider-service**

Registers and manages utility providers (electricity, gas, water, etc).

- Provides gRPC APIs for other services.
- Supports adding and fetching provider details.

### 4. **bill-aggregation-service**

Aggregates bills from linked accounts across providers.

- Calls provider APIs (via mock data in `mock-api-service`).
- Fetches and caches bills using Redis.
- Talks to `account-linking-service` and `utility-provider-service` via gRPC.

### 5. **broker-service**

Acts as a gateway or entry point for external requests.

- Handles incoming API traffic.
- Emits events (e.g., using RabbitMQ or Kafka).
- Contains middleware for auth, rate limiting.

### 6. **listener-service**

Listens to events emitted by the `broker-service`.

- Logs or triggers further async actions.
- Example use: refresh bills.

### 7. **mock-api-service**

Simulates external utility provider APIs.

- Returns mock billing data.
- Used by `bill-aggregation-service` to fetch sample bills.

## 🛠️ Project Structure

```
project/
├── Makefile
├── docker-compose.yml
├── init-scripts/
│   └── 01-init.sh
└── ...microservices
```

## 🧪 How to Run

### Prerequisites

- [Docker](https://www.docker.com/)
- [Make](https://www.gnu.org/software/make/)

### Step-by-Step

1. **Navigate to the project folder:**

   ```bash
   cd project
   ```

2. **Spin up all services using Docker Compose:**
   ```bash
   make up_build
   ```
   This command builds and starts all microservices defined in the `docker-compose.yml`.

### Additional Commands (Optional)

- `make down` — Stop and remove containers.
- `make restart` — Rebuild and restart all services.

---

## 📦 Service Communication

- Services communicate via:
  - **HTTP/REST** (e.g., user APIs)
  - **gRPC** (for fast internal service-to-service communication)
  - **Event-driven architecture** using **RabbitMQ/Kafka** (broker → listener).

## 📂 Database Migrations

Each service has its own migrations folder:

```
internal/database/migrations/
├── ...up.sql
├── ...down.sql
```

You can run them via Go code or integrated tooling (e.g., `golang-migrate`).

