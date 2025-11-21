.PHONY: swagger-merge migratedown migratenew migrateup

swagger-merge:
	swagger-cli validate ${dir}/bundler.yaml
	swagger-cli bundle ${dir}/bundler.yaml --outfile ${dir}/openapi.yaml --type yaml

swagger-generate:
	go generate internal/generate.go

migratedown:
	dbmate -d ./db/postgres/dbmate/migrations -u "postgres://${DB_USERNAME}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=$(if $(DB_SSL_MODE),$(DB_SSL_MODE),disable)" down

migratenew:
	dbmate -d ./db/postgres/dbmate/migrations -u "postgres://${DB_USERNAME}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=$(if $(DB_SSL_MODE),$(DB_SSL_MODE),disable)" new ${name}

migrateup:
	dbmate -d ./db/postgres/dbmate/migrations -u "postgres://${DB_USERNAME}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=$(if $(DB_SSL_MODE),$(DB_SSL_MODE),disable)" up
