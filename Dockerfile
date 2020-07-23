ARG GOLANG_VERSION=1.14.6-buster
ARG DEBIAN_VERSION=buster-slim
FROM golang:$GOLANG_VERSION as build
ARG BINARY_NAME=protoc-gen-es-mapping
WORKDIR /workspace
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o $BINARY_NAME cmd/schema-generator.go

FROM debian:$DEBIAN_VERSION
ARG PROTOC_VERSION=3.12.3
ARG BINARY_NAME=protoc-gen-es-mapping
ENV ENV_BINARY_NAME=protoc-gen-es-mapping
ENV ENV_BINARY_SUFFIX_NAME=es-mapping
RUN apt update && apt install wget unzip make -y -q
WORKDIR /protoc
RUN wget https://github.com/protocolbuffers/protobuf/releases/download/v$PROTOC_VERSION/protoc-$PROTOC_VERSION-linux-x86_64.zip && unzip protoc-$PROTOC_VERSION-linux-x86_64.zip && chmod +x bin/protoc && rm protoc-$PROTOC_VERSION-linux-x86_64.zip
COPY --from=build /workspace/$BINARY_NAME ./
ENTRYPOINT /protoc/bin/protoc --plugin $ENV_BINARY_NAME --proto_path /config/ --${ENV_BINARY_SUFFIX_NAME}_opt='template_in=/config/mapping.template' --${ENV_BINARY_SUFFIX_NAME}_out=/output/ observabilitySchema.proto