# 🏦 bank-sim gRPC Service

A minimalist gRPC service written in Go for managing user accounts. This project showcases clean architecture, Protocol Buffers integration, environment-based configuration, and a modern client-server setup.

---
## 🚀 Quickstart
Generate .pb.go files if missing:
```
make pb
```
Quickly launch both server and client with:
```
make quickstart
```
---

## 🛠 Tech Stack

- **Go 1.24+**
- **gRPC**
- **Protocol Buffers**
- **UUID** (`github.com/gofrs/uuid`)
- **godotenv** — for `.env`-based configuration
- **Go test** — built-in testing framework 

---

## 📁 Project Structure
```
.
└── bankAccount/
    ├── api/
    │   └── proto/               # gRPC contracts
    │       ├── account/         # account service (v1, v2)
    │       ├── common/          # common types (Money, etc)
    │       ├── reporting/       # reporting service
    │       ├── transaction/     # transaction service
    │       └── user/            # user service
    ├── client/
    │   └── main.go
    ├── server/
    │   ├── account/
    │   │   └── account_service.go
    │   └── main.go
    ├── tests/
    │   ├── integration/
    │   │   └── client_test.go
    │   └── unit/
    │       └── account_test.go
    ├── README.md
    ├── go.mod
    ├── go.sum
    ├── .env.example
    └── .gitignore
```
---
## ⚙️ Configuration

Before running the project, copy `.env.example` to `.env` and fill in your values:

```
cp .env.example .env
```
Example .env:
```
GRPC_ADDRESS=localhost:50051
LOGIN=your_login_here
EMAIL=your_email_here
```
---

## 🖥️ Server
```
make run-server
```

## 📬 Client
```
make run-client
```
The client creates a new account and immediately retrieves it by ID. All parameters are loaded from .env

---
## ✅ Tests
```
make test
```

## 📌 Sample Output
```
2025/09/18 13:54:43 Account ID: a53a5586-60df-49ed-9a82-510d6ba77dd5 added successfully
2025/09/18 13:54:43 Account: id:"a53a5586-60df-49ed-9a82-510d6ba77dd5" login:"galadeat" email:"user@example.com"
```
---
## 🧠 Features

- Creating accounts with unique UUIDs

- Retrieving accounts by ID

- In-memory storage using Go maps

- Environment-based configuration

- Uses the modern gRPC client API (grpc.NewClient, available in gRPC-Go v1.60+ with Go 1.23+)

## 🔮 Future Plans

This project will evolve into a more realistic banking simulation. Planned features include:

- Persistent storage (e.g., PostgreSQL or MongoDB)

- REST gateway via grpc-gateway

- Authentication and TLS encryption

- Transaction support (deposits, withdrawals, transfers)

- Dockerization and CI/CD pipelines
