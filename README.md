# 🛒 JOBPORTAL APPLICATION (Gin + PostgreSQL + Bun + AIR + JWT etc)

A clean architecture boilerplate for building large-scale applications using **Gin**, **PostgreSQL**, and **Bun ORM**.  
Supports modular structure with `internal/` folder (DTO, Models, Routes, Utils, etc).

---

## 🚀 Features
- Gin framework for HTTP API
- PostgreSQL as database
- Bun ORM for database operations
- Modular folder structure (`internal/app`, `internal/dto`, `internal/models`, etc.)
- Configurable with `config/`
- Auto reload support with [air](https://github.com/cosmtrek/air)

---

## 📂 Project Structure
```text
ecommerce/
│── cmd/
│   ├── commands/   # CLI commands (make, seed etc.)
│   └── main.go     # Application entrypoint
│
│── config/         # App configuration (config.yaml, db setup)
│── internal/
│   ├── app/        # Services / business logic
│   ├── dto/        # DTOs (request/response)
│   ├── models/     # Database models
│   ├── pkg/        # Reusable packages
│   ├── routes/     # API routes
│   └── utils/      # Helpers (pagination, common utils)
│
│── templates/      # HTML templates
│── migrations/     # Database migrations
│── static/         # Static files
│── tmp/            # Compiled files (ignored in git)
│── .air.toml       # Air configuration for live reload
│── go.mod / go.sum # Go modules

## 🛠️ Requirements
- Go **1.21+**
- PostgreSQL **15+**
- [Air](https://github.com/cosmtrek/air) (for hot reload, optional)

---

## ⚙️ Installation & Setup

### 1️⃣ Clone the repo
```bash
git clone https://github.com/sohidulislam353/gin-postgresql-boilerplat.git
cd your cloning project

###  Open Your terminal or gitbash
```bash
go mod tidy

3️⃣ Configure database

--- Update config/config.yaml with your PostgreSQL credentials:

database:
  host: "localhost"
  port: 5432
  user: "postgres"
  password: "password"
  name: "ecommerce_db"

--- RUN migrations
goose -dir ./migrations postgres "postgresql://your_db_username:your_db_password@localhost:5432/ecommerce?sslmode=disable" up  [run all migrations]

--- Open terminal run the command
air

** Enjoy Go Application **