package elasticsearch

import (
	"bytes"
	"encoding/json"
	"github.com/Masterminds/sprig"
	"io/ioutil"
)
import "text/template"

func RenderTemplate(templateInPath string, fields string,jsonKeysDotNotation []string) *bytes.Buffer {
	var bufferOutputFile bytes.Buffer
	templateFile, err := ioutil.ReadFile(templateInPath)
	template := template.Must(template.New("").Funcs(sprig.TxtFuncMap()).Parse(string(templateFile)))
	if err != nil {
		panic(err)
	}
	jsonKeys,err:=json.Marshal(jsonKeysDotNotation)
	if err != nil {
		panic(err)
	}
	context := map[string]interface{}{
		"fields": fields,
		"jsonkeys": string(jsonKeys),
	}
	err = template.Execute(&bufferOutputFile, context)
	if err != nil {
		panic(err)
	}
	return &bufferOutputFile
}
