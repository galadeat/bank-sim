# 🏦 bank-sim 

A simulation project demonstrating interaction between two microservices (user and account). It highlights clean architecture, gRPC integration, and a modern client–server setup.

---
## 🚀 Quickstart

Quickly launch both server and client with:
```
make quickstart
```
---
## 📌 REPL UI
The client includes an interactive REPL that simulates account management.  
Through simple menus you can:
- Create and manage users
- Open and manage accounts
- Perform deposits and withdrawals
- Query balances in real time

---
## 🛠 Tech Stack

- **Go 1.24+**
- **gRPC**


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
    ├── cmd/
    │   ├── client
    │   └── server
    ├── internal/
    │   ├── account
    │   ├── repl
    │   └── user
    ├── mocks/
    ├── pkg/
    │   ├── clients
    │   └── logger
    ├── tests/
    │   └── integration
    ├── README.md
    ├── go.mod
    ├── go.sum
    ├── .env.example
    └── .gitignore
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


---
## ✅ Tests
```
make tests
```
---
## 🧠 Features


- **Create** accounts with unique UUIDs  
- **Retrieve** accounts by ID  
- **Store** data in memory using Go maps  
- **Interact** through an intuitive REPL for better UX  
- **Deposit** and **Withdraw** money from accounts  
- **Communicate** via the modern gRPC client API  

## 🔮 Future Plans

This project will evolve into a more realistic banking simulation. Planned features include:

- Persistent storage (e.g., PostgreSQL or MongoDB)

- REST gateway via grpc-gateway

- Authentication and TLS encryption

- Transaction support (deposits, withdrawals, transfers)

- Dockerization and CI/CD pipelines
