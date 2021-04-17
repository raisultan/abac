up_services:
	docker-compose up --build -d

stop_services:
	docker-compose stop

down_services:
	docker-compose down

migrate_up:
	docker exec -it $$(docker ps | grep web_ | awk '{{ print $$1 }}') sh -c "migrate -source file:/app/migrations -database \$$POSTGRES_URL up"

migrate_down:
	docker exec -it $$(docker ps | grep web_ | awk '{{ print $$1 }}') sh -c "migrate -source file:/app/migrations -database \$$POSTGRES_URL down"
