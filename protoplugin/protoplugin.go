package protoplugin

import (
	"bytes"
	"fmt"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"
	"strings"
)

const (
	templatePath           = "examples/elasticsearch/input/mapping.template"
	templateInPathParam        = "template_in_path"
	templateOutFilenameParam        = "template_out_filename"
	templateOutputFileName = "template.json"
)

type Parameter struct {
	TemplateInPath string
	TemplateOutFilename string
}

func ReqInit(data []byte) *pluginpb.CodeGeneratorRequest {
	req := &pluginpb.CodeGeneratorRequest{}
	if err := proto.Unmarshal(data, req); err != nil {
		panic(err)
	}
	return req
}
func OutputStructSerialized(templateOutFilename string,buffer *bytes.Buffer) []byte {
	resp := &pluginpb.CodeGeneratorResponse{}
	bufferStr := buffer.String()
	outputFile := pluginpb.CodeGeneratorResponse_File{Name: &templateOutFilename, Content: &bufferStr}
	resp.File = []*pluginpb.CodeGeneratorResponse_File{&outputFile}
	output, err := proto.Marshal(resp)
	if err != nil {
		panic(err)
	}
	return output
}

func GetParameters(param string) (Parameter,error) {
	parameter := Parameter{}
	var err error

	keyValueParams := strings.Split(param, ";")
	for _,keyValueParam := range keyValueParams {
		tab := strings.Split(keyValueParam, "=")
		if tab[0] == templateInPathParam {
			parameter.TemplateInPath = tab[1]
		} else if tab[0] == templateOutFilenameParam {
			parameter.TemplateOutFilename = tab[1]
		} else {
			err=fmt.Errorf("bad parameter %v",param)
		}
	}
	return parameter,err
}
