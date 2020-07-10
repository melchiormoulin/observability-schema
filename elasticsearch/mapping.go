package elasticsearch

import (
	"encoding/json"
	"fmt"
	pb "github.com/melchiormoulin/observability-schema/schema"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"os"
)

type Mapping struct {
	fieldsMapping    map[string]json.RawMessage
	protoJson        protojson.MarshalOptions
	fieldsDefinition string
	formatIndent     string
	formatPrefix     string
}

func getElasticsearchType(options *descriptorpb.FieldOptions) proto.Message {
	if proto.HasExtension(options, pb.E_ElasticsearchField) {
		esFieldConfig := proto.GetExtension(options, pb.E_ElasticsearchField)
		return esFieldConfig.(proto.Message)
	} else if proto.HasExtension(options, pb.E_ElasticsearchFieldString) {
		esFieldConfigString := proto.GetExtension(options, pb.E_ElasticsearchFieldString)
		return esFieldConfigString.(proto.Message)
	}
	panic(fmt.Errorf("bad protobuf option type"))
}

func MappingInit(withTimestampField bool, formatIndent string, formatPrefix string) Mapping {
	protoJson := protojson.MarshalOptions{EmitUnpopulated: true, UseProtoNames: true}
	fieldsMapping := make(map[string]json.RawMessage)
	mapping := Mapping{
		fieldsMapping: fieldsMapping,
		protoJson:     protoJson,
		formatIndent:  formatIndent,
		formatPrefix:  formatPrefix,
	}
	if withTimestampField {
		mapping.addTimestampField()
	}
	return mapping
}
func (mapping *Mapping) addTimestampField() {
	fieldDefinition := pb.ElasticsearchFieldString{Type: "keyword", DocValues: true, Index: true}
	mapping.addField("@timestamp", &fieldDefinition)
}
func (mapping *Mapping) FieldsDefinition(protofile []*descriptorpb.FileDescriptorProto) string {
	for _, p := range protofile {
			//fmt.Fprintf(os.Stderr, ">>>>> name %+v\n", p)

		messageTypes := p.GetMessageType()
		for _, messageType := range messageTypes {
			if messageType.GetName() == "ObservabilitySchema" {
				fields := messageType.GetField()
				for _, field := range fields {
					mapping.parseField(field)
				}

			}
		}
	}
	return mapping.String()
}
func (mapping *Mapping) String() string {
	fieldsDefinition, err := json.MarshalIndent(mapping.fieldsMapping, mapping.formatPrefix, mapping.formatIndent)
	if err != nil {
		panic(err)
	}
	return string(fieldsDefinition)
}
func (mapping *Mapping) parseField(field *descriptorpb.FieldDescriptorProto) {
	fieldDefinition := getElasticsearchType(field.GetOptions())
	mapping.addField(field.GetName(), fieldDefinition)
}

func (mapping *Mapping) addField(fieldName string, fieldDefinition proto.Message) {
	fieldsDefinitionBytes, _ := mapping.protoJson.Marshal(fieldDefinition) //Can't use the basic encoding/json because we can't use EmitUnpopulated: true with the basic json package.
	mapping.fieldsMapping[fieldName] = json.RawMessage(string(fieldsDefinitionBytes))
	fmt.Fprintf(os.Stderr, "field %+v : %+v\n", fieldName, fieldDefinition)
}
