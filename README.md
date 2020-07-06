# observability-schema

Issues :

When you are logging accross multiple applications fields name you can have :
1. A bad data coherence for example the same data with different names `user-id user_id userId`
2. A mapping conflict for example `user_id` with `long` type and `string` type

The goal of this tool is to provide a common schema for your observability fields definitions across multiple applications.

How this tool solves those issues ?

1. It provides a protobuf schema in order to define all fields names with their corresponding types for all applications.
2. The tool is the protoc plugin that parse schema and generate an elasticsearch mapping.


How to use it ?

Install `golang` `protoc` and `make` and launch the `make` command.
Look at the generated elasticsearch template in output/
