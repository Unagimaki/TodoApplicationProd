.PHONY: db-up

db-up:
	docker compose up -d
migrate-up:
	docker run --rm --network host \
		-v "$(CURDIR)/migrations:/migrations" \
		migrate/migrate \
		-path=/migrations \
		-database "postgres://todo_user:todo_password@localhost:5432/todo_db?sslmode=disable" \
		up