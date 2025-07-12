#!/bin/bash

set -e

echo "Initializing Kafka KRaft cluster..."

mkdir -p /var/lib/kafka/data
chown -R appuser:appuser /var/lib/kafka/data
chmod -R 755 /var/lib/kafka/data

if [ ! -f /var/lib/kafka/data/meta.properties ]; then
    echo "Formatting Kafka storage directory..."
    kafka-storage.sh format -t 4L6g3nShT-eMCtK--X86sw -c /etc/kafka/kafka.properties
fi

echo "Starting Kafka..."
exec kafka-server-start.sh /etc/kafka/kafka.properties 