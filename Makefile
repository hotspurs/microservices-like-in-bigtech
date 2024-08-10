start:
	docker-compose -f ./services/docker-compose.yml up --build
stop:
	docker-compose -f ./services/docker-compose.yml down