# Go Automation Runner

A containerized Go API for executing predefined operational jobs on a Linux server.

## Features
- REST API
- API key authentication
- Health check endpoint
- Designed for Docker and AWS EC2

## Running locally
```bash
export API_KEY="api_key"
go run cmd/server/main.go
