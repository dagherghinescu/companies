#!/bin/bash

# Topic configuration
TOPIC_NAME="companies-events"
PARTITIONS=1
REPLICATION=1
BROKER="localhost:9092"

# Get the container ID for the Kafka container
CONTAINER_ID=$(docker ps --filter "ancestor=confluentinc/cp-kafka:7.5.0" --format "{{.ID}}")

if [ -z "$CONTAINER_ID" ]; then
  echo "Kafka container not running!"
  exit 1
fi

echo "Kafka container ID: $CONTAINER_ID"

# Create the topic inside the container
docker exec -it "$CONTAINER_ID" kafka-topics \
  --create \
  --bootstrap-server "$BROKER" \
  --replication-factor $REPLICATION \
  --partitions $PARTITIONS \
  --topic "$TOPIC_NAME"

if [ $? -eq 0 ]; then
  echo "Topic '$TOPIC_NAME' created successfully!"
else
  echo "Failed to create topic '$TOPIC_NAME'. It may already exist."
fi

# List topics to confirm
docker exec -it "$CONTAINER_ID" kafka-topics \
  --list \
  --bootstrap-server "$BROKER"
