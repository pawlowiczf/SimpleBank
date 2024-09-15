FROM golang:1.22-alpine3.20 AS builder
WORKDIR /app 
COPY . .
RUN go build -o main main.go 
RUN apk add curl 
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.1/migrate.linux-amd64.tar.gz | tar xvz

FROM alpine:3.13 
WORKDIR /app 
COPY --from=builder /app/main . 
COPY --from=builder /app/migrate .
COPY db/migration ./migration
COPY app.env .
COPY start.sh .
COPY wait-for.sh .

EXPOSE 8080
# CMD [ "/app/main" ]
# ENTRYPOINT [ "/app/start.sh" ]

