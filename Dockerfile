FROM golang:1.24 AS builder
WORKDIR /app
COPY go.mod go.sum .env ./
COPY ./cmd  ./cmd
COPY ./certs  ./certs
COPY ./pkg ./pkg
COPY ./db ./db
RUN go mod tidy && go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/project /app/cmd/main.go

FROM golang:1.24-alpine AS production
WORKDIR /app
COPY --from=builder /app/bin/project /app/.env ./
COPY --from=builder /app/certs ./certs
COPY --from=builder /app/db ./db
ENTRYPOINT [ "/app/project" ]