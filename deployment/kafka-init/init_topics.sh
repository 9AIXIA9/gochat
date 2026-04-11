#!/bin/sh
set -eu

BOOTSTRAP_SERVER="${KAFKA_BOOTSTRAP_SERVERS:-kafka:9092}"
PARTITIONS="${KAFKA_PARTITIONS:-1}"
REPLICATION_FACTOR="${KAFKA_REPLICATION_FACTOR:-1}"
TOPICS_FILE="${KAFKA_TOPICS_FILE:-/topics.txt}"

resolve_kafka_topics_cmd() {
  if command -v kafka-topics.sh >/dev/null 2>&1; then
    echo "kafka-topics.sh"
    return 0
  fi
  if command -v kafka-topics >/dev/null 2>&1; then
    echo "kafka-topics"
    return 0
  fi
  if [ -x /opt/kafka/bin/kafka-topics.sh ]; then
    echo "/opt/kafka/bin/kafka-topics.sh"
    return 0
  fi

  if [ -x /usr/bin/kafka-topics ]; then
    echo "/usr/bin/kafka-topics"
    return 0
  fi

  if [ -x /usr/bin/kafka-topics.sh ]; then
    echo "/usr/bin/kafka-topics.sh"
    return 0
  fi

  if [ -x /opt/bitnami/kafka/bin/kafka-topics.sh ]; then
    echo "/opt/bitnami/kafka/bin/kafka-topics.sh"
    return 0
  fi

  echo "kafka-topics command not found" >&2
  return 1
}

wait_for_broker() {
  cmd="$1"
  retries="${KAFKA_WAIT_RETRIES:-30}"
  sleep_seconds="${KAFKA_WAIT_SECONDS:-2}"

  i=1
  while [ "$i" -le "$retries" ]; do
    if "$cmd" --bootstrap-server "$BOOTSTRAP_SERVER" --list >/dev/null 2>&1; then
      return 0
    fi
    echo "[$i/$retries] waiting for Kafka broker at $BOOTSTRAP_SERVER"
    sleep "$sleep_seconds"
    i=$((i + 1))
  done

  echo "Kafka broker is not ready: $BOOTSTRAP_SERVER" >&2
  return 1
}

create_topic_if_missing() {
  cmd="$1"
  topic="$2"

  if "$cmd" --bootstrap-server "$BOOTSTRAP_SERVER" --list | grep -Fx "$topic" >/dev/null 2>&1; then
    echo "topic already exists: $topic"
    return 0
  fi

  echo "creating topic: $topic"
  "$cmd" --bootstrap-server "$BOOTSTRAP_SERVER" \
    --create \
    --if-not-exists \
    --topic "$topic" \
    --partitions "$PARTITIONS" \
    --replication-factor "$REPLICATION_FACTOR"
}

load_topics() {
  if [ -n "${KAFKA_TOPICS:-}" ]; then
    printf '%s\n' "$KAFKA_TOPICS" | tr ' ' '\n' | tr -d '\r'
    return 0
  fi

  if [ ! -f "$TOPICS_FILE" ]; then
    echo "topics file not found: $TOPICS_FILE" >&2
    return 1
  fi

  # Read topics from file while skipping blanks and comments.
  grep -v '^[[:space:]]*#' "$TOPICS_FILE" | sed '/^[[:space:]]*$/d' | tr -d '\r'
}

main() {
  kafka_topics_cmd="$(resolve_kafka_topics_cmd)"

  wait_for_broker "$kafka_topics_cmd"

  load_topics | while IFS= read -r topic; do
    [ -z "$topic" ] && continue
    create_topic_if_missing "$kafka_topics_cmd" "$topic"
  done

  echo "Kafka topic initialization completed"
}

main
