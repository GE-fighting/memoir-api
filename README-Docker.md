# Memoir API Docker Setup

This directory contains Docker configuration files to run the Memoir API backend service.

## Files Created

- `Dockerfile` - Multi-stage build for the Go application
- `docker-compose.yml` - Simple setup for the API service only
- `README-Docker.md` - This documentation

## Quick Start

1. **Build and start the API service:**
   ```bash
   docker-compose up --build
   ```

2. **Start in background:**
   ```bash
   docker-compose up -d --build
   ```

3. **View logs:**
   ```bash
   docker-compose logs -f api
   ```

4. **Stop the service:**
   ```bash
   docker-compose down
   ```

## Service

- **API**: Memoir API server (port 5000)

## Environment Configuration

The docker-compose.yml directly uses your existing `.env` file - no additional configuration needed! Just make sure your `.env` file contains the correct settings for your external database and Redis servers.

## Development Workflow

1. **Make code changes**
2. **Rebuild and restart:**
   ```bash
   docker-compose up --build
   ```

3. **View API logs:**
   ```bash
   docker-compose logs -f api
   ```

## Database Management

- **Run migrations manually (if needed):**
  ```bash
  docker-compose exec api ./migrate -action=up
  ```

- **Access the container:**
  ```bash
  docker-compose exec api sh
  ```

## Manual Migration

If you need to run database migrations, you can do it manually:
```bash
# Enter the container
docker-compose exec api sh

# Run migrations
./migrate -action=up
```

## Troubleshooting

- **Check service status:**
  ```bash
  docker-compose ps
  ```

- **View logs:**
  ```bash
  docker-compose logs api
  ```

- **Restart the service:**
  ```bash
  docker-compose restart api
  ```

- **Check connectivity to external services:**
  ```bash
  docker-compose exec api nc -z your-db-host 5432
  docker-compose exec api nc -z your-redis-host 6379
  ```

## Production Notes

1. Ensure your `.env` file has production-ready values
2. Change the JWT secret in production
3. Configure proper CORS origins
4. Add your Aliyun credentials if using cloud storage
5. Make sure your external database and Redis are accessible from the container
6. Consider using Docker secrets for sensitive data

## External Dependencies

This setup assumes you have external services running:
- PostgreSQL database (configured in your `.env` file)
- Redis cache (configured in your `.env` file)

The container will attempt to connect to these services using the configuration from your `.env` file.
