package elasticsearch

import (
	pb "github.com/melchiormoulin/observability-schema/schema"
	"google.golang.org/protobuf/encoding/protojson"
	"testing"
)

func TestMappingInit(t *testing.T) {
	mapping :=MappingInit(false,"  ","")
	expectedFormatIndent := "  "
	if expectedFormatIndent != mapping.formatIndent  {
		t.Errorf("the format indent has not been set correctly expected : %+v find : %+v",expectedFormatIndent,mapping.formatIndent)
	}
	expectedFormatPrefix := ""
	if expectedFormatPrefix != mapping.formatPrefix {
		t.Errorf("the format prefix has not been set correctly expected : %+v find : %+v",expectedFormatPrefix,mapping.formatPrefix)
	}
	expectedProtoJson :=  protojson.MarshalOptions{EmitUnpopulated: true,UseProtoNames: true}

	if mapping.protoJson != expectedProtoJson {
		t.Errorf("the protojson marshaler has not been set correctly expected : %+v find : %+v",expectedProtoJson,mapping.protoJson)
	}
	if mapping.fieldsMapping == nil {
		t.Errorf("the fieldsMapping map is nil")
	}
}
func TestMappingInitWithTimestamp(t *testing.T) {
	mapping :=MappingInit(true,"  ","")
	if mapping.fieldsMapping["@timestamp"] !=nil {
		t.Errorf("@timestamp key should be set")
	}
}

func TestAddField(t *testing.T) {
	expectedFieldDefinition :=  pb.ElasticsearchFieldString{Type:"keyword",DocValues: true, Index: true}
	mapping :=MappingInit(true,"  ","")
	mapping.addField("@timestamp",&expectedFieldDefinition)
	//tmp,_ := mapping.fieldsMapping["@timestamp"].(*pb.ElasticsearchFieldString)
	//if tmp.Type != expectedFieldDefinition.Type || tmp.DocValues != expectedFieldDefinition.DocValues || tmp.Index != expectedFieldDefinition.Index || tmp.Norms != expectedFieldDefinition.Norms ||  tmp.Store != expectedFieldDefinition.Store{
	//	t.Errorf("error to add field")

	//}
}
