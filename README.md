# TaskFlow

TaskFlow is a full-stack task management application built using Go, PostgreSQL, Docker, and React. It provides a structured system for managing projects and tasks with authentication and a modern UI.

---

## Overview

The application allows users to:

* Register and authenticate using JWT
* Create and manage projects
* Create, update, and delete tasks within projects
* Interact with a RESTful API
* Run the entire system using Docker

---

## Tech Stack

### Backend

* Go (Golang)
* Chi Router
* PostgreSQL
* JWT Authentication
* bcrypt for password hashing

### Frontend

* React
* Axios
* React Router
* Tailwind CSS

### DevOps

* Docker
* Docker Compose

---

## Project Structure

```
TaskFlow/
├── backend/
│   ├── cmd/
│   │   └── main.go
│   ├── internal/
│   │   ├── db/
│   │   ├── handlers/
│   │   ├── middleware/
│   │   ├── models/
│   │   └── services/
│   ├── migrations/
│   ├── .env
│   ├── Dockerfile
│   ├── go.mod
│   └── go.sum
│
├── frontend/
│   ├── src/
│   ├── public/
│   ├── package.json
│   └── tailwind.config.js
│
├── docker-compose.yml
└── .gitignore
```

---

## Getting Started

### Prerequisites

* Docker and Docker Compose
* Node.js and npm

---

### Run the Application

#### 1. Start Backend and Database

From the project root:

```
docker-compose up --build
```

This starts:

* Backend server on port 8080
* PostgreSQL database on port 5432

---

#### 2. Start Frontend

In a separate terminal:

```
cd frontend
npm install
npm start
```

---

## Application URLs

* Frontend: http://localhost:3000
* Backend API: http://localhost:8080

---

## API Endpoints

### Authentication

* POST /auth/register
* POST /auth/login

### Projects

* POST /projects
* GET /projects
* GET /projects/{id}
* PATCH /projects/{id}
* DELETE /projects/{id}

### Tasks

* POST /projects/{id}/tasks
* GET /projects/{id}/tasks
* PATCH /tasks/{id}
* DELETE /tasks/{id}

---

## Notes

* JWT-based authentication is implemented on the backend.
* Protected routes require an Authorization header with Bearer token.
* Docker setup ensures consistent environment for backend and database.
* Frontend UI is implemented with a focus on clean layout and usability.

---

## Author

Rohit Kumar
