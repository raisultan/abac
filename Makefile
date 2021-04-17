up_services:
	docker-compose up --build -d

stop_services:
	docker-compose stop

restart_services:
	docker-compose stop
	docker-compose up --build -d

down_services:
	docker-compose down

migrate_up:
	docker exec -it $$(docker ps | grep server_ | awk '{{ print $$1 }}') sh -c "migrate -source file:/app/migrations -database \$$POSTGRES_URL up"

migrate_down:
	docker exec -it $$(docker ps | grep server_ | awk '{{ print $$1 }}') sh -c "migrate -source file:/app/migrations -database \$$POSTGRES_URL down"
