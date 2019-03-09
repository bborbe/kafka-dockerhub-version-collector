# Kafka Dockerhub Version Collector

Publishes available Dockerhub versions to a Kafka topic.

## Run version collector

```bash
go run main.go \
-kafka-brokers=kafka:9092 \
-kafka-topic=application-version-available \
-kafka-schema-registry-url=http://schema-registry:8081 \
-repositories=library/traefik,library/ubuntu \
-v=2
```
