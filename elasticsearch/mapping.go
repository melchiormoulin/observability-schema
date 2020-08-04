package elasticsearch

import (
	"encoding/json"
	"fmt"
	pb "github.com/melchiormoulin/observability-schema/schema"
	"github.com/tidwall/sjson"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"os"
	"strings"
)

//Mapping struct is to get the generated mapping
type Mapping struct {
	fieldsMapping    map[string]json.RawMessage // Elasticsearch Mapping with all fields name in keys with their elasticsearch definitions in value, rawMessage because of protobuf json serialization
	protoJSON        protojson.MarshalOptions   // needed for json serialization
	fieldsDefinition string                     // the json string of the mapping
	formatIndent     string                     // format param for fieldsDefinition
	formatPrefix     string                     //format param for fieldsDefinition
}

//Only two types are allowed for now for elasticsearch
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

//Get a new instance of Mapping struct
func MappingInit(withTimestampField bool, formatIndent string, formatPrefix string) Mapping {
	protoJson := protojson.MarshalOptions{EmitUnpopulated: true, UseProtoNames: true}
	fieldsMapping := make(map[string]json.RawMessage)
	mapping := Mapping{
		fieldsMapping: fieldsMapping,
		protoJSON:     protoJson,
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
	fieldsDefinitionBytes, err := mapping.protoJSON.Marshal(&fieldDefinition) //Can't use the basic encoding/json because we can't use EmitUnpopulated: true with the basic json package.
	if err!=nil {
		panic(err)
	}
	mapping.fieldsMapping["@timestamp"]= fieldsDefinitionBytes
}
func (mapping *Mapping) FieldsDefinition(protofile []*descriptorpb.FileDescriptorProto) string {
	for _, p := range protofile {
		messageTypes := p.GetMessageType()
		for _, messageType := range messageTypes {
			if messageType.GetName() == "ObservabilitySchema" {
				mapping.parseNestedTypes(messageType.GetNestedType(), []string{})
				fmt.Fprintf(os.Stderr, "nested type root level\n")
				mapping.parseFields(messageType.GetField(), mapping.fieldsMapping)
			}
		}
	}
	return mapping.String()
}

func (mapping *Mapping) parseFields(fields []*descriptorpb.FieldDescriptorProto, jsonKv map[string]json.RawMessage) {
	for _, field := range fields {
		fieldsDefinitionBytes := mapping.parseField(field)
		jsonKv[*field.Name] = fieldsDefinitionBytes
		fmt.Fprintf(os.Stderr, "field %+v : %+v\n", *field.Name, string(fieldsDefinitionBytes))
	}
}

func (mapping *Mapping) parseNestedTypes(nestedTypes []*descriptorpb.DescriptorProto, jsonPath []string) {
	if nestedTypes == nil { //recursive
		return
	}
	for _, nestedType := range nestedTypes {
		fmt.Fprintf(os.Stderr, "nested type %+v \n", *nestedType.Name)
		tmpFieldMap := make(map[string]json.RawMessage)
		mapping.parseFields(nestedType.GetField(), tmpFieldMap)
		mapping.parseNestedType(tmpFieldMap, jsonPath, *nestedType.Name)
		children := nestedType.GetNestedType()
		jsonPathChildren := getJsonPathChildren(children, jsonPath, *nestedType.Name)
		mapping.parseNestedTypes(children, jsonPathChildren) // recursive
	}
}

func (mapping *Mapping) parseNestedType(tmpFieldMap map[string]json.RawMessage, jsonPathChildren []string, currentNodeName string) {
	myJson, _ := json.Marshal(tmpFieldMap)
	str := mapping.String()
	var err error
	currentNodeName = currentNodeName +".properties"
	str, err = sjson.SetRaw(str, getJsonPath(jsonPathChildren, currentNodeName), string(myJson))
	if err!=nil {
		panic(err)
	}
	json.Unmarshal([]byte(str), &mapping.fieldsMapping)
}

func getJsonPathChildren(childrenNestedTypes []*descriptorpb.DescriptorProto, jsonPath []string, currentNode string) []string {
	if len(childrenNestedTypes) > 0 {
		currentNode = currentNode + ".properties"
		jsonPath = append(jsonPath, currentNode)
	} else {
		jsonPath = nil
	}
	return jsonPath
}

func getJsonPath(jsonPathFather []string, currentNodeName string) string {
	jsonPath := currentNodeName
	if len(jsonPathFather) > 0 {
		jsonPath = strings.Join(jsonPathFather, ".") + "." + jsonPath
	}
	return jsonPath
}

func (mapping *Mapping) String() string {
	fieldsDefinition, err := json.MarshalIndent(mapping.fieldsMapping, mapping.formatPrefix, mapping.formatIndent)
	if err != nil {
		panic(err)
	}
	return string(fieldsDefinition)
}
func (mapping *Mapping) parseField(field *descriptorpb.FieldDescriptorProto) json.RawMessage {
	fieldDefinition := getElasticsearchType(field.GetOptions())
	fieldsDefinitionBytes, err := mapping.protoJSON.Marshal(fieldDefinition) //Can't use the basic encoding/json because we can't use EmitUnpopulated: true with the basic json package.
	if err!=nil {
		panic(err)
	}
	return fieldsDefinitionBytes
}

