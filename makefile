# Makefile

DOCKER_COMPOSE=docker-compose

up: 
	$(DOCKER_COMPOSE) build
	$(DOCKER_COMPOSE) up -d

down:
	$(DOCKER_COMPOSE) down

kafka-topic:
	./create_kafka_topic.sh

kafka-consume:
	./kafka_consume.sh

test:
	./tests.sh