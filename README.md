# distributed Tasks Processing Service

A service for handling various tasks. The peculiarity lies in the distributed nature of the system. 

> **THE PROJECT IS EDUCATIONAL**

## Features

- Transactional Outbox pattern
- Idempotent workers
- Message broker integration
- Retry mechanism
- Graceful shutdown

This project demonstrates:

- Transactional Outbox Pattern
- Idempotent Workers
- At-least-once message delivery

## Architecture

The system consists of:

- API service
- Outbox relay
- Message broker
- Worker service (node)

![](/.source/architecture.png)

## Tech Stack

- Go
- PostgreSQL
- Kafka
- Docker

## Project Structure

```
cmd/
api/
worker/
internal/
domain/
service/
repository/
pkg/
configs/
scripts/
```

## 6. Getting Started

### Requirements

- Go 1.22+
- Docker
- PostgreSQL

### Run

```sh
git clone https://github.com/Sapiyulla/dJPS.git
cd dJPS
docker compose up
```

## Configuration

| Variable | Description |
|---------|-------------|
| DB_URL | PostgreSQL connection |

## Roadmap

- User API
- Task API
- Outbox + Relay
- Worker 
