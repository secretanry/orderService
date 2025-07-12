#!/bin/bash
set -e

if [ -n "$KAFKA_ADVERTISED_LISTENERS" ]; then
  sed -i "s|advertised.listeners=.*|advertised.listeners=$KAFKA_ADVERTISED_LISTENERS|g" /opt/kafka/config/kraft/server.properties
fi

if [ -n "$KAFKA_CONTROLLER_QUORUM_VOTERS" ]; then
  sed -i "s|controller.quorum.voters=.*|controller.quorum.voters=$KAFKA_CONTROLLER_QUORUM_VOTERS|g" /opt/kafka/config/kraft/server.properties
fi

if [ ! -f /tmp/kraft-combined-logs/meta.properties ]; then
  echo "Initializing KRaft storage..."
  export KAFKA_CLUSTER_ID="$(/opt/kafka/bin/kafka-storage.sh random-uuid)"
  /opt/kafka/bin/kafka-storage.sh format -t $KAFKA_CLUSTER_ID -c /opt/kafka/config/kraft/server.properties
fi

echo "Starting Kafka..."
/opt/kafka/bin/kafka-server-start.sh /opt/kafka/config/kraft/server.properties &
KAFKA_PID=$!

check_kafka_ready() {
  /opt/kafka/bin/kafka-broker-api-versions.sh --bootstrap-server localhost:9092 >/dev/null 2>&1
}

echo "Waiting for Kafka to start..."
attempts=0
max_attempts=60
while ! check_kafka_ready && [ $attempts -lt $max_attempts ]; do
  echo "Attempt $((attempts + 1))/$max_attempts: Kafka not ready yet..."
  sleep 2
  attempts=$((attempts + 1))
done

if [ $attempts -eq $max_attempts ]; then
  echo "Kafka failed to start within $max_attempts attempts"
  exit 1
fi

echo "Kafka is ready!"

if [ -n "$KAFKA_TOPIC" ]; then
  echo "Creating topic: $KAFKA_TOPIC"
  /opt/kafka/bin/kafka-topics.sh --create \
    --topic "$KAFKA_TOPIC" \
    --partitions "$PARTITIONS" \
    --replication-factor "$REPLICATION_FACTOR" \
    --bootstrap-server localhost:9092
fi

wait $KAFKA_PID