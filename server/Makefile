db:
	@docker compose up -d go_db

build:
	@echo "Building the docker image..."
	@docker compose build	-q
	@echo "Done!"

run:
	@echo "Running the docker image..."
	@powershell -Command try { docker-compose up musicapp --watch} finally { docker-compose stop musicapp-db pgadmin }


psql:
	@docker compose exec musicapp-db psql -h musicapp-db -U root -d musicapp
