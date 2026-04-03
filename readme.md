# 🚀 AI Task Processor (Backend in Go)

> Not just another CRUD project — this is a backend system built from scratch to understand how things actually work under the hood.

---

## 🧠 Why this project exists

Most tutorials teach:

* “How to build APIs”

Very few teach:

* How backend systems actually work internally

So instead of relying on frameworks and ORMs, this project focuses on:

👉 Understanding core backend fundamentals
👉 Building everything from scratch
👉 Learning by breaking and fixing

---

## ⚙️ Tech Stack

* **Language:** Go (Golang)
* **Database:** PostgreSQL
* **Containerization:** Docker
* **Auth:** JWT (custom implementation)
* **Encryption:** bcrypt
* **Architecture:** Clean layered architecture

---

## 🧱 Architecture

```text
Client
  ↓
Router
  ↓
Middleware (Auth + RBAC)
  ↓
Handler (HTTP layer)
  ↓
Service (Business logic)
  ↓
Repository (DB interaction)
  ↓
PostgreSQL (Docker)
```

---

## 🔐 Features Implemented

### 👤 Authentication System

* User Signup
* User Login
* Password hashing (bcrypt)
* JWT generation & validation

---

### 🧠 Middleware System

* Custom Auth Middleware
* Role-based Authorization (RBAC)
* Context-based user propagation

---

### 🔒 Role-Based Access Control (RBAC)

Roles:

* `USER`
* `ADMIN`
* `MOD` (extendable)

Examples:

* `/admin/users` → ADMIN only
* `/mod/...` → ADMIN + MOD
* `/me` → any authenticated user

---

### 📡 API Endpoints

| Method | Endpoint       | Description                  |
| ------ | -------------- | ---------------------------- |
| POST   | `/createuser`  | Signup                       |
| POST   | `/login`       | Login                        |
| GET    | `/me`          | Get current user             |
| PATCH  | `/update`      | Update user                  |
| DELETE | `/delete`      | Delete user                  |
| GET    | `/admin/users` | Fetch all users (ADMIN only) |

---

### 🧩 Dynamic Update System

* Partial updates using pointer fields
* Only non-nil fields are updated in DB
* Built dynamic SQL query (no ORM)

---

### 📦 JSON Response Standardization

```json
{
  "success": true,
  "message": "Response message",
  "payload": {}
}
```

---

### ⚡ Context Usage

* Passing `user_id` and `role` via request context
* Clean separation between layers
* No global state

---

## 🐳 Database Setup (Docker)

PostgreSQL runs locally via Docker.

---

## 🔐 JWT Design

JWT contains:

```json
{
  "user_id": 1,
  "role": "ADMIN",
  "exp": ...
}
```

👉 No DB calls required for authorization
👉 Fully stateless authentication

---

## 🚨 Security Decisions

* Role is NOT accepted from client
* All users default to `USER`
* Admin role is assigned manually or via protected routes
* Middleware handles authorization (not handlers)

---

## 🎯 Current Focus

This project started as an **AI Task Processor**,
but evolved into a **deep backend learning system**.

---

## 🤖 AI Task Processing (Planned)

Basic/dummy task processing system:

* Submit task → queue simulation
* Process task (mock AI response)
* Store result
* Retrieve processed tasks

👉 Goal is NOT AI itself
👉 Goal is backend system design

---

## 🚧 Upcoming Features

* [ ] Role-based admin endpoints
* [ ] Background workers (goroutines)
* [ ] Task queue simulation
* [ ] Context timeouts
* [ ] Repository interfaces
* [ ] Panic recovery middleware
* [ ] Logging & observability

---

## 💡 Key Learnings

* How middleware actually works
* How JWT authentication flows internally
* How context propagates request-scoped data
* How to structure production-like backend code
* Why avoiding ORMs initially builds deeper understanding

---

## ⚠️ Why No Frameworks?

This project intentionally avoids:

* Gin
* Echo
* Fiber
* ORMs like GORM

👉 To understand:

* Routing
* Middleware chaining
* SQL interactions
* Request lifecycle

---

## 🧠 Philosophy

> Don’t just use tools. Understand them.

---

## 🔗 Getting Started

```bash
git clone <your-repo>
cd ai-task-processor
go run cmd/server/main.go
```

---

## 📌 Final Note

This is not about building features fast.

It’s about building systems that are:

* Scalable ⚙️
* Secure 🔐
* Understandable 🧠

