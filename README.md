# Project Time Tracker
Internal tool for tracking time spent on various projects, tasks and other things.

## Prerequisites
- Docker and Docker Compose installed
- Make (for running Makefile commands)

## Environment Setup
Create a `.envrc` file in the project root with the required environment variables:

```bash
# Required
export TURSO_DATABASE_URL=
export TURSO_AUTH_TOKEN=

# Optional
export SERVER_ENV=development
export SERVER_ADDR=9090
export SERVER_READ_TIMEOUT=10s
export SERVER_WRITE_TIMEOUT=30s
export SERVER_IDLE_TIMEOUT=1m
```
## Running the project

### Development mode
Start the development server with hot reloading

```bash
make docker-dev
```

The application will be available at `http://localhost:9090`

## API Endpoints

### Public

 - `POST /v1/auth/register` - Register user
 - `POST /v1/auth/login` - Login

### Authed
 - `GET /v1/me/categories` - List followed categories
 - `PUT /v1/me/categories/{id}/follow` - Follows category
 - `PUT /v1/me/categories/{id}/unfollow` - Unfollows category
 - `POST /v1/me/time_entries` - Make new time entry
 - `DELETE /v1/me/time_entries/{id}` - Delete a time entry
 - `GET /v1/me/time_entries/day/{date}` - Get summary for date (YYYY-MM-DD)
 - `GET /v1/me/time_entries/month/{year-month}` - Get summary for month (YYYY-MM)
