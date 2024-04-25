FROM golang:1.22 AS build

RUN apt-get update && apt-get install -y \
    libc6-dev \
    gcc \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server/main.go

FROM alpine:3.16

RUN apk --no-cache add libc6-compat libstdc++

RUN adduser -D appuser

WORKDIR /app

COPY --from=build /app/main .

RUN chown -R appuser:appuser /app

USER appuser

EXPOSE 3010

HEALTHCHECK --interval=30s --timeout=30s --start-period=5s --retries=3 \
CMD ["./main", "--healthcheck"]

# Run the application when the container starts
CMD ["./main"]
