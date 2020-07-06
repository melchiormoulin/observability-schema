package elasticsearch

import (
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	pb "github.com/melchiormoulin/observability-schema/schema"
	"google.golang.org/protobuf/types/descriptorpb"
	"os"
)

type Mapping struct {
	fieldsMapping    map[string]json.RawMessage
	jsonPbConfig     jsonpb.Marshaler
	fieldsDefinition string
}

func getElasticsearchType(options *descriptorpb.FieldOptions) proto.Message {
	esFieldConfig, _ := proto.GetExtension(options, pb.E_ElasticsearchField)
	esFieldConfigString, _ := proto.GetExtension(options, pb.E_ElasticsearchFieldString)
	var esType proto.Message
	if esFieldConfig != nil {
		esType = esFieldConfig.(*pb.ElasticsearchField)
	} else if esFieldConfigString != nil {
		esType = esFieldConfigString.(*pb.ElasticsearchFieldString)
	} else {
		panic(fmt.Errorf("bad protobuf option type"))
	}
	return esType
}

func MappingInit() Mapping {
	jsonPbConfig := jsonpb.Marshaler{EmitDefaults: true, OrigName: true}
	fieldsMapping := make(map[string]json.RawMessage)
	mapping := Mapping{
		fieldsMapping: fieldsMapping,
		jsonPbConfig:  jsonPbConfig,
	}
	mapping.addTimestampField()
	return mapping
}
func (mapping *Mapping) addTimestampField() {
	fieldDefinition := pb.ElasticsearchFieldString{Type:"keyword",DocValues: true, Index: true}
	mapping.addField("@timestamp",&fieldDefinition)
}
func (mapping *Mapping) FieldsDefinition(protofile []*descriptorpb.FileDescriptorProto) string {
	for _, p := range protofile {
		//	fmt.Fprintf(os.Stderr, ">>>>> name %+v\n", p.GetName())
		messageTypes := p.GetMessageType()
		for _, messageType := range messageTypes {
			if messageType.GetName() == "ObservabilitySchema" {
				fields := messageType.GetField()
				//fmt.Fprintf(os.Stderr, ">>>>>>>>>>> messageTYpe  %+v\n", messageType)
				for _, field := range fields {
					mapping.parseField(field)
				}

			}
		}
	}
	return mapping.String()
}
func (mapping *Mapping) String() string {
	fieldsDefinition, err := json.MarshalIndent(mapping.fieldsMapping, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(fieldsDefinition)
}
func (mapping *Mapping) parseField(field *descriptorpb.FieldDescriptorProto) {
	fieldDefinition := getElasticsearchType(field.GetOptions())
	mapping.addField(field.GetName(),fieldDefinition)

}

func(mapping *Mapping) addField(fieldName string,fieldDefinition proto.Message) {
	fieldDefinitionBytes, err := mapping.jsonPbConfig.MarshalToString(fieldDefinition)
	if err != nil {
		panic(fmt.Errorf("error during protobuf marshaling %+v", err))
	}
	mapping.fieldsMapping[fieldName] = json.RawMessage(fieldDefinitionBytes)
	fmt.Fprintf(os.Stderr, "field %+v : %+v\n", fieldName, fieldDefinitionBytes)

}
