# start app
.PHONY : run
run:
	go mod download && go run main.go --config="config/dev.yml"

.PHONY : docker-up
docker-up:
	docker-compose -f docker-compose.yml up --build

# Таргет для применения миграций базы данных
migrate:
	@echo "Applying database migrations..."
	goose -dir ./migrations postgres "host=localhost port=5432 user=brands password=pgpwdbrands dbname=brands sslmode=disable" up
	@echo "Database migrations applied successfully."
