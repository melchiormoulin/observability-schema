package main

import (
	"fmt"
	"github.com/melchiormoulin/observability-schema/elasticsearch"
	"github.com/melchiormoulin/observability-schema/protoplugin"
	"io/ioutil"
	"os"
)

func main() {
	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(fmt.Errorf("error in reading stdin %+v", err))
	}
	input := protoplugin.ReqInit(data)
	parameters,err := protoplugin.GetParameters(input.GetParameter())
	if err!= nil {
		panic(fmt.Errorf("bad template_in path %v",err))
	}
	fmt.Fprintf(os.Stderr, "generating for file %+v with params %+v\n", input.FileToGenerate, input.GetParameter())
	mapping := elasticsearch.MappingInit(true, "  ", "")
	fieldsDefinition := mapping.FieldsDefinition(input.GetProtoFile())
	buffer := elasticsearch.RenderTemplate(parameters.TemplateInPath, fieldsDefinition,mapping.JsonDotPaths)
	os.Stdout.Write(protoplugin.OutputStructSerialized(parameters.TemplateOutFilename,buffer))
}
