# Stage 1: Build the frontend
FROM node:20-alpine AS frontend-builder

WORKDIR /frontend

COPY frontend/package*.json ./
RUN npm ci --legacy-peer-deps

COPY frontend/ ./
RUN npm run build

# Stage 2: Build the Go backend
FROM golang:1.26.1-alpine AS backend-builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o messenger-server ./cmd/server/main.go

# Stage 3: Run the application
FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=backend-builder /app/messenger-server .
COPY --from=backend-builder /app/migrations ./migrations
COPY --from=frontend-builder /frontend/dist ./frontend/dist

EXPOSE 8080

CMD ["./messenger-server"]