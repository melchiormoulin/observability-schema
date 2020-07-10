package elasticsearch

import (
	"encoding/json"
	pb "github.com/melchiormoulin/observability-schema/schema"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"testing"
)

func TestMappingInit(t *testing.T) {
	mapping := MappingInit(false, "  ", "")
	expectedFormatIndent := "  "
	if expectedFormatIndent != mapping.formatIndent {
		t.Errorf("the format indent has not been set correctly expected : %+v find : %+v", expectedFormatIndent, mapping.formatIndent)
	}
	expectedFormatPrefix := ""
	if expectedFormatPrefix != mapping.formatPrefix {
		t.Errorf("the format prefix has not been set correctly expected : %+v find : %+v", expectedFormatPrefix, mapping.formatPrefix)
	}
	expectedProtoJson := protojson.MarshalOptions{EmitUnpopulated: true, UseProtoNames: true}

	if mapping.protoJson != expectedProtoJson {
		t.Errorf("the protojson marshaler has not been set correctly expected : %+v find : %+v", expectedProtoJson, mapping.protoJson)
	}
	if mapping.fieldsMapping == nil {
		t.Errorf("the fieldsMapping map is nil")
	}
}
func TestMappingInitWithTimestamp(t *testing.T) {
	mapping := MappingInit(true, "  ", "")
	if mapping.fieldsMapping["@timestamp"] != nil {
		t.Errorf("@timestamp key should be set")
	}
}

func TestAddField(t *testing.T) {
	mapping := MappingInit(true, "  ", "")
	expectedFieldDefinitionStruct := pb.ElasticsearchFieldString{Type: "keyword", DocValues: true, Index: true}
	mapping.addField("@timestamp", &expectedFieldDefinitionStruct)
	fieldsDefinitionBytes, _ := mapping.protoJson.Marshal(&expectedFieldDefinitionStruct)
	expectedFieldDefinitionTmp := json.RawMessage(fieldsDefinitionBytes)
	fieldDefinition, _ := mapping.fieldsMapping["@timestamp"].MarshalJSON()
	expectedFieldDefinition, _ := expectedFieldDefinitionTmp.MarshalJSON()
	if string(fieldDefinition) != string(expectedFieldDefinition) {
		t.Errorf("@timestamp field definition should be present")
	}
}

func TestGetElasticsearchType(t *testing.T) {
	options := descriptorpb.FieldOptions{}
	fieldStringType := pb.ElasticsearchFieldString{Type: "keyword", DocValues: true, Index: true}
	proto.SetExtension(&options, pb.E_ElasticsearchFieldString, &fieldStringType)
	msg := getElasticsearchType(&options)
	if msg == nil {
		t.Errorf("msg should not be nil")
	}
	options2 := descriptorpb.FieldOptions{}
	fieldStringType2 := pb.ElasticsearchField{Type: "keyword", DocValues: true, Index: true}
	proto.SetExtension(&options2, pb.E_ElasticsearchField, &fieldStringType2)
	msg = getElasticsearchType(&options2)
	if msg == nil {
		t.Errorf("msg should not be nil")
	}
}
func TestGetElasticsearchTypePanic(t *testing.T) {
	options := descriptorpb.FieldOptions{}
	assert.Panics(t,func() { getElasticsearchType(&options) },"should panic has it is not an elasticsearchtype")
}

func TestParseField(t *testing.T) {
	userId := "user_id"
	options := descriptorpb.FieldOptions{}
	fieldStringType := pb.ElasticsearchFieldString{Type: "keyword", DocValues: true, Index: true}
	proto.SetExtension(&options, pb.E_ElasticsearchFieldString, &fieldStringType)
	desc := descriptorpb.FieldDescriptorProto{Name: &userId, Options: &options}
	mapping := MappingInit(false, "  ", "")
	mapping.parseField(&desc)
	if mapping.fieldsMapping[userId] == nil {
		t.Errorf("the field %+v should be present", userId)
	}
}

func TestFieldDefinitions(t *testing.T) {
	userId := "user_id"
	options := descriptorpb.FieldOptions{}
	fieldStringType := pb.ElasticsearchFieldString{Type: "keyword", DocValues: true, Index: true}
	proto.SetExtension(&options, pb.E_ElasticsearchFieldString, &fieldStringType)
	desc := descriptorpb.FieldDescriptorProto{Name: &userId, Options: &options}
	mapping := MappingInit(false, "  ", "")
	observabilitySchema := "ObservabilitySchema"
	descProto := descriptorpb.DescriptorProto{Name: &observabilitySchema, Field: []*descriptorpb.FieldDescriptorProto{&desc}}
	fileProto := descriptorpb.FileDescriptorProto{MessageType: []*descriptorpb.DescriptorProto{&descProto}}
	mapping.FieldsDefinition([]*descriptorpb.FileDescriptorProto{&fileProto})
	if mapping.fieldsMapping[userId] == nil {
		t.Errorf("the field %+v should be present", userId)
	}
}
