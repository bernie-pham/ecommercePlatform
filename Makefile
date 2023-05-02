# Start variable area
elastic_url=http://localhost:9200

# End variable area
.PHONY: server, createdb, migrateup, migratedown, sqlc, asynqmon, protogen, initelastic

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



initelastic_dev:
	docker run \
      --name elasticsearch \
      --net bridge -d -p 9200:9200 \
      -e discovery.type=single-node \
      -e ES_JAVA_OPTS="-Xms1g -Xmx1g"\
      -e xpack.security.enabled=false \
      -it 6b84ffb4a623

elastic_create_ecommerce_tag_setting:
	curl --location --request PUT '$(elastic_url)/ecommerce_tag' \
	--header 'Content-Type: application/json' \
	--data '{ \
		"settings": { \
			"analysis": { \
				"analyzer": { \
					"tag_analyzer": { \
						"filter": [ \
							"lowercase", \
							"asciifolding", \
							"custom_stem", \
							"unique" \
						], \
						"char_filter": [ \
							"tag_special_filter" \
						], \
						"type": "custom", \
						"tokenizer": "standard" \
					} \
				}, \
				"char_filter": { \
				"tag_special_filter": { \
					"pattern": "[|&,'\'']", \
					"type": "pattern_replace", \
					"replacement": " " \
				} \
				}, \
				"filter":{ \
					"custom_stem":{ \
						"type":"stemmer", \
						"language":"english" \
					} \
				} \
			} \
		} \
	}'

elastic_create_ecommerce_product_setting:
	curl --location --request PUT '$(elastic_url)/ecommerce_product' \
	--header 'Content-Type: application/json' \
	--data '{ \
		"settings": { \
			"analysis": { \
				"analyzer": { \
					"product_name_analyzer": { \
						"type":      "custom", \
						"tokenizer": "standard", \
						"char_filter": [ \
							"single_char_filter", \
							"useless_filter" \
						], \
						"filter": [ \
							"lowercase", \
							"asciifolding", \
							"custom_stem", \
							"stop", \
							"unique", \
							"remove_digit" \
						] \
					} \
				}, \
				"char_filter": { \
					"single_char_filter": { \
						"type": "pattern_replace", \
						"pattern": "(\\s[a-zA-Z]{1}\\s|\\s[0-9]{1}\\s|[,'\''\\\/!])", \
						"replacement": " " \
					}, \
					"useless_filter": { \
						"type": "pattern_replace", \
						"pattern": "([(]{0,1}[0-9]{0,}[.'\''\"-+,]{1}[0-9]{0,}[”'\''\"]{0,1}[)'\''\"]{0,1})", \
						"replacement": " " \
					} \
				}, \
				"filter":{ \
					"custom_stem":{ \
						"type":"stemmer", \
						"language":"english" \
					}, \
					"remove_digit": { \
						"type":"keep_types", \
						"types": [ "<NUM>" ], \
						"mode": "exclude" \
					} \
				} \
			} \
		} \
	}'
elastic_update_ecommerce_product_mapping:
	curl --location --request PUT '$(elastic_url)/ecommerce_product/_mapping' \
	--header 'Content-Type: application/json' \
	--data '{  \
		"properties": { \
			"product": { \
				"properties": { \
					"name": { \
						"type": "text", \
						"analyzer": "product_name_analyzer" \
					}, \
					"id": { \
						"type": "text", \
						"index": "false" \
					} \
				} \
			} \
		} \
	}'


elastic_update_ecommerce_tag_mapping:
	curl --location --request PUT '$(elastic_url)/ecommerce_tag/_mapping' \
	--header 'Content-Type: application/json' \
	--data '{  \
		"properties": { \
			"tag": { \
				"properties": { \
					"tag": { \
						"type": "text", \
						"analyzer": "tag_analyzer" \
					}, \
					"product_id": { \
						"type": "text", \
						"index": "false" \
					} \
				} \
			} \
		} \
	}'
initelasticindex:
	$(MAKE) elastic_create_ecommerce_tag_setting
	$(MAKE) elastic_create_ecommerce_product_setting
	$(MAKE) elastic_update_ecommerce_tag_mapping
	$(MAKE) elastic_update_ecommerce_product_mapping
