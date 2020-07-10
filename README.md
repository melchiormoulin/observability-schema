# observability-schema
A protoc plugin that generates elasticsearch template mapping from a protobuf schema.

When you are logging across multiple applications fields name you can have :
1. A bad data coherence for example the same data with different names `user-id user_id userId`
2. A mapping conflict for example `user_id` with `long` type and `string` type

The purpose of this plugin is to provide a common schema for your observability fields definitions across multiple applications.

This repo provides:

1. It provides a protobuf schema in order to define all fields names with their corresponding types for all applications.
2. The tool is the protoc plugin that parse schema and generate an elasticsearch template mapping.

How to use it ?

Install `golang` `protoc` and `make` and launch the `make` command.
Look at the input and output examples of the plugin  in `examples/elasticsearch/`