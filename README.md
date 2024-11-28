# go-test-task-medods

> url/ip в .env файле, по дефолту https://172.28.0.2:443 - api, 172.28.0.3:5432 user/pass = postgres/postgres

Шаги для запуска

1) cd deployment
2) docker build -t medods:latest -f Dockerfile.api ..
3) docker compose up -d

# API endpoints

# GET /api/auth/tokens?userGUID=<userGUID>


### response
```json
{"accessToken":"eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzI5MjIzNDUsImlwIjoiMTcyLjI4LjAuMSIsInVzZXJfZ3VpZCI6ImY0YWUyZjMxLTI3YWQtNGI0Yi1hMDg5LThiZjExZWFjMWQxOSJ9.Gx9XD1ArOwiKR3ewx-VaQXr7x0XAJ9FCtA3odTmMFtzVwnLypsmzJxzoyfa2m0hMRON4XE6V8gy5eGLSKqNH7w","refreshToken":"JDJhJDEwJDBjS2ZPZ2FIaE5ldy5jSmJhVml6Uk9RSU1CTnQuajhXMXZnNFJFUDBVTWR0MWlkSHdkWnJl"}
```
# POST /api/auth/refresh

## request
### BODY JSON

```json
{"refreshToken":"JDJhJDEwJGZCQVZpby42Mi5xWi52MDR2WGVHR09iMTNHTmRhdW9SVXAvVzFKR1VvYy5WUFpSeWkyVVJH"}
```

### response
```json
{"status":"success"}
```

{"refreshToken":"TOKEN"}
# POST /api/auth/test
> тестовая ручка для генерации guid в postgres
### response
```json
{"guid":"ee9234be-51f7-4f49-8f1c-7e1ac2fae902","status":"success"}
```