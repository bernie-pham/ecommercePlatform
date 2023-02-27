.PHONY: server, createdb, migrateup, migratedown, sqlc, asynqmon, protogen

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

protogen:
	rm -f pb/*.go
	protoc --proto_path=proto --go_out=pb  --go_opt=paths=source_relative \
	--go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	proto/*.proto 