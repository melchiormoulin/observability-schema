package protoplugin

import (
	"bytes"
	"github.com/golang/protobuf/proto"
	plugin_go "github.com/golang/protobuf/protoc-gen-go/plugin"
	"strings"
)

func ReqInit(data []byte) *plugin_go.CodeGeneratorRequest {

	req := &plugin_go.CodeGeneratorRequest{}

	if err := proto.Unmarshal(data, req); err != nil {
		panic(err)
	}
	return req

}
func OutputStructSerialized(buffer *bytes.Buffer) []byte {
	resp := &plugin_go.CodeGeneratorResponse{}
	bufferStr := buffer.String()
	fileName := "template.json"
	outputFile := plugin_go.CodeGeneratorResponse_File{Name: &fileName, Content: &bufferStr}
	resp.File = []*plugin_go.CodeGeneratorResponse_File{&outputFile}
	output, err := proto.Marshal(resp)
	if err != nil {
		panic(err)
	}
	return output
}

func TemplatePath(param string) string {
	keyValueParam := strings.Split(param, "=")
	template := "examples/elasticsearch/mapping.template"
	if len(keyValueParam) == 2 && keyValueParam[0] == "template_in" {
		template = keyValueParam[1]
	}
	return template
}