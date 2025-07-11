#!/bin/bash
set -e

if [ -n "$KAFKA_ADVERTISED_LISTENERS" ]; then
  sed -i "s|advertised.listeners=.*|advertised.listeners=$KAFKA_ADVERTISED_LISTENERS|g" /opt/kafka/config/kraft/server.properties
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
  /opt/kafka/bin/kafka-broker-api-versions.sh --bootstrap-server localhost:9092
}

echo "Waiting for Kafka to start..."
while ! check_kafka_ready; do
  sleep 1
done

if [ -n "$KAFKA_TOPIC" ]; then
  echo "Creating topic: $KAFKA_TOPIC"
  /opt/kafka/bin/kafka-topics.sh --create \
    --topic "$KAFKA_TOPIC" \
    --partitions "$PARTITIONS" \
    --replication-factor "$REPLICATION_FACTOR" \
    --bootstrap-server localhost:9092
fi

wait $KAFKA_PID