# chess-online
full stack chess online using go, websocket, reactjs
## Demo
![Image](https://github.com/user-attachments/assets/d4ab8af5-1527-4faa-8878-bf1d3748964d)

## Tech Stack

- Backend: Go
- Frontend: React
- Database: MySQL
- Realtime: WebSocket
- Container: Docker


## Run with Docker

### 1. Create `.env` file

Before running the project, create a `.env` file in the root directory and add the following configuration:

```env
# Server
SERVER_PORT=8080
SERVER_API_KEY=SERVER_API_KEY

# MySQL
# Do not change this
DB_HOST=mysql
# Do not change this
DB_PORT=3306
DB_NAME=chess_online
DB_ROOT_PASSWORD=123456

# JWT
JWT_SECRET_KEY=YOUR_SECRET_KEY

# Valid time units: s (seconds), m (minutes), h (hours)
# Note: 'd' (days) is NOT supported. Use hours instead (e.g., 24h)

# Access token time to live
ACCESS_TOKEN_TTL=1000h

# Refresh token time to live
REFRESH_TOKEN_TTL=72h
```

### 2. Run docker-compose
```
docker-compose up -d --build
```