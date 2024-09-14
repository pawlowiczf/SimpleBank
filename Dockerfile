FROM golang:1.22-alpine3.20 AS builder
WORKDIR /app 
COPY . .
RUN go build -o main main.go 

FROM alpine:3.13 
WORKDIR /app 
COPY --from=builder /app/main . 
COPY --from=builder /app/app.env .

EXPOSE 8080
CMD [ "/app/main" ]