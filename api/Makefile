.PHONY: up down clean migrate migrate-dev drop drop-dev dev-db-pass run rund gen model cm

db-password = password #$(shell gcloud secrets versions access latest --secret="dev-db-password")

up:
		docker-compose up -d

down:
		docker-compose down

clean:
		docker-compose down
		docker volume prune -f
		docker-compose up -d

migrate:
	@echo "up or down?: "; \
	read ACTION; \
	migrate -source file://ent/migrate/migrations -database 'postgres://root:password@localhost:5436/fabric?sslmode=disable' $$ACTION

migrate-dev:
	@echo "up or down?: "; \
	read ACTION; \
	migrate -source file://ent/migrate/migrations -database 'postgres://root:$(db-password)@localhost:5432/fabric?sslmode=disable' $$ACTION

drop:
	migrate -source file://ent/migrate/migrations -database 'postgres://root:password@localhost:5436/fabric?sslmode=disable' drop

drop-dev:
	migrate -source file://ent/migrate/migrations -database 'postgres://root:$(db-password)@localhost:5432/fabric?sslmode=disable' drop

dev-db-pass:
	@echo $(db-password)

run:
		AUTO_MIGRATIONS=true go run cmd/server.go

rund:
		AUTO_MIGRATIONS=true DEBUG=true go run cmd/server.go

gen:
		rm gql/ent.graphql
		go generate ./...

model:
	@echo "Name of model you want to create?: "; \
	read NAME; \
	go run entgo.io/ent/cmd/ent new --target ./schema $$NAME

cm:
		@echo "Name of migration you want to create?: "; \
		read NAME; \
		go run ent/migrate/main.go $$NAME

proxy:
	./bin/cloud_sql_proxy -instances=fabric-375222:us-west1:convey-api-dev=tcp:5432