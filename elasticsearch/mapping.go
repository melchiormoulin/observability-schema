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
	fieldsMapping    map[string]json.RawMessage // Elasticsearch Mapping with all fields name in keys with their definitions in value, rawMessage because of protobuf json serialization
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
	mapping.addField("@timestamp", &fieldDefinition)
}
func (mapping *Mapping) FieldsDefinition(protofile []*descriptorpb.FileDescriptorProto) string {
	for _, p := range protofile {
		messageTypes := p.GetMessageType()
		for _, messageType := range messageTypes {
			if messageType.GetName() == "ObservabilitySchema" {
				nestedTypes := messageType.GetNestedType()

				var oldNestTypesName []string

				for nestedTypes != nil {
					for _, nestedType := range nestedTypes {
						tmpFieldMap := make(map[string]json.RawMessage)
						fields := nestedType.GetField()
						for _, field := range fields {
							fieldsDefinitionBytes := mapping.parseField(field)
							tmpFieldMap[*field.Name] = fieldsDefinitionBytes
							fmt.Fprintf(os.Stderr, "from nested field  %+v : %+v\n", *field.Name, string(fieldsDefinitionBytes))
						}
						myJson, _ := json.Marshal(tmpFieldMap)
						if len(oldNestTypesName) == 0 {
							mapping.fieldsMapping[*nestedType.Name] = myJson

						} else {
							tmpFieldMap2 := make(map[string]json.RawMessage)
							tmpFieldMap2[*nestedType.Name] = myJson
						//	out := map[string]interface{}{}
						//	if len(oldNestTypesName) > 1 {
								str := mapping.String()
								path := strings.Join(oldNestTypesName, ".") + "." + *nestedType.Name
								str, _ = sjson.SetRaw(str, path, string(myJson))
								json.Unmarshal([]byte(str), &mapping.fieldsMapping)
								fmt.Fprintf(os.Stderr, "DEBUG  str: %+v \n", mapping.fieldsMapping)

					//		}
							//else {
							//	json.Unmarshal(mapping.fieldsMapping[oldNestTypesName[0]], &out)
							//	for key, value := range tmpFieldMap2 {
							//		out[key] = value
							//	}
							//	tmpJson, _ := json.Marshal(out)
							//	mapping.fieldsMapping[oldNestTypesName[0]] = tmpJson
							//}
						}
						nestedTypes = nestedType.GetNestedType()
						if len(nestedTypes) == 0 {
							oldNestTypesName = nil
						} else {
							oldNestTypesName = append(oldNestTypesName, *nestedType.Name)

						}
					}
				}

				fields := messageType.GetField()
				for _, field := range fields {
					fieldsDefinitionBytes := mapping.parseField(field)
					mapping.fieldsMapping[*field.Name] = fieldsDefinitionBytes
					fmt.Fprintf(os.Stderr, "basic field %+v : %+v\n", *field.Name, string(fieldsDefinitionBytes))
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
func (mapping *Mapping) parseField(field *descriptorpb.FieldDescriptorProto) json.RawMessage {
	fieldDefinition := getElasticsearchType(field.GetOptions())
	return mapping.addField(field.GetName(), fieldDefinition)
}

func (mapping *Mapping) addField(fieldName string, fieldDefinition proto.Message) json.RawMessage {
	fieldsDefinitionBytes, _ := mapping.protoJSON.Marshal(fieldDefinition) //Can't use the basic encoding/json because we can't use EmitUnpopulated: true with the basic json package.
	return fieldsDefinitionBytes
}
