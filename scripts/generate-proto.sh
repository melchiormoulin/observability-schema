protoc --proto_path ../schema --go_out ../schema --go_opt=paths=source_relative ../schema/elasticsearch.proto
protoc --proto_path ../schema --go_out ../schema --go_opt=paths=source_relative ../schema/observabilitySchema.proto
