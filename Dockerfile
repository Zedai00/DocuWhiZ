# ==== FRONTEND BUILD ====
FROM node:20-alpine AS frontend

WORKDIR /app/frontend
COPY frontend/ ./
RUN npm install && npm run build

# ==== BACKEND BUILD ====
FROM golang:1.24-alpine AS backend

WORKDIR /app/backend
COPY backend/ ./
COPY --from=frontend /app/frontend/dist ./dist

RUN go mod tidy
RUN go build -o server .

# ==== FINAL STAGE ====
FROM alpine:latest

# Install pdftotext (poppler-utils) and libc
RUN apk add --no-cache poppler-utils libc6-compat

WORKDIR /app

# Copy server binary
COPY --from=backend /app/backend/server .

# Copy frontend dist folder
COPY --from=backend /app/backend/dist ./dist

ENV PORT=8000
EXPOSE 8000

CMD ["./server"]

