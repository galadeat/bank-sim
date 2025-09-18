# 🏦 bank-sim gRPC Service

A minimalist gRPC service written in Go for managing user accounts. This project showcases clean architecture, Protocol Buffers integration, environment-based configuration, and a modern client-server setup.

---

## 🚀 Tech Stack

- **Go 1.24+**
- **gRPC**
- **Protocol Buffers**
- **UUID** (`github.com/gofrs/uuid`)
- **godotenv** — for `.env`-based configuration

---

## 📁 Project Structure
```
└── bankAccount/
    ├── api/
    │   └── proto/
    │       └── account/
    │           └── v1/
    │               ├── account.proto
    │               ├── account.pb.go
    │               └── account_grpc.pb.go
    ├── client/
    │   └── main.go
    ├── server/
    │   ├── internal/
    │   │   └── account_serivce.go
    │   └── main.go
    ├── README.md
    ├── go.mod
    ├── go.sum
    ├── .env.example
    └── .gitignore

```

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

## 🧪 Protobuf Generation

If .pb.go files are missing, generate them using:
```
protoc --go_out=paths=source_relative:. \
       --go-grpc_out=paths=source_relative:. \
       api/proto/account/v1/account.proto
```

## 🖥️ Run the Server
```
go run ./server
```

## 📬 Run the Client
```
go run ./client
```
The client creates a new account and immediately retrieves it by ID. All parameters are loaded from .env


## 📌 Sample Output
```
2025/09/18 13:54:43 Account ID: a53a5586-60df-49ed-9a82-510d6ba77dd5 added successfully
2025/09/18 13:54:43 Account: id:"a53a5586-60df-49ed-9a82-510d6ba77dd5" login:"galadeat" email:"user@example.com"
```

## 🧠 Features

- Create accounts with unique UUIDs

- Retrieve accounts by ID

- In-memory storage using Go maps

- Environment-based configuration

- Uses the modern gRPC API (grpc.NewClient)

## 🔮 Future Plans

This project will evolve into a more realistic banking simulation. Planned features include:

- Persistent storage (e.g., PostgreSQL or MongoDB)

- REST gateway via grpc-gateway

- Authentication and TLS encryption

- Transaction support (deposits, withdrawals, transfers)

- Dockerization and CI/CD pipelines
