# GOAL
Goal is a collections of helpers to speed up microservice development

## Motivation 
I wrote this library to organize basic functions that I need for writing microservices in a single place. Goal is accordingly a common library that includes the following features.

- It has an easy to use ETCD client that allows for 
  - Automatic configuration reading 
  - Bootstrapping configurations 
- It provides NATS codecs for
  - Protobuf 
  - Protobuf + ZSTD
- It provides an in-memory cache which has a built-in TTL
- It provides high performance collections
  - Queue
  - Stack
- It provides a high performance DI container that support the following life-cycles
  - Singleton 
  - Transient 
  - Scoped 
- It provides a wrapper around the original http client in Go that 
  - Ensures connection reuse 
  - Simplifies all http operations in a single function 
- It provides an advanced logger that
  - Sends measurements to InfluxDB 
  - Has failover features 
  - Records performance benchmarks at runtime without any performance penalty 
- It provides a Protobuf util that
  - Marshals Protobuf to Map 
  - Unmarshals Map to Protobuf 



