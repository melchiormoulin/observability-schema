package protoplugin

import (
	"bytes"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"
	"strings"
)

const (
	templatePath = "examples/elasticsearch/input/mapping.template"
	templateInParam        = "template_in"
	templateOutputFileName = "template.json"
)

func ReqInit(data []byte) *pluginpb.CodeGeneratorRequest {
	req := &pluginpb.CodeGeneratorRequest{}
	if err := proto.Unmarshal(data, req); err != nil {
		panic(err)
	}
	return req
}
func OutputStructSerialized(buffer *bytes.Buffer) []byte {
	resp := &pluginpb.CodeGeneratorResponse{}
	bufferStr := buffer.String()
	fileName := templateOutputFileName
	outputFile := pluginpb.CodeGeneratorResponse_File{Name: &fileName, Content: &bufferStr}
	resp.File = []*pluginpb.CodeGeneratorResponse_File{&outputFile}
	output, err := proto.Marshal(resp)
	if err != nil {
		panic(err)
	}
	return output
}

func TemplatePath(param string) string {
	keyValueParam := strings.Split(param, "=")
	template := templatePath
	if len(keyValueParam) == 2 && keyValueParam[0] == templateInParam {
		template = keyValueParam[1]
	}
	return template
}
