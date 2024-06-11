# Ticket Management

## Overview

A simple microservices system with the following functionalities:

- User:
  - Create a user with information such as full name, username, gender, date of birth, and password.
  - Each user is created with either manager or employee role.
- Tickets:

  - Users are allowed to create tickets with information such as title, content, and priority. Newly created tickets are in the "Pending" state.

  - Only managers can approve/reject tickets. Approved tickets will be in the "Confirmed" state, while rejected tickets will be in the "Rejected" state.

  - To create and confirm tickets, users must provide a token for authentication. The token is generated based on the username and password.

## Architecture Diagram

![Architecture Diagram](architecture_diagram.png)

The system consists of 3 microservices:

- BFF (Service Backend for Frontend): Communicates with the IAM and TM services to handle client requests.
- IAM (Service Identity and Access Management): Manages user data and authenticates users.
- TM (Service Ticket Management): Manages tickets.

````markdown
## Endpoints

### BFF

#### Create User

- Endpoint: `POST /api/v1/create-user`
- Description: Creates a new user with the specified information.
- Request Body:

```json
{
  "full_name": "string",
  "username": "string",
  "gender": "string", // ("male", "female", "other")
  "birth_date": "YYYY-MM-DD",
  "password": "string",
  "role": "string" // ("manager", "employee")
}
```
````

- Status Codes:
  - 201 Created: User created successfully.
  - 400 Bad Request: Invalid input.
  - 409 Conflict: Username already exists.
  - 500 Internal Server Error: Server error.

#### Get Token

- Endpoint: `POST /api/v1/get-token`
- Description: Generates a token for the user based on the username and password.
- Request Body:

```json
{ "username": "string", "password": "string" }
```

- Response:

```json
{ "token": "JWT" }
```

- Status Codes:
  - 200 OK: Token generated successfully.
  - 401 Unauthorized: Invalid credentials.
  - 500 Internal Server Error: Server error.

#### Create Ticket

- Endpoint: `POST /api/v1/create-ticket`
- Description: Creates a new ticket with the specified information.
- Request Body:

```json
{
  "user_id": "UUID",
  "title": "string",
  "content": "string",
  "priority": "string" // ("low", "medium", "high")
}
```

- Response:

```json
{
  "message": "Ticket created successfully",
  "ticket_id": "UUID"
}
```

- Status Codes:
  - 201 Created: Ticket created successfully.
  - 400 Bad Request: Invalid input.
  - 401 Unauthorized: Missing or invalid token.
  - 500 Internal Server Error: Server error.

#### Approve/Reject Ticket

- Endpoint: `POST /api/v1/update-ticket-status`
- Description: Approves or rejects the specified ticket.
- Request

Headers

: `Authorization: Bearer <token>`

- Request Body:

```json
{
  "ticket_id": "UUID",
  "status": "string" // either "approved" or "rejected"
}
```

- Response:

```json
{
  "message": "Ticket status updated successfully"
}
```

- Status Codes:
  - 200 OK: Ticket status updated successfully.
  - 400 Bad Request: Invalid input.
  - 401 Unauthorized: Missing or invalid token.
  - 403 Forbidden: Only managers can approve/reject tickets.
  - 500 Internal Server Error: Server error.

## Getting Started

### Prerequisites

- Go
- PostgreSQL

### Installation

1. Clone the repository
2. Install the dependencies
3. Set up the environment variables

### Running the Services

Each service can be run individually using the `main.go` file in their respective `cmd` directories.

For example, to run the `iam_service`, navigate to the `iam_service/cmd` directory and run:

```sh
go run main.go
```
