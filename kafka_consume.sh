#!/usr/bin/env bash

TOPIC=${1:-companies-events}
CONTAINER_ID=${2:-$(docker ps --filter "ancestor=confluentinc/cp-kafka:7.5.0" --format "{{.ID}}" | head -n1)}

if [ -z "$CONTAINER_ID" ]; then
  echo "Kafka container not found. Please start Kafka or provide container ID."
  exit 1
fi

echo "Using Kafka container: $CONTAINER_ID"
echo "Consuming messages from topic: $TOPIC"

docker exec -it "$CONTAINER_ID" kafka-console-consumer \
  --bootstrap-server localhost:9092 \
  --topic "$TOPIC" \
  --from-beginning \
  --timeout-ms 5000
