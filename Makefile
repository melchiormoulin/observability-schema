BINARY_SUFFIX_NAME=es-mapping
BINARY_NAME=protoc-gen-$(BINARY_SUFFIX_NAME) # if change, change the .gitignore
ELASTICSEARCH_INPUT_GO_TEMPLATE='template_in=examples/elasticsearch/input/mapping.template'
ELASTICSEARCH_OUTPUT_TEMPLATE_DIR='examples/elasticsearch/output/'
SCHEMA_PATH_DIR=examples/elasticsearch/input
SCHEMA_NAME=observabilitySchema.proto
all: build run
build:
	go build -o $(BINARY_NAME) cmd/schema-generator.go
run:
	protoc --plugin $(BINARY_NAME) --proto_path $(SCHEMA_PATH_DIR) --$(BINARY_SUFFIX_NAME)_opt=$(ELASTICSEARCH_INPUT_GO_TEMPLATE) --$(BINARY_SUFFIX_NAME)_out=$(ELASTICSEARCH_OUTPUT_TEMPLATE_DIR) $(SCHEMA_NAME)
fmt:
	go fmt ./... -v