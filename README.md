# Logging infrastructure
Simple log parsing and ingestion using well known tech stacks

## Problem statement
A microservice generates unstructured log using this pattern/format
```
[timestamp] [service_name] [status_code] [response_time_ms] [user_id]
[transaction_id] [additional_info]
```

example entry
```
2023-08-15 13:45:00 checkout 200 120ms user1234 tx5678 Purchased iPhone 13
```

We need to store the incoming logs as structured format (i.e: JSON) thus it can be queried and visualized to provide better debugging and monitoring experience for developer.

## Components
For this case, we can use elasticsearch for log storage as it offer full-text search with good performance, it also has Kibana to simplify log visualization. For logging agent, we can use fluentbit for starting point as it has small overhead (memory and CPU)
- Elasticsearch (v8.7.1)
- Kibana (v8.7.1)
- Fluentbit (v2.2.2)
- Docker + compose (25.0.0), to run the solution locally

## Architecture
### Design
For small - medium workload, we can start with the very simple architecture, it uses minimum part to be maintained.
Fluentbit being used for "ETL" pipeline, it reads each entry (line) from the log file, transform/structure the log using regex, then store it to elasticsearch.
![Simple logging infrastructure](/assets/simple.png)
![Regex pattern to extract string part](/assets/regex.png)


For medium-large workloads, we can consider adding a buffer layer that holds the entry temporarily before doing further processing or storing it in persistent storage, we can use kafka, pub/sub, or similar tools, so the architecture will be like this
![Advanced logging infrastructure](/assets/with_buffer.png)
- The first fluentbit will only responsible for shipping the raw log to buffer layer, we should reduce operation here since we want to maximize its throughput for sending large volume log entry. We can consider tuning chunk/buffer config of fluentbit to avoid backpressure scenario also reducing network connection count.
- Kafka will be used as buffer layer that holds the logs temporarily, we can differentiate topics for each microservice and configure partitions accordingly
- We can use fluentbit or custom code to consume kafka topic and do some operations before store the log in elasticsearch.
- Eventhough we can store huge log streams in Kafka, without proper consumers we will be increasing the Kafka lag, so we can adjust the number of consumers accordingly  

### Alerting
We can utilize kibana Observability to set any alerting rules and sent it to custom destination (e.g: email, webhook) but we need a license (at least Plati­num) in order to use the feature. Other alternative is ElastAlert which we can use for free. Both of alerting systems are relying on ES data.

## Run the solution
Make sure you have docker and docker-compose installed on your machine.

p.s: If you're using minikube/rancher desktop/ or any VM based docker as replacement to docker desktop, make sure to mount this project directory to the VM, in my case i use minikube
```sh
$ minikube -p docker-vm mount $(pwd):$(pwd)
```

Run the container
```sh
$ docker compose up
```

## Cleanup
Delete all docker compose components (including its volume)
```sh
$ docker compose down -v
```

## TODO
- Create elasticsearch index pattern / data view using code
- Create kibana dashboard using code

## Appendix
- Kibana log view
![Kibana log view](/assets/kibana_01.png)
![Kibana log view](/assets/kibana_02.png)
- Kibana dashboard
![Kibana dashboard](/assets/kibana_03.png)