package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	plugin_go "github.com/golang/protobuf/protoc-gen-go/plugin"
	pb "github.com/melchiormoulin/observability-schema/schema"
	"github.com/melchiormoulin/observability-schema/elasticsearch"
	"google.golang.org/protobuf/types/descriptorpb"
	"io/ioutil"
	"os"
)

var template = flag.String("template", "mapping.template", "elasticsearch template to render.")

func getElasticsearchType(options *descriptorpb.FieldOptions) (proto.Message, error) {
	esFieldConfig, _ := proto.GetExtension(options, pb.E_ElasticsearchField)
	esFieldConfigString, _ := proto.GetExtension(options, pb.E_ElasticsearchFieldString)
	var esType proto.Message
	if esFieldConfig != nil {
		esType = esFieldConfig.(*pb.ElasticsearchField)
	}
	if esFieldConfigString != nil {
		esType = esFieldConfigString.(*pb.ElasticsearchFieldString)
	}
	return esType, nil
}

func main() {
	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return
	}
	req := &plugin_go.CodeGeneratorRequest{}
	resp := &plugin_go.CodeGeneratorResponse{}
	fieldsMapping := make(map[string]json.RawMessage)
	if err := proto.Unmarshal(data, req); err != nil {
		return
	}
	fmt.Fprintf(os.Stderr, ">>>>>>>>>>> generating for file %+v\n", req.FileToGenerate)
	protofile := req.ProtoFile
	for _, p := range protofile {
		//	fmt.Fprintf(os.Stderr, ">>>>> name %+v\n", p.GetName())
		messageTypes := p.GetMessageType()
		for _, messageType := range messageTypes {
			if messageType.GetName() == "ObservabilitySchema" {
				//fmt.Fprintf(os.Stderr, ">>>>>>>>>>> messageTYpe  %+v\n", messageType)
				fields := messageType.GetField()
				for _, field := range fields {
					options := field.GetOptions()
					jsonPbConfig := jsonpb.Marshaler{EmitDefaults: true, OrigName: true}
					esType, _ := getElasticsearchType(options)
					esTypebytes, _ := jsonPbConfig.MarshalToString(esType)
					fieldConfiguration := json.RawMessage(esTypebytes)
					name := field.GetName()
					fmt.Fprintf(os.Stderr, ">>>>>>>>>>> field %+v : %+v\n", field.GetName(), esTypebytes)
					fieldsMapping[name] = fieldConfiguration

				}
			}
		}
	}
	fieldsDefinition, err := json.MarshalIndent(fieldsMapping, "", "  ")
	fields := string(fieldsDefinition)
	fileName := "template.json"
	var buffer bytes.Buffer
	elasticsearch.Rendertemplate(*template, fields,&buffer)
	bufferStr := buffer.String()
	outputFile := plugin_go.CodeGeneratorResponse_File{Name: &fileName,Content:&bufferStr}
	resp.File = []*plugin_go.CodeGeneratorResponse_File {&outputFile}
	marshalled, err := proto.Marshal(resp)
	if err != nil {
		panic(err)
	}
	os.Stdout.Write(marshalled)



}
