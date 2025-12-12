FROM golang:1.24 AS builder
WORKDIR /app

COPY . .
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o gophkeeper ./cmd/server


FROM alpine
WORKDIR /app
RUN apk --no-cache add ca-certificates

COPY --from=builder /app/gophkeeper /app/gophkeeper

COPY  database ./database

ENV KEY="secretGophkeeperKey"
ENV ADDRESS="0.0.0.0:8080"
#ENV DBCONNSTR="postgresql://postgres:passw@postgresdb:5432/gophkeeper?sslmode=disable"

EXPOSE 8080
CMD ["./gophkeeper"]