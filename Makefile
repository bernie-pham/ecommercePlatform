.PHONY: server, createdb, migrateup, migratedown, sqlc, asynqmon

server:
	go run main.go

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root ecommerce_platform

migrateup:
	migrate -path db/migration -database "postgres://root:secret@localhost:5432/ecommerce_platform?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgres://root:secret@localhost:5432/ecommerce_platform?sslmode=disable" -verbose down

sqlc:
	sqlc generate

asynqmon: 
	docker run --rm --name asynqmon \
    	-p 8000:8080 hibiken/asynqmon --redis-addr=host.docker.internal:6379