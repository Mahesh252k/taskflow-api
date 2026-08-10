# TaskFlow API

TaskFlow API is a RESTful Task Management API built with **Go**, **Gin**, **GORM**, and **MySQL**.

This project is focused on building a backend API using real-world Go development practices, including layered architecture, database relationships, middleware, pagination, filtering, sorting, and search.

## Features

* Create tasks
* Get all tasks
* Get a task by ID
* Update tasks
* Delete tasks
* MySQL database integration
* GORM ORM
* User and Task relationship
* Pagination
* Filter tasks by completion status
* Filter tasks by user
* Sort tasks
* Search tasks by title and description
* Authentication middleware
* Request logging middleware
* Environment variable configuration
* Database migrations and seed data
* Structured error handling

## Tech Stack

* Go
* Gin
* GORM
* MySQL
* godotenv
* REST API

## Project Structure

```text
taskflow-api/
├── config/
├── database/
├── handlers/
├── middleware/
├── models/
├── routes/
├── services/
├── .env
├── .gitignore
├── go.mod
├── go.sum
└── main.go
```

> The `.env` file is used locally and should not be committed to Git.

## API Endpoints

| Method | Endpoint     | Description       |
| ------ | ------------ | ----------------- |
| GET    | `/tasks`     | Get all tasks     |
| GET    | `/tasks/:id` | Get a task by ID  |
| POST   | `/tasks`     | Create a new task |
| PUT    | `/tasks/:id` | Update a task     |
| DELETE | `/tasks/:id` | Delete a task     |

## Query Parameters

The `GET /tasks` endpoint supports pagination, filtering, sorting, and search.

### Pagination

```text
GET /tasks?page=1&limit=10
```

### Filter by Completion Status

```text
GET /tasks?done=false
```

### Filter by User

```text
GET /tasks?user_id=1
```

### Search

Searches both task title and description.

```text
GET /tasks?search=go
```

### Sorting

```text
GET /tasks?sort=title&order=asc
```

Supported sort fields:

```text
id
title
created_at
```

Supported order values:

```text
asc
desc
```

### Combining Query Parameters

Filters, search, sorting, and pagination can be combined.

```text
GET /tasks?page=1&limit=10&done=false&user_id=1&search=go&sort=created_at&order=desc
```

## Example Task

```json
{
  "title": "Learn Go REST API",
  "description": "Practice building APIs using Gin and GORM",
  "done": false,
  "user_id": 1
}
```

## Environment Variables

Create a `.env` file in the project root.

```env
DB_USER=your_mysql_user
DB_PASSWORD=your_mysql_password
DB_HOST=127.0.0.1
DB_PORT=3306
DB_NAME=taskdb

APP_PORT=8080
APP_NAME=TaskFlow
```

Do not commit your `.env` file.

Add it to `.gitignore`:

```text
.env
```

## Running the Project

Clone the repository:

```bash
git clone <your-repository-url>
cd taskflow-api
```

Install dependencies:

```bash
go mod download
```

Configure your `.env` file and make sure MySQL is running.

Start the API:

```bash
go run main.go
```

The application will run locally on:

```text
http://localhost:8080
```

## Example Request

```text
GET /tasks?page=1&limit=5&search=go&sort=title&order=asc
```

This request:

1. Searches tasks containing `go`
2. Sorts matching tasks by title
3. Uses ascending order
4. Returns the first page
5. Limits the response to 5 tasks

## Purpose

This project is being developed incrementally to teach students like production-style backend development with Go.

The goal is to continue improving the API by adding more real-world backend concepts such as validation, authentication, testing, API documentation, containerization, and deployment.
