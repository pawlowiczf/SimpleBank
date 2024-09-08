postgres:
	docker run --name postgres12 --network bank-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=password -d postgres:12-alpine

createdb: 
	docker exec -ti postgres12 createdb --username=root --owner=root simple_bank 

dropdb: 
	docker exec -ti postgres12 dropdb simple_bank

migrateup: 
	migrate -path db/migration -database "postgresql://root:password@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:password@localhost:5432/simple_bank?sslmode=disable" -verbose down

sqlc: 
	sqlc generate

test:
	go test -v -cover ./...

server: 
	go run main.go 