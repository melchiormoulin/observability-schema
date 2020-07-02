all: build run
build:
	go build -o protoc-gen-estemplate cmd/schema-generator.go
run:
	protoc --plugin protoc-gen-estemplate --proto_path schema --estemplate_out=examples/elasticsearch/ observabilitySchema.proto
