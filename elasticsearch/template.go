package elasticsearch

import (
	"bytes"
	"github.com/Masterminds/sprig"
	"io/ioutil"
)
import "text/template"



func Rendertemplate(filename string,fields string, bufferOutputFile *bytes.Buffer) {
	templateFile, err := ioutil.ReadFile(filename)
	template := template.Must(template.New("").Funcs(sprig.TxtFuncMap()).Parse(string(templateFile)))
	if err != nil {
		panic(err)
	}
	context := map[string]string{
		"fields": fields,
	}
	err = template.Execute(bufferOutputFile, context)
	if err != nil {
		panic(err)
	}
}
